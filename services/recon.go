package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jtom38/dvb/domain"
)

const (
	ErrFileAlreadyExists = "a file already exists with the requested path"
)

type ReconClient struct {
	config  domain.Config
	details domain.RunDetails
}

// The recon client goes and check to see what we need before any actions are taken.
// Once the recon has done its job all the other services will be able run without needing to validate along the way.
func NewReconClient(config domain.Config) ReconClient {
	c := ReconClient{
		config: config,
	}
	return c
}

// Scout is the main logic loop and reports back its findings to
func (c ReconClient) DockerScout(container domain.ContainerDocker) (domain.RunDetails, error) {
	var res domain.RunDetails
	var backup domain.RunBackupDetails

	// set the container name
	res.ContainerName = container.Name

	for {
		backup, err := c.NewBackupDetails(container.Directory)
		if err != nil {
			return res, err
		}

		// make sure that we are able to use the generated name and path
		err = c.ValidateBackupDetails(backup)
		if err != nil {
			break
		}
	}
	res.Backup = backup

	// build the details for our dest
	if c.config.Destination.Local.Path != "" {

	}

	// Generate Local Dest details

	tarDir, err := c.getDirectoryPath(container.Tar.Directory)
	if err != nil {
		return c.details, err
	}
	c.details.Backup.Directory = tarDir
	//c.TestBackupName()

	return c.details, nil
}

// This will generate new backup details and store them.
func (c ReconClient) NewBackupDetails(backupDir string) (domain.RunBackupDetails, error) {
	var d domain.RunBackupDetails

	name := uuid.NewString()
	ext := "tar"
	nameAndExt := fmt.Sprintf("%v.%v", name, ext)
	backupDir, err := ReplaceAllConfigVariables(backupDir)
	if err != nil {
		return d, err
	}

	d = domain.RunBackupDetails{
		Directory:             backupDir,
		FileName:              name,
		Extension:             ext,
		FileNameWithExtension: nameAndExt,
		FullFilePath:          filepath.Join(backupDir, nameAndExt),
	}

	return d, nil
}

func (c ReconClient) ValidateBackupDetails(details domain.RunBackupDetails) error {
	_, err := os.Stat(details.FullFilePath)
	if err != nil {
		return nil
	}
	return errors.New(ErrFileAlreadyExists)
}

func (c ReconClient) GetLocalDestDetails(config domain.ConfigDestLocal, backup domain.RunBackupDetails) (domain.RunDetailsDestLocal) {
	var d domain.RunDetailsDestLocal

	dir, err := ReplaceAllConfigVariables(config.Path)
	if err != nil {
		
	}
	d = domain.RunDetailsDestLocal{
		Directory: dir,
		FileName:  backup.FileName,
		Extension: backup.Extension,

	}
	return d
}

// Test to make sure the path works before we commit the values to the RunDetails.
func (c ReconClient) TestBackupName(dir, name, ext string) error {
	fileName := fmt.Sprintf("%v.%v", name, ext)
	p := filepath.Join(dir, fileName)
	_, err := os.Stat(p)
	if err == nil {
		return errors.New("unable to validate path")
	}

	return nil
}

func (c ReconClient) getDirectoryPath(value string) (string, error) {
	if value == "$PWD" {
		workingDirectory, err := os.Getwd()
		if err != nil {
			return "", err
		}

		return workingDirectory, nil
	}
	return value, nil
}

// This will update the filename if one already exists with a number appended
func (c ReconClient) GetValidFileName(config domain.ConfigContainerTar, directory string) string {
	var tempName string
	var t string

	ogName := config.Pattern
	backupName := config.Pattern

	if strings.Contains(config.Pattern, "{{date}}") {
		backupName = c.replaceDatePlaceholder(config.Pattern)
		ogName = backupName
	}

	i := 0
	tempName = fmt.Sprintf("%v.%v", ogName, c.details.Backup.Extension)
	t = fmt.Sprintf("%v/%v", directory, tempName)

	for {

		_, err := os.Stat(t)
		if err != nil {
			return backupName
		}

		backupName = fmt.Sprintf("%v.%v", ogName, i)
		tempName = fmt.Sprintf("%v.%v", backupName, c.details.Backup.Extension)
		t = fmt.Sprintf("%v/%v", directory, tempName)
		i = i + 1
	}
}

func (c ReconClient) replaceDatePlaceholder(pattern string) string {
	backupName := pattern
	todayString := time.Now().Format("20060102")
	backupName = strings.ReplaceAll(backupName, "{{date}}", todayString)
	return backupName
}
