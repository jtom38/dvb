package alerts

import (
	"crypto/tls"
	"errors"

	"github.com/jtom38/dvb/domain"
	"gopkg.in/gomail.v2"
)

const (
	EmailSubjectSuccess = "Backup Successful"
	EmailSubjectError   = "Backup Error"
)

type SmtpClient struct {
	Config domain.ConfigAlertEmail

	subject string
	body    string
}

func NewSmtpClient(config domain.ConfigAlertEmail) SmtpClient {
	return SmtpClient{
		Config: config,
	}
}

func (c *SmtpClient) SetSubject(value string) {
	c.subject = value
}

func (c *SmtpClient) SetBody(value string) {
	c.body = value
}

func (c SmtpClient) SendAlert() error {

	if c.subject == "" {
		return errors.New("subject was null and needs a value")
	}

	if c.body == "" {
		return errors.New("no body was found")
	}

	dial := gomail.NewDialer(c.Config.Account.Host, c.Config.Account.Port, c.Config.Account.Username, c.Config.Account.Password)

	dial.TLSConfig = &tls.Config{
		InsecureSkipVerify: c.Config.Account.UseTls,
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", c.Config.From)
	msg.SetAddressHeader("To", c.Config.To, "test")
	msg.SetHeader("Subject", c.subject)
	msg.SetBody("text/html", c.body)

	err := dial.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
