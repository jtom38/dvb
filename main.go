package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jtom38/dvb/domain"
	"github.com/jtom38/dvb/services"
	"gopkg.in/yaml.v3"
)

func main() {
	var flagConfigPath string

	flag.StringVar(&flagConfigPath, "config", "config.yaml", "Defines what yaml config file to use")
	flag.Parse()

	// load the config file into memory
	config := MustLoadConfig(flagConfigPath)

	for _, container := range config.Backup.Docker {
		ProcessDockerContainers(config, container)
	}
}

func ProcessDockerContainers(config domain.Config, container domain.ContainerDocker,) error {
	logs := domain.Logs{}
	logs.Add("The container backup has started.")

	details := domain.RunDetails{
		ContainerName:   container.Name,
	}

	// TODO Review the storage location
	err := ReviewStorageLocation(config.Destination)
	if err != nil {
		logs.Error(err)
		SendAlert(config.Alert, logs)
		return err
	}
	logs.Add(fmt.Sprintf("Was able to access '%v'", config.Destination.Local.Path))

	// Start the backup process on the container
	backupClient := services.NewBackupClient()
	details, err = backupClient.BackupDockerVolume(details, container)
	if err != nil {
		logs.Error(err)
		SendAlert(config.Alert, logs)
		return err
	}
	logs.Add(fmt.Sprintf("Backup was created. '%v.tar'", details.BackupFileName))

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

	retain := services.NewRetainClient(config.Destination.Local, details.ContainerName, config.Destination.Retain.Days)
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

func ReviewStorageLocation(config domain.ConfigDest) error {
	if config.Local.Path != "" {
		log.Printf("Moving backup to %v", config.Local.Path)
		_, err := os.Stat(config.Local.Path)
		if err != nil {
			return err
		}
	}

	return nil
}

func MoveFile(details domain.RunDetails, config domain.ConfigDest) error {
	var err error
	if config.Local.Path != "" {
		local := services.NewMoveClient(details.BackupFileName, details.BackupPath, details.ContainerName, config.Local.Path)
		err = local.Move()
		if err != nil {
			return err
		}

		// Remove the old file
		err = os.Remove(details.BackupPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func SendAlert(config domain.ConfigAlert, logs domain.Logs) {
	discordAlert := services.NewDiscordAlertClient(config.Discord.Webhooks, config.Discord.Username)
	m := strings.Join(logs.Message, "\r\n> ")
	discordAlert.ReplaceContent(m)
	discordAlert.Send()
}
