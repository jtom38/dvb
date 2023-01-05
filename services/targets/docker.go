package targets

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jtom38/dvb/domain"
	"github.com/jtom38/dvb/services/common"
	"github.com/jtom38/dvb/services/lib"
)

type DockerClient struct {
	FileExtension string
}

func NewDockerClient() DockerClient {
	c := DockerClient{
		FileExtension: "tar",
	}
	return c
}

// This will return the location of the new file on disk if it was successful
func (c DockerClient) BackupDockerVolume(details domain.RunDetails, config domain.ContainerDocker) error {
	client := lib.NewDockerCliClient()

	log.Printf("> Checking for %v", config.Name)
	inspect, err := client.InspectContainer(config.Name)
	if err != nil {
		return errors.New(inspect)
	}

	log.Print("> Stopping container")
	out, err := client.StopContainer(config.Name)
	if err != nil {
		return errors.New(out)
	}

	log.Printf("Backup will generate as '%v'", details.Backup.FileNameWithExtension)

	// backup volume
	log.Print("> Starting to backup the volume")
	backedResults, err := client.BackupDockerVolume(lib.DockerBackupVolumeParams{
		ContainerName:  config.Name,
		BackupFolder:   details.Backup.LocalDirectory,
		BackupFilename: details.Backup.FileName,
		TargetFolder:   details.Backup.TargetDirectory,
	})
	if err != nil {
		return errors.New(backedResults)
	}

	// start container
	log.Print("> Starting container")
	out, err = client.StartContainer(config.Name)
	if err != nil {
		return errors.New(out)
	}

	return nil
}

//func (c DockerClient) GetDirectoryPath(value string) (string, error) {
//	if value == "$PWD" {
//		workingDirectory, err := os.Getwd()
//		if err != nil {
//			return "", err
//		}
//
//		return workingDirectory, nil
//	}
//	return value, nil
//}

//func (c DockerClient) ReplaceDatePlaceholder(pattern string) string {
//	backupName := pattern
//	todayString := time.Now().Format("20060102")
//	backupName = strings.ReplaceAll(backupName, "{{date}}", todayString)
//	return backupName
//}

// This will update the filename if one already exists with a number appended
func (c DockerClient) GetValidFileName(config domain.ConfigContainerTar, directory string) (string, error) {
	var tempName string
	var t string

	ogName := config.Pattern
	backupName := config.Pattern

	if strings.Contains(config.Pattern, "{{date}}") {
		backupName, err := common.ReplaceAllConfigVariables(config.Pattern)
		if err != nil {
			return backupName, err
		}
		//backupName = c.ReplaceDatePlaceholder(config.Pattern)
		ogName = backupName
	}

	i := 0
	tempName = fmt.Sprintf("%v.%v", ogName, c.FileExtension)
	t = fmt.Sprintf("%v/%v", directory, tempName)

	for {

		_, err := os.Stat(t)
		if err != nil {
			return backupName, nil
		}

		backupName = fmt.Sprintf("%v.%v", ogName, i)
		tempName = fmt.Sprintf("%v.%v", backupName, c.FileExtension)
		t = fmt.Sprintf("%v/%v", directory, tempName)
		i = i + 1
	}
}
