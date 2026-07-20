package repository

import (
	"backend_institutions/SignInNotification/smtp"
	"backend_institutions/utilities"
)

type EmailRepository struct{}

func NewEmailRepository() *EmailRepository {
	return &EmailRepository{}
}

func (r *EmailRepository) SendMail(email, subject, body string) error {
	err:= smtp.SendSignInEmail(email, subject, body)
	if err != nil {
		_ = utilities.WriteEmailLog(email, subject, false, err.Error())
		return err
	}

	_ = utilities.WriteEmailLog(email, subject, true, "")
	return  nil
}
