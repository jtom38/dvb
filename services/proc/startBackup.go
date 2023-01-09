package proc

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"

	"github.com/jtom38/dvb/domain"
	"github.com/jtom38/dvb/services/alerts"
	"github.com/jtom38/dvb/services/dest"
	"github.com/jtom38/dvb/services/discovery"
	"github.com/jtom38/dvb/services/lib"
	"github.com/jtom38/dvb/services/targets"
)

type StartBackupParams struct {
	ConfigPath string
	Daemon     bool
}

type StartBackupClient struct {
	Config domain.Config
	Params StartBackupParams
}

func NewStartBackupClient(params StartBackupParams) StartBackupClient {
	return StartBackupClient{
		Params: params,
	}
}

func (c StartBackupClient) RunProcess() error {
	// Get the config file loaded into memory.
	// Need this is we are running once or as a daemon
	config, err := c.LoadConfig(c.Params.ConfigPath)
	if err != nil {
		return err
	}
	c.SetConfig(config)

	// If daemon is requested from param or config check
	if c.Config.Daemon.Cron != "" {
		log.Print("Daemon mode was requested.")
		log.Printf("Backups will start at '%v'", c.Config.Daemon.Cron)
		c.RunDaemon()
	} else {
		err = c.RunSingle()
		if err != nil {
			return nil
		}
	}

	return nil
}

// This runs the tool once and closes down once its finished.
func (c StartBackupClient) RunSingle() error {
	// Process all requested docker containers
	for _, container := range c.Config.Backup.Docker {
		err := c.ProcessDockerContainers(container)
		if err != nil {
			log.Print(err)
		}
	}

	return nil
}

func (c StartBackupClient) RunDaemon() error {
	// Set when we want to run the backup job
	cronClient := cron.New()
	_, err := cronClient.AddFunc(c.Config.Daemon.Cron, func() {
		log.Print("Cron was triggered")
		go c.RunSingle()
	})
	if err != nil {
		log.Print(err)
	}

	cronClient.Start()

	// Check if we get a request to stop the app
	ch := make(chan os.Signal, 6)
	signal.Notify(ch,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	for {
		req := <-ch
		switch req {
		case syscall.SIGTERM:
			fallthrough
		case syscall.SIGINT:
			cronClient.Stop()
			signal.Stop(ch)
			return nil
		case syscall.SIGQUIT:
			signal.Stop(ch)
			os.Exit(0)
		}
	}
}

func (c *StartBackupClient) SetConfig(config domain.Config) {
	c.Config = config
}

func (c StartBackupClient) ProcessDockerContainers(container domain.ContainerDocker) error {
	logs := domain.NewLogs()
	logs.Add("The container backup has started.")

	details := domain.RunDetails{
		ContainerName: container.Name,
	}

	// Based on the destination path, lets figure out what we should name the file
	recon := discovery.NewReconClient(c.Config)
	details, err := recon.DockerScout(container)
	if err != nil {
		logs.Error(err)
		return err
	}

	// Start the backup process on the container
	backupDockerClient := targets.NewDockerClient()
	err = backupDockerClient.BackupDockerVolume(details, container)
	if err != nil {
		logs.Error(err)
		c.SendAlert(c.Config.Alert, logs, true)
		return err
	}
	logs.Add(fmt.Sprintf("Backup was created. '%v.tar'", details.Backup.FileName))

	// run any post reboot requests after a backup was made
	c.postRebootContainer(container.Post.Reboot)

	err = c.MoveFile(details, c.Config.Destination)
	if err != nil {
		logs.Error(err)
		c.SendAlert(c.Config.Alert, logs, true)
		return err
	}
	logs.Add("Backup was moved.")

	// Check if we need to remove any old backups
	if c.Config.Destination.Retain.Days == 0 {
		return nil
	}
	if c.Config.Destination.Local.Path == "" {
		return nil
	}

	retain := dest.NewLocalRetainClient(c.Config.Destination.Local, details.ContainerName, c.Config.Destination.Retain.Days)
	err = retain.Check(".tar")
	if err != nil {
		logs.Error(err)
		c.SendAlert(c.Config.Alert, logs, true)
		return err
	}

	logs.Add(fmt.Sprintf("No errors reported backing up '%v' 🎉", container.Name))

	c.SendAlert(c.Config.Alert, logs, false)
	return nil
}

func (c StartBackupClient) LoadConfig(path string) (domain.Config, error) {
	var config domain.Config

	content, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func (c StartBackupClient) MoveFile(details domain.RunDetails, config domain.ConfigDest) error {
	var err error
	if details.Dest.Local.Directory != "" {
		local := dest.NewLocalClient(details.Backup.FileName, details.Backup.FullFilePath, details.ContainerName, config.Local.Path)
		err = local.Move(details)
		if err != nil {
			return err
		}

		// Remove the old file
		err = os.Remove(details.Backup.FullFilePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c StartBackupClient) SendAlert(config domain.ConfigAlert, logs domain.Logs, isError bool) {
	var err error

	if len(config.Discord.Webhooks) >= 1 {
		log.Print("Sending discord alert")
		err = c.sendDiscordAlert(config.Discord, logs, isError)
		if err != nil {
			log.Print(err)
		}
	}

	if config.Email.Account.Username != "" && config.Email.Account.Password != "" {
		log.Print("Sending email alert")
		err = c.sendEmailAlert(config.Email, logs)
		if err != nil {
			log.Print(err)
		}
	}
}

func (c StartBackupClient) sendDiscordAlert(config domain.ConfigAlertDiscord, logs domain.Logs, isError bool) error {
	var color string

	discordAlert := alerts.NewDiscordEmbedMessage(config)
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	if isError {
		color = alerts.DiscordErrorColor
	} else {
		color = alerts.DiscordSuccessColor
	}

	discordAlert.AppendFields(alerts.DiscordEmbedFieldParams{
		Name:   "Server",
		Value:  hostname,
		Inline: true,
	})

	discordAlert.AppendFields(alerts.DiscordEmbedFieldParams{
		Name:   "Container",
		Value:  "Placeholder",
		Inline: true,
	})

	m := strings.Join(logs.Message, "\n")
	discordAlert.SetBody(alerts.DiscordEmbedBodyParams{
		Title:       "Backup Results",
		Color:       color,
		Description: m,
	})

	err = discordAlert.SendPayload()
	if err != nil {
		return err
	}

	log.Print("> OK")
	return nil
}

func (c StartBackupClient) sendEmailAlert(config domain.ConfigAlertEmail, logs domain.Logs) error {
	m := strings.Join(logs.Message, "<br>")

	client := alerts.NewSmtpClient(config)
	client.SetSubject(alerts.EmailSubjectSuccess)
	client.SetBody(m)
	err := client.SendAlert()
	if err != nil {
		return err
	}
	log.Print("> OK")
	return nil
}

func (c StartBackupClient) postRebootContainer(names []string) {
	if len(names) == 0 {
		return
	}

	client := lib.NewDockerCliClient()
	log.Print("Running Post Reboot requests")

	for _, name := range names {
		log.Printf("Stopping '%v'", name)
		output, err := client.StopContainer(name)
		if err != nil {
			log.Print(output)
		}

		log.Printf("Starting '%v'", name)
		output, err = client.StartContainer(name)
		if err != nil {
			log.Print(output)
		}
	}
}
