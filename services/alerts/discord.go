package alerts

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jtom38/dvb/domain"
)

var (
	ErrDiscordContentLengthTooLong = "the message length is greater then 2000 characters"
)

type DiscordAlertClient struct {
	webhooks []string
	username string

	messages []string

	message domain.DiscordMessage
}

func NewDiscordAlertClient(Webhooks []string, Username string) DiscordAlertClient {
	c := DiscordAlertClient{
		webhooks: Webhooks,
		username: Username,
	}
	c.message = c.NewMessage()
	return c
}

func (c DiscordAlertClient) NewMessage() domain.DiscordMessage {
	m := domain.DiscordMessage{
		Username: &c.username,
	}
	return m
}

func (c DiscordAlertClient) validateWebhooks() error {
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

func (c *DiscordAlertClient) SetMessage(Message domain.DiscordMessage) {
	c.message = Message
}

func (c DiscordAlertClient) sendPayload(msg domain.DiscordMessage) error {
	err := c.validateWebhooks()
	if err != nil {
		return err
	}

	if len(*c.message.Content) >= 2000 {
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

func (c *DiscordAlertClient) ReplaceContent(content string) {
	c.message.Content = &content
}

func (c *DiscordAlertClient) UpdateContent(content string) {
	c.messages = append(c.messages, content)
}

func (c DiscordAlertClient) Send() error {
	return c.sendPayload(c.message)
}

func (c DiscordAlertClient) SendMessage(msg domain.DiscordMessage) error {
	return c.sendPayload(msg)
}
