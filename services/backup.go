package services

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jtom38/dvb/domain"
)

type BackupClient struct {
	FileExtension string
}

func NewBackupClient() BackupClient {
	c := BackupClient{
		FileExtension: "tar",
	}
	return c
}

// This will return the location of the new file on disk if it was successful
func (c BackupClient) BackupDockerVolume(details domain.RunDetails, config domain.ContainerDocker) (domain.RunDetails, error) {
	client := NewDockerCliClient()

	log.Printf("> Checking for %v", config.Name)
	inspect, err := client.InspectContainer(config.Name)
	if err != nil {
		return details, errors.New(inspect)
	}

	log.Print("> Stopping container")
	out, err := client.StopContainer(config.Name)
	if err != nil {
		return details, errors.New(out)
	}

	log.Print("> Determining backup name")

	// Check if we are going to dump into the working directory
	tarDirectory, err := c.GetDirectoryPath(config.Tar.Directory)
	if err != nil {
		return details, err
	}
	details.BackupDirectory = tarDirectory

	backupName := c.GetValidFileName(config.Tar, tarDirectory)
	details.BackupFileName = backupName
	details.BackupExtension = ".tar"

	log.Printf("Backup will generate as '%v'", backupName)

	// backup volume
	log.Print("> Starting to backup the volume")
	backedResults, err := client.BackupDockerVolume(BackupVolumeParams{
		ContainerName:  config.Name,
		BackupFolder:   tarDirectory,
		BackupFilename: backupName,
		TargetFolder:   config.Directory,
	})
	if err != nil {
		return details, errors.New(backedResults)
	}

	path := fmt.Sprintf("%v/%v.%v", tarDirectory, backupName, c.FileExtension)
	_, err = os.Stat(path)
	if err != nil {
		return details, err
	}

	// The file exists, so we will return the location we tested
	details.BackupPath = path

	// start container
	log.Print("> Starting container")
	out, err = client.StartContainer(config.Name)
	if err != nil {
		return details, errors.New(out)
	}

	return details, nil
}

func (c BackupClient) GetDirectoryPath(value string) (string, error) {
	if value == "$PWD" {
		workingDirectory, err := os.Getwd()
		if err != nil {
			return "", err
		}

		return workingDirectory, nil
	}
	return value, nil
}

func (c BackupClient) ReplaceDatePlaceholder(pattern string) string {
	backupName := pattern
	todayString := time.Now().Format("20060102")
	backupName = strings.ReplaceAll(backupName, "{{date}}", todayString)
	return backupName
}

// This will update the filename if one already exists with a number appended
func (c BackupClient) GetValidFileName(config domain.ConfigContainerTar, directory string) string {
	var tempName string
	var t string

	ogName := config.Pattern
	backupName := config.Pattern

	if strings.Contains(config.Pattern, "{{date}}") {
		backupName = c.ReplaceDatePlaceholder(config.Pattern)
		ogName = backupName
	}

	i := 0
	tempName = fmt.Sprintf("%v.%v", ogName, c.FileExtension)
	t = fmt.Sprintf("%v/%v", directory, tempName)

	for {

		_, err := os.Stat(t)
		if err != nil {
			return backupName
		}

		backupName = fmt.Sprintf("%v.%v", ogName, i)
		tempName = fmt.Sprintf("%v.%v", backupName, c.FileExtension)
		t = fmt.Sprintf("%v/%v", directory, tempName)
		i = i + 1
	}
}

func GetValidFileName(directory, name, extension string) string {
	var tempName string
	var t string

	i := 0
	tempName = fmt.Sprintf("%v.%v", name, extension)
	t = fmt.Sprintf("%v/%v", directory, tempName)

	for {

		_, err := os.Stat(t)
		if err != nil {
			return tempName
		}

		tempName = fmt.Sprintf("%v.%v.%v", name, i, extension)
		t = fmt.Sprintf("%v/%v", directory, tempName)
		i = i + 1
	}
}
