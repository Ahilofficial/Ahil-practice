package service

import (
	"context"
	"fmt"
	"os"

	"notification_service/proto/notificationpb"
	"notification_service/repository"
	"notification_service/smtp"
)

type NotificationService struct {
	notificationpb.UnimplementedNotificationServiceServer
	loggerRepo *repository.LoggerRepo
}

func NewNotificationService(loggerRepo *repository.LoggerRepo) *NotificationService {
	return &NotificationService{loggerRepo: loggerRepo}
}

func (s *NotificationService) SendEmail(ctx context.Context, req *notificationpb.EmailRequest) (*notificationpb.EmailResponse, error) {
	// Read SMTP configuration from environment
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	email := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	if host == "" || port == "" || email == "" || password == "" {
		err := fmt.Errorf("SMTP configuration variables (SMTP_HOST, SMTP_PORT, SMTP_EMAIL, SMTP_PASSWORD) are not fully configured")
		_ = s.loggerRepo.WriteLog(req.To, req.Subject, false, err.Error())
		return &notificationpb.EmailResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	mailer := smtp.NewMailer(host,port,email, password)
	err := mailer.SendEmail(req.To, req.Subject, req.Body)

	if err != nil {
		_ = s.loggerRepo.WriteLog(req.To, req.Subject, false, err.Error())
		return &notificationpb.EmailResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to send email: %v", err),
		}, nil
	}

	_ = s.loggerRepo.WriteLog(req.To, req.Subject, true, "")
	return &notificationpb.EmailResponse{
		Success: true,
		Message: "Email sent successfully",
	}, nil
}
