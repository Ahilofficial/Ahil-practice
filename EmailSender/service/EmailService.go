package service

import (
	"context"
	"fmt"

	"backend_institutions/EmailSender/notificationpb"
	"backend_institutions/EmailSender/repository"
	"backend_institutions/utilities"
)

type NotificationService struct {
	notificationpb.UnimplementedSendMailServer
	repo *repository.EmailRepository
}

func NewNotificationService(repo *repository.EmailRepository) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

func (s *NotificationService) SendMail(
	ctx context.Context,
	req *notificationpb.MailRequest,
) (*notificationpb.MailResponse, error) {

	var subject, body string
	name := req.Subject
	emailType := req.Body

	switch emailType {
	case "signup":
		subject = "Welcome to Backend Institutions!"
		body = fmt.Sprintf("<h1>Hello %s,</h1><p>Thank you for registering on our platform. Your account is now active!</p>", name)

	case "signin":
		subject = "Sign In Notification"
		body = fmt.Sprintf("<h1>Hello %s,</h1><p>You have successfully signed in to our platform.</p>", name)

	
	}

	err := s.repo.SendMail(req.To, subject, body)

	if err != nil {
		_ = utilities.WriteEmailLog(req.To, subject, false, err.Error())
		return nil, err
	}

	_ = utilities.WriteEmailLog(req.To, subject, true, "")

	return &notificationpb.MailResponse{
		Success: true,
		Message: "Mail sent successfully",
	}, nil
}
