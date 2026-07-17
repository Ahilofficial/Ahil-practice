package service

import (
	"context"

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

	err := s.repo.SendMail(
		req.To,
		req.Subject, 
		req.Body,
	)

	if err != nil {
		_ = utilities.WriteEmailLog(req.To, req.Subject, false, err.Error())
		return nil, err
	}
	
	_ = utilities.WriteEmailLog(req.To, req.Subject, true, "")

	return &notificationpb.MailResponse{
		Message: "Mail sent successfully",
	}, nil
}
