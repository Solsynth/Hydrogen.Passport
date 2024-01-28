package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"net/textproto"

	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
)

func SendMail(target string, subject string, content string) error {
	mail := &email.Email{
		To:      []string{target},
		From:    viper.GetString("mailer.name"),
		Subject: subject,
		Text:    []byte(content),
		Headers: textproto.MIMEHeader{},
	}
	return mail.SendWithTLS(
		fmt.Sprintf("%s:%d", viper.GetString("mailer.smtp_host"), viper.GetInt("mailer.smtp_port")),
		smtp.PlainAuth(
			"",
			viper.GetString("mailer.username"),
			viper.GetString("mailer.password"),
			viper.GetString("mailer.smtp_host"),
		),
		&tls.Config{ServerName: viper.GetString("mailer.smtp_host")},
	)
}

func SendMailHTML(target string, subject string, content string) error {
	mail := &email.Email{
		To:      []string{target},
		From:    viper.GetString("mailer.name"),
		Subject: subject,
		HTML:    []byte(content),
		Headers: textproto.MIMEHeader{},
	}
	return mail.SendWithTLS(
		fmt.Sprintf("%s:%d", viper.GetString("mailer.smtp_host"), viper.GetInt("mailer.smtp_port")),
		smtp.PlainAuth(
			"",
			viper.GetString("mailer.username"),
			viper.GetString("mailer.password"),
			viper.GetString("mailer.smtp_host"),
		),
		&tls.Config{ServerName: viper.GetString("mailer.smtp_host")},
	)
}
