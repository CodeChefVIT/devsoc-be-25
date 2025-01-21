package utils

import (
	"sync"

	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"gopkg.in/gomail.v2"
)

type mails struct {
	Dialer *gomail.Dialer
	Conn   gomail.SendCloser
}

var (
	mailers chan mails
	mu      sync.Mutex
)

func InitMailer() {
	mailers = make(chan mails, len(Config.SmtpCreds))

	for _, creds := range Config.SmtpCreds {
		dialer := gomail.NewDialer(Config.EmailHost, Config.EmailPort, creds.User, creds.Password)
		mail, err := createConnection(dialer)
		if err != nil {
			logger.Infof("Unable to connect to SMTP server due to error: %v", err)
			continue
		}
		mailers <- mail
		logger.Infof("Created mailer with: %s", creds.User)
	}
	logger.Infof("Successfully initialized mailers")
}

func createConnection(dialer *gomail.Dialer) (mails, error) {
	conn, err := dialer.Dial()
	if err != nil {
		return mails{}, err
	}

	return mails{
		Dialer: dialer,
		Conn:   conn,
	}, nil
}

func reinitMailers() {
	mu.Lock()
	defer mu.Unlock()

	for len(mailers) > 0 {
		mail := <-mailers
		mail.Conn.Close()
	}

	logger.Infof("Reinitializing mailers")
	InitMailer()
}

func SendEmail(to, subject, body string, attachments ...string) error {
	var mail mails

	select {
	case mail = <-mailers:
	default:
		logger.Errorf("No mailer available, reinitializing mailers")
		reinitMailers()
		mail = <-mailers
	}

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

	if err := gomail.Send(mail.Conn, m); err != nil {
		logger.Errorf("Unable to send mail: %v", err)
		_ = mail.Conn.Close()
		return err
	}
	mailers <- mail
	logger.Infof("Sent mail using: " + Config.SendingEmail)
	return nil
}
