package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", host, port),
		auth,
		from,
		[]string{to},
		[]byte(message),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
