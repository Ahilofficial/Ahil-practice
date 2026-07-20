package service

import (
	"context"

	"backend_institutions/EmailSender/notificationpb"
	"backend_institutions/SignInNotification/repository"
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

	err := s.repo.SendMail(
		req.To,
		req.Subject,
		req.Body,
	)

	if err != nil {
		
		return nil, err
	}

	_ = utilities.WriteEmailLog(req.To, req.Subject, true, "")

	return &notificationpb.MailResponse{
		Message: "Sign-in mail sent successfully",
	}, nil
}
