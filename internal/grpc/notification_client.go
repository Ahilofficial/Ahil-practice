package grpc

import (
	"context"
	"log"
	"os"

	"backend_institutions/internal/notificationpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	NotificationClient notificationpb.NotificationServiceClient
	NotificationConn   *grpc.ClientConn
)

func ConnectNotificationService() error {
	port := os.Getenv("NOTIFICATION_GRPC_PORT")
	if port == "" {
		port = "15052"
	}

	conn, err := grpc.NewClient(
		"localhost:"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	NotificationConn = conn
	NotificationClient = notificationpb.NewNotificationServiceClient(conn)

	log.Printf("Connected to Notification Service on port %s\n", port)
	return nil
}

func SendEmail(to, subject, body string) error {
	_, err := NotificationClient.SendEmail(
		context.Background(),
		&notificationpb.EmailRequest{
			To:      to,
			Subject: subject,
			Body:    body,
		},
	)
	return err
}
