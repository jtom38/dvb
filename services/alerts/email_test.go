package alerts_test

import (
	"os"
	"testing"

	"github.com/jtom38/dvb/domain"
	"github.com/jtom38/dvb/services/alerts"
)

func getEmailConfig() (domain.ConfigAlertEmail, error) {
	var config domain.ConfigAlert

	// Check if we have a config file we can load
	_, err := os.Stat("config.yaml")
	if err == nil {
		config, err = readConfig("config.yaml")
		if err != nil {
			return config.Email, err
		}
	}

	return config.Email, nil
}

func TestEmailSetTo(t *testing.T) {
	cfg, err := getEmailConfig()
	if err != nil {
		t.Error(err)
	}

	client := alerts.NewSmtpClient(cfg)
	client.SetSubject("Unit Test")
	client.SetBody("Hello World")
	err = client.SendAlert()
	if err != nil {
		t.Error(err)
	}
}
