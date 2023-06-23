package Email

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type MailProcessor struct {
	from     string
	fromName string
	bodyType string
	smtpHost string
	smtpPort int
	smtpUser string
	smtpPass string
}

func NewMailProcessor() *MailProcessor {
	return &MailProcessor{
		from:     viper.GetString("SMTP.From"),
		fromName: viper.GetString("SMTP.Name"),
		bodyType: "text/html",
		smtpHost: viper.GetString("SMTP.Host"),
		smtpPort: viper.GetInt("SMTP.Port"),
		smtpUser: viper.GetString("SMTP.Username"),
		smtpPass: viper.GetString("SMTP.Password"),
	}
}

func (m *MailProcessor) Send(to string, subject string, content string) error {
	message := gomail.NewMessage(gomail.SetCharset("UTF-8"))
	message.SetAddressHeader("From", m.from, m.fromName)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "[Rustdesk]"+subject)
	message.SetBody(m.bodyType, content)
	dialer := gomail.NewDialer(m.smtpHost, m.smtpPort, m.smtpUser, m.smtpPass)
	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("error occoured when sending mail: %w", err)
	}
	return nil
}
