package services_test

import (
	"testing"

	"github.com/jtom38/dvb/domain"
	"github.com/jtom38/dvb/services"
)

func getConfig() domain.Config {
	var containers []domain.ContainerDocker
	containers = append(containers, domain.ContainerDocker{
		Name:      "webdav",
		Directory: "/var/lib/dav",
		Tar: domain.ConfigContainerTar{
			Directory: "{{PWD}}",
			Pattern:   "webdav-data-{{date}}",
		},
	})

	c := domain.Config{
		Backup: domain.BackupConfig{
			Docker: containers,
		},
		Destination: domain.ConfigDest{
			Retain: domain.ConfigRetain{
				Days: 10,
			},
			Local: domain.ConfigDestLocal{
				Path: "~",
			},
		},
	}

	return c
}

func TestReconNewBackupDetails(t *testing.T) {
	config := getConfig()
	c := services.NewReconClient(config)
	backupDetails, err := c.NewBackupDetails(config.Backup.Docker[0].Tar.Directory)
	if err != nil {
		t.Error(err)
	}

	if backupDetails.FullFilePath == "" {
		t.Error("Full File Path was missing and not generated")
	}
}

func TestReconValidateBackupDetails(t *testing.T) {
	config := getConfig()
	c := services.NewReconClient(config)
	details, err := c.NewBackupDetails(config.Backup.Docker[0].Tar.Directory)
	if err != nil {
		t.Error(err)
	}

	err = c.ValidateBackupDetails(details)
	if err != nil {
		t.Error(err)
	}

}
