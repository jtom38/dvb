package discovery_test

import (
	"testing"

	"github.com/jtom38/dvb/domain"
	"github.com/jtom38/dvb/services/discovery"
)

func getConfig() domain.Config {
	var containers []domain.ContainerDocker
	containers = append(containers, domain.ContainerDocker{
		Name:      "webdav",
		Directory: "/var/lib/dav",
		Tar: domain.ConfigContainerTar{
			Directory: "{{PWD}}",
			Pattern:   "data-{{DATE}}",
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
				Path: "{{USERDIR}}",
			},
		},
	}

	return c
}

func TestReconNewBackupDetails(t *testing.T) {
	config := getConfig()
	c := discovery.NewReconClient(config)
	backupDetails, err := c.NewBackupDetails(config.Backup.Docker[0].Directory, config.Backup.Docker[0].Name, config.Backup.Docker[0].Tar.Directory)
	if err != nil {
		t.Error(err)
	}

	if backupDetails.FullFilePath == "" {
		t.Error("Full File Path was missing and not generated")
	}
}

func TestReconValidateBackupDetails(t *testing.T) {
	config := getConfig()
	c := discovery.NewReconClient(config)
	backupDetails, err := c.NewBackupDetails(config.Backup.Docker[0].Directory, config.Backup.Docker[0].Name, config.Backup.Docker[0].Tar.Directory)
	if err != nil {
		t.Error(err)
	}

	err = c.ValidateBackupDetails(backupDetails)
	if err != nil {
		t.Error(err)
	}
}

func TestReconGetLocalDestDetails(t *testing.T) {
	config := getConfig()

	c := discovery.NewReconClient(config)
	backupDetails, err := c.NewBackupDetails(config.Backup.Docker[0].Directory, config.Backup.Docker[0].Name, config.Backup.Docker[0].Tar.Directory)
	if err != nil {
		t.Error(err)
	}
	_, err = c.GetLocalDestDetails(discovery.LocalDetailsParam{
		Container:     config.Backup.Docker[0],
		BackupDetails: backupDetails,
		DestLocal:     config.Destination.Local,
	})

	if err != nil {
		t.Error(err)
	}
}

func TestReconValidateLocalDestDetails(t *testing.T) {
	config := getConfig()

	c := discovery.NewReconClient(config)
	backupDetails, err := c.NewBackupDetails(config.Backup.Docker[0].Directory, config.Backup.Docker[0].Name, config.Backup.Docker[0].Tar.Directory)
	if err != nil {
		t.Error(err)
	}

	dest, err := c.GetLocalDestDetails(discovery.LocalDetailsParam{
		Container:     config.Backup.Docker[0],
		BackupDetails: backupDetails,
		DestLocal:     config.Destination.Local,
	})
	if err != nil {
		t.Error(err)
	}

	err = c.ValidateLocalDestDetails(dest)
	if err != nil {
		t.Error(err)
	}
}

func TestReconDockerScout(t *testing.T) {
	config := getConfig()

	c := discovery.NewReconClient(config)
	details, err := c.DockerScout(config.Backup.Docker[0])
	if err != nil {
		t.Error(err)
	}

	if details.Backup.FullFilePath == "" {
		t.Error("Backup.FullFilePath is missing")
	}

	if details.Dest.Local.FullFilePath == "" {
		t.Error("Dest.Local.FullFilePath is missing")
	}
}
