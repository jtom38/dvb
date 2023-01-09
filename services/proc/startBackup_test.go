package proc_test

import (
	"testing"

	"github.com/jtom38/dvb/domain"
	"github.com/jtom38/dvb/services/proc"
)

func TestStartBackup(t *testing.T) {
	c := proc.NewStartBackupClient(proc.StartBackupParams{
		Daemon:     true,
		ConfigPath: "config.yaml",
	})

	err := c.RunProcess()
	if err != nil {
		t.Error(err)
	}
}

func TestStartBackupDiscordAlert(t *testing.T) {
	c := proc.NewStartBackupClient(proc.StartBackupParams{
		Daemon: false,
		ConfigPath: "config.yaml",
	})

	c.SendAlert(domain.ConfigAlert{
		SendOnlyOnError: false,
		Discord: domain.ConfigAlertDiscord{
			Username: "unit-test",
			Webhooks: "",
		},
	})
}
