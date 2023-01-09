package alerts_test

import (
	"os"
	"testing"

	"github.com/jtom38/dvb/domain"
	"github.com/jtom38/dvb/services/alerts"
	"gopkg.in/yaml.v3"
)

func getDiscordConfig() (domain.ConfigAlert, error) {
	var config domain.ConfigAlert

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

func readConfig(name string) (domain.ConfigAlert, error) {
	var config domain.Config

	content, err := os.ReadFile("config.yaml")
	if err != nil {
		return config.Alert, err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return config.Alert, err
	}

	return config.Alert, nil
}

func TestDiscordBasicMessage(t *testing.T) {
	cfg, err := getDiscordConfig()
	if err != nil {
		t.Error(err)
	}

	c := alerts.NewDiscordBasicAlertClient(cfg.Discord)
	c.SetContent("hi")
	err = c.SendPayload()
	if err != nil {
		t.Error(err)
	}
}

func TestDiscordEmbedMessage(t *testing.T) {
	cfg, err := getDiscordConfig()
	if err != nil {
		t.Error(err)
	}

	c := alerts.NewDiscordEmbedMessage(cfg.Discord)
	c.SetAuthor("bingo", "", "")
	c.SetBody(alerts.DiscordEmbedBodyParams{
		Title:       "Unit...",
		Description: "...TESTING!!!",
		Color:       alerts.DiscordSuccessColor,
		TitleURL:    "https://github.com/jtom38/dvb",
	})
	c.AppendFields(alerts.DiscordEmbedFieldParams{
		Name:  "New",
		Value: "Thing",
	})

	c.SetFooter(alerts.DiscordEmbedFooterParams{
		Text: "DVB",
	})

	err = c.SendPayload()
	if err != nil {
		t.Error(err)
	}

}
