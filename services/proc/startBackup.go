package proc

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jtom38/dvb/domain"
	"github.com/jtom38/dvb/services/alerts"
	"github.com/jtom38/dvb/services/dest"
	"github.com/jtom38/dvb/services/discovery"
	"github.com/jtom38/dvb/services/targets"
	"gopkg.in/yaml.v3"
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
	config, err := c.LoadConfig(c.Params.ConfigPath)
	if err != nil {
		return err
	}
	c.SetConfig(config)

	// Process all requested docker containers
	for _, container := range c.Config.Backup.Docker {
		err := c.ProcessDockerContainers(container)
		if err != nil {
			log.Print(err)
		}
	}

	return nil
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
		c.SendAlert(c.Config.Alert, logs)
		return err
	}
	logs.Add(fmt.Sprintf("Backup was created. '%v.tar'", details.Backup.FileName))

	err = c.MoveFile(details, c.Config.Destination)
	if err != nil {
		logs.Error(err)
		c.SendAlert(c.Config.Alert, logs)
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
		c.SendAlert(c.Config.Alert, logs)
		return err
	}

	logs.Add(fmt.Sprintf("No errors reported backing up '%v'", container.Name))
	c.SendAlert(c.Config.Alert, logs)
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

func (c StartBackupClient) SendAlert(config domain.ConfigAlert, logs domain.Logs) {
	discordAlert := alerts.NewDiscordAlertClient(config.Discord.Webhooks, config.Discord.Username)
	m := strings.Join(logs.Message, "\r\n> ")
	discordAlert.ReplaceContent(m)
	discordAlert.Send()
}
