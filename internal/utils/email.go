package utils

import (
	"github.com/sirupsen/logrus"
	"go_project/internal/config"
	"gopkg.in/gomail.v2"
)

// SendEmail отправляет письмо с заданными темой и текстом на указанный адрес.
func SendEmail(to string, subject string, body string) error {
	cfg, err := config.LoadConfig("config")
	if err != nil {
		logrus.Fatal("cannot load config:", err)
	}
	var smtpHost = cfg.SMTP.Host
	var smtpPort = cfg.SMTP.Port
	var smtpUser = cfg.SMTP.User
	var smtpPass = cfg.SMTP.Pass
	var smtpFrom = cfg.SMTP.From
	m := gomail.NewMessage()
	m.SetHeader("From", smtpFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	return d.DialAndSend(m)
}
