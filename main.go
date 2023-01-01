package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jtom38/dvb/domain"
	"github.com/jtom38/dvb/services"
	"github.com/jtom38/dvb/services/alerts"
	"github.com/jtom38/dvb/services/cmd"
	"github.com/jtom38/dvb/services/dest"
	"github.com/jtom38/dvb/services/targets"
	"gopkg.in/yaml.v3"
)

func main() {
	cmd.Execute()
}

func ProcessDockerContainers(config domain.Config, container domain.ContainerDocker) error {
	logs := domain.NewLogs()
	logs.Add("The container backup has started.")

	details := domain.RunDetails{
		ContainerName: container.Name,
	}

	// Based on the destination path, lets figure out what we should name the file
	recon := services.NewReconClient(config)
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
		SendAlert(config.Alert, logs)
		return err
	}
	logs.Add(fmt.Sprintf("Backup was created. '%v.tar'", details.Backup.FileName))

	err = MoveFile(details, config.Destination)
	if err != nil {
		logs.Error(err)
		SendAlert(config.Alert, logs)
		return err
	}
	logs.Add("Backup was moved.")

	// Check if we need to remove any old backups
	if config.Destination.Retain.Days == 0 {
		return nil
	}
	if config.Destination.Local.Path == "" {
		return nil
	}

	retain := dest.NewLocalRetainClient(config.Destination.Local, details.ContainerName, config.Destination.Retain.Days)
	err = retain.Check(".tar")
	if err != nil {
		logs.Error(err)
		SendAlert(config.Alert, logs)
		return err
	}

	logs.Add(fmt.Sprintf("No errors reported backing up '%v'", container.Name))
	SendAlert(config.Alert, logs)
	return nil
}

func MustLoadConfig(path string) domain.Config {
	var config domain.Config

	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func MoveFile(details domain.RunDetails, config domain.ConfigDest) error {
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

func SendAlert(config domain.ConfigAlert, logs domain.Logs) {
	discordAlert := alerts.NewDiscordAlertClient(config.Discord.Webhooks, config.Discord.Username)
	m := strings.Join(logs.Message, "\r\n> ")
	discordAlert.ReplaceContent(m)
	discordAlert.Send()
}
