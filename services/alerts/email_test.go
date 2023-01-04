package alerts_test

import (
	"os"
	"testing"

	"github.com/jtom38/dvb/domain"
	"github.com/jtom38/dvb/services/alerts"
	"gopkg.in/yaml.v3"
)

func getEmailConfig() (domain.ConfigAlertEmail, error) {
	var config domain.ConfigAlertEmail

	// Check if we have a config file we can load
	_, err := os.Stat("config.yaml")
	if err == nil {
		config, err = readConfig("config.yaml")
		if err != nil {
			return config, err
		}
	}

	return config, nil
}

func readConfig(name string) (domain.ConfigAlertEmail, error) {
	var config domain.Config

	content, err := os.ReadFile("config.yaml")
	if err != nil {
		return config.Alert.Email, err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return config.Alert.Email, err
	}

	return config.Alert.Email, nil
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
