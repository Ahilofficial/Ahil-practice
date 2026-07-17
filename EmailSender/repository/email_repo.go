package repository

import "backend_institutions/EmailSender/smtp"

type EmailRepository struct{}

func NewEmailRepository() *EmailRepository {
	return &EmailRepository{}
}

func (r *EmailRepository) SendMail(email, subject, body string) error {
	return smtp.SendEmail(email, subject, body)
}