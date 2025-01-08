package utils

import (
	logger "github.com/CodeChefVIT/devsoc-be-24/pkg/logging"
	"gopkg.in/gomail.v2"
)

type mails struct {
	Dialer *gomail.Dialer
	Conn   gomail.SendCloser
}

var mailers chan mails

func InitMailer() {
	mailers = make(chan mails, len(Config.SmtpCreds))

	for _, creds := range Config.SmtpCreds {
		dialer := gomail.NewDialer(Config.EmailHost, Config.EmailPort, creds.User, creds.Password)
		mail, err := createConnection(dialer)
		if err != nil {
			logger.Infof("Unable to connect to SMTP server due to error: ", err.Error())
			continue
		}
		mailers <- mail
		logger.Infof("Created mailer with: " + creds.User)
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

func SendEmail(to, subject, body string) error {
	mail := <-mailers

	m := gomail.NewMessage()
	m.SetHeader("From", Config.SendingEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := gomail.Send(mail.Conn, m); err != nil {
		logger.Errorf("Unable to send mail, trying to reconnect to server: ", err.Error())
		_ = mail.Conn.Close()
		mail, err = createConnection(mail.Dialer)
		if err != nil {
			logger.Infof("Unable to connect, skipping: ", err.Error())
			return err
		}
		err = gomail.Send(mail.Conn, m)
		if err != nil {
			return err
		}
	}
	mailers <- mail
	logger.Infof("Sent mail using: " + Config.SendingEmail)
	return nil
}
