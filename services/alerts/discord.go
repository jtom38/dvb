package alerts

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/jtom38/dvb/domain"
)

const (
	DiscordErrorColor   = 16711680
	DiscordSuccessColor = 65290

	ErrDiscordContentLengthTooLong = "the message length is greater then 2000 characters"
	ErrDiscordEmbedLenthTooLong    = "the message length is greater then 4096 characters"
)

type DiscordBasicAlertClient struct {
	webhooks []string
	username string

	message domain.DiscordBasicMessage
}

func NewDiscordBasicAlertClient(params domain.ConfigAlertDiscord) DiscordBasicAlertClient {
	c := DiscordBasicAlertClient{
		webhooks: params.Webhooks,
		username: params.Username,
	}
	c.message = c.NewMessage()
	return c
}

func (c DiscordBasicAlertClient) NewMessage() domain.DiscordBasicMessage {
	m := domain.DiscordBasicMessage{
		Username: c.username,
	}
	return m
}

func (c *DiscordBasicAlertClient) SetContent(content string) {
	c.message.Content = content
}

func (c DiscordBasicAlertClient) validateWebhooks() error {
	if len(c.webhooks) == 0 {
		return errors.New("no webhooks given to post to")
	}

	// Check to confirm the webhook uri is correct
	for _, uri := range c.webhooks {
		if strings.Contains(uri, "https://discord.com/api/webhooks/") {
			continue
		}

		return errors.New("invalid uri given to post to")
	}

	return nil
}

func (c DiscordBasicAlertClient) SendPayload() error {
	err := c.validateWebhooks()
	if err != nil {
		return err
	}

	if len(c.message.Content) >= 2000 {
		return errors.New(ErrDiscordContentLengthTooLong)
	}

	// Convert the message to a io.reader object
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(c.message)

	// Send the message
	for _, Url := range c.webhooks {

		resp, err := http.Post(Url, "application/json", buffer)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Check for 204
		if resp.StatusCode != 204 {
			errMsg, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			return fmt.Errorf(string(errMsg))
		}
	}

	return nil
}

type DiscordEmbedClient struct {
	webhooks []string
	username string

	embed domain.DiscordEmbed
}

// Generates a new Discord Embed Client
func NewDiscordEmbedMessage(config domain.ConfigAlertDiscord) DiscordEmbedClient {
	c := DiscordEmbedClient{
		webhooks: config.Webhooks,
		username: config.Username,
	}

	return c
}

// A message can have multiple embeds, Generate a blank one
func (c *DiscordEmbedClient) SetAuthor(Name, URL, IconUrl string) {
	c.embed.Author.Name = Name
	c.embed.Author.IconUrl = IconUrl
	c.embed.Author.Url = URL
}

type DiscordEmbedBodyParams struct {
	Title       string
	Description string
	Color       int
	TitleURL    string
}

func (c *DiscordEmbedClient) SetBody(params DiscordEmbedBodyParams) {
	c.embed.Title = params.Title
	c.embed.Url = params.TitleURL
	c.embed.Description = params.Description
	c.embed.Color = params.Color
}

type DiscordEmbedFieldParams struct {
	Name   string
	Value  string
	Inline bool
}

func (c *DiscordEmbedClient) AppendFields(params DiscordEmbedFieldParams) {
	f := domain.DiscordField{
		Name:   params.Name,
		Value:  params.Value,
		Inline: params.Inline,
	}
	c.embed.Fields = append(c.embed.Fields, f)
}

type DiscordEmbedFooterParams struct {
	Text    string
	IconUrl string
	//TimeStamp string
}

func (c *DiscordEmbedClient) SetFooter(params DiscordEmbedFooterParams) {
	c.embed.Footer.Value = params.Text
	c.embed.Footer.IconUrl = params.IconUrl
}

func (c DiscordEmbedClient) SendPayload() error {
	var embeds []domain.DiscordEmbed

	err := c.validateWebhooks()
	if err != nil {
		return err
	}

	if len(c.embed.Description) >= 4096 {
		return errors.New(ErrDiscordEmbedLenthTooLong)
	}

	embeds = append(embeds, c.embed)
	msg := domain.DiscordBasicMessage{
		Username: c.username,
		Content:  "",
		Embeds:   embeds,
	}

	// Convert the message to a io.reader object
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(msg)

	log.Print(buffer)

	// Send the message
	for _, Url := range c.webhooks {

		resp, err := http.Post(Url, "application/json", buffer)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Check for 204
		if resp.StatusCode != 204 {
			errMsg, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			return fmt.Errorf(string(errMsg))
		}
	}

	return nil
}

func (c DiscordEmbedClient) validateWebhooks() error {
	if len(c.webhooks) == 0 {
		return errors.New("no webhooks given to post to")
	}

	// Check to confirm the webhook uri is correct
	for _, uri := range c.webhooks {
		if strings.Contains(uri, "https://discord.com/api/webhooks/") {
			continue
		}

		return errors.New("invalid uri given to post to")
	}

	return nil
}
