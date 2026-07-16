package smtp

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	Host     string
	Port     string
	Email    string
	Password string
}

func NewMailer(host, port, email, password string) *Mailer {
	return &Mailer{
		Host:     host,
		Port:     port,
		Email:    email,
		Password: password,
	}
}

func (m *Mailer) SendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
	auth := smtp.PlainAuth("", m.Email, m.Password, m.Host)

	// Construct SMTP headers and body
	header := make(map[string]string)
	header["From"] = m.Email
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	err := smtp.SendMail(addr, auth, m.Email, []string{to}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
