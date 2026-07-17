package service

import "backend_institutions/EmailSender/repository"

type EmailService struct {
	repo *repository.EmailRepository
}

func NewEmailService(repo *repository.EmailRepository) *EmailService {
	return &EmailService{
		repo: repo,
	}
}

