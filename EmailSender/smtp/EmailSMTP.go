package smtp

import (
	"net/smtp"
	"os"
)

func SendEmail(email string, subject string, body string) error {
	myemail := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	message := "Hello welcome to my application"
	to := email
	address := host + ":" + port
	auth := smtp.PlainAuth("", myemail, password, host)
	err := smtp.SendMail(address, auth, myemail, []string{to}, []byte(message))
	if err != nil {
		return err
	}
	return nil
}
