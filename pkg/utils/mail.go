package utils

import (
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"gopkg.in/gomail.v2"
)

var mailers chan *gomail.Dialer

func InitMailer() {
	if len(Config.SmtpCreds) == 0 {
		logger.Errorf("No smtp creds found")
		return
	}

	mailers = make(chan *gomail.Dialer, len(Config.SmtpCreds))

	for _, creds := range Config.SmtpCreds {
		dialer := gomail.NewDialer(Config.EmailHost, Config.EmailPort, creds.User, creds.Password)
		mailers <- dialer
		logger.Infof("Created dialer with: %s", creds.User)
	}
	logger.Infof("Successfully initialized dialers")
}

func SendEmail(to, subject, body string, attachments ...string) error {

	dialer := <-mailers
	defer func() {
		mailers <- dialer
	}()
	m := gomail.NewMessage()
	m.SetHeader("From", Config.SendingEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	for _, attachment := range attachments {
		if attachment != "" {
			m.Attach(attachment)
		}
	}

	if err := dialer.DialAndSend(m); err != nil {
		logger.Errorf("Unable to send mail: %v", err)
		return err
	}
	logger.Infof("Sent mail using: " + dialer.Username)
	return nil
}
