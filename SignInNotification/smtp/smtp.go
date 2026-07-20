package smtp

import (
	"net/smtp"
	"os"
)

func SendSignInEmail(email string, subject string, body string) error {
	myemail := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" +
		body + "\r\n")

	address := host + ":" + port
	auth := smtp.PlainAuth("", myemail, password, host)
	return smtp.SendMail(address, auth, myemail, []string{email}, msg)
}
