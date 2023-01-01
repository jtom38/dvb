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
	"github.com/jtom38/dvb/services/common"
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
	var (
		err       error
		res       domain.RunDetails
		backup    domain.RunBackupDetails
		destLocal domain.RunDetailsDestLocal
	)

	// set the container name
	res.ContainerName = container.Name

	for {
		backup, err = c.NewBackupDetails(container.Directory, container.Name, container.Tar.Directory)
		if err != nil {
			return res, err
		}

		// make sure that we are able to use the generated name and path
		err = c.ValidateBackupDetails(backup)
		if err == nil {
			res.Backup = backup
			break
		}
	}

	// build the details for our dest
	if c.config.Destination.Local.Path != "" {
		destLocal, err = c.getLocalDestValues(GetLocalDestValuesParam{
			Container: container,
			Backup:    backup,
			Dest:      c.config.Destination.Local,
		})
		if err != nil {
			return res, err
		}
		res.Dest.Local = destLocal
	}

	return res, nil
}

type GetLocalDestValuesParam struct {
	Container domain.ContainerDocker
	Backup    domain.RunBackupDetails
	Dest      domain.ConfigDestLocal
}

// This controls the loop that will check to make sure the generated values are all correct.
// If it it finds a combo that will not work then it will generate a new value
func (c ReconClient) getLocalDestValues(params GetLocalDestValuesParam) (domain.RunDetailsDestLocal, error) {
	var (
		dest    domain.RunDetailsDestLocal
		err     error
		counter int
	)

	for {
		dest, err = c.GetLocalDestDetails(LocalDetailsParam{
			Container:     params.Container,
			BackupDetails: params.Backup,
			DestLocal:     params.Dest,
			Counter:       counter,
		})
		if err != nil {
			return dest, err
		}

		err = c.ValidateLocalDestDetails(dest)
		if err == nil {
			return dest, nil
		}
		counter = counter + 1
	}

}

// This will generate new backup details and store them.
func (c ReconClient) NewBackupDetails(targetDir, folderName, destDir string) (domain.RunBackupDetails, error) {
	var d domain.RunBackupDetails

	name := uuid.NewString()
	ext := "tar"
	nameAndExt := fmt.Sprintf("%v.%v", name, ext)
	d.ServiceName = folderName

	destDir = filepath.Join(destDir, folderName)

	destDir, err := common.ReplaceAllConfigVariables(destDir)
	if err != nil {
		return d, err
	}

	d = domain.RunBackupDetails{
		TargetDirectory:       targetDir,
		LocalDirectory:        destDir,
		FileName:              name,
		Extension:             ext,
		FileNameWithExtension: nameAndExt,
		FullFilePath:          filepath.Join(destDir, nameAndExt),
		ServiceName:           folderName,
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

type LocalDetailsParam struct {
	Counter       int
	Container     domain.ContainerDocker
	BackupDetails domain.RunBackupDetails
	DestLocal     domain.ConfigDestLocal
}

// func (c ReconClient) GetLocalDestDetails(config domain.ConfigDestLocal, backup domain.RunBackupDetails) domain.RunDetailsDestLocal {
func (c ReconClient) GetLocalDestDetails(params LocalDetailsParam) (domain.RunDetailsDestLocal, error) {
	var d domain.RunDetailsDestLocal
	var fileName string

	dir := filepath.Join(params.DestLocal.Path, params.BackupDetails.ServiceName)

	dir, err := common.ReplaceAllConfigVariables(dir)
	if err != nil {
		return d, err
	}

	// append the counter value to the name
	fileName = fmt.Sprintf("%v.%v", params.Container.Tar.Pattern, params.Counter)

	fileName, err = common.ReplaceAllConfigVariables(fileName)
	if err != nil {
		return d, err
	}

	fileWithExt := fmt.Sprintf("%v.%v", fileName, params.BackupDetails.Extension)

	d = domain.RunDetailsDestLocal{
		Directory:             dir,
		FileName:              fileName,
		Extension:             params.BackupDetails.Extension,
		FileNameWithExtension: fileWithExt,
		FullFilePath:          filepath.Join(dir, fileWithExt),
	}
	return d, nil
}

func (c ReconClient) ValidateLocalDestDetails(details domain.RunDetailsDestLocal) error {
	_, err := os.Stat(details.FullFilePath)
	if err != nil {
		return nil
	}
	return errors.New(ErrFileAlreadyExists)
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

//func (c ReconClient) getDirectoryPath(value string) (string, error) {
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
