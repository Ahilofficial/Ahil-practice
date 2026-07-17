package service

import (
	"context"

	"backend_institutions/EmailSender/notificationpb"
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

	err := s.emailService.SendMail(
		req.To,
		req.Subject, 
		req.Body,
	)

	if err != nil {
		return nil, err
	}

	return &notificationpb.MailResponse{
		Message: "Mail sent successfully",
	}, nil
}
