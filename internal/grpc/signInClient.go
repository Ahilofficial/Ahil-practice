package grpc

import (
	"backend_institutions/EmailSender/notificationpb"
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	SignInClient     notificationpb.SendMailClient
	SignInConnection *grpc.ClientConn
)

func ConnectSignInService() error {
	port := os.Getenv("SIGNIN_GRPC_PORT")
	if port == "" {
		port = "15053"
	}
	conn, err := grpc.NewClient("localhost:"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	SignInConnection = conn
	SignInClient = notificationpb.NewSendMailClient(conn)
	log.Printf("Connected to SignIn Notification Service on port %s\n", port)
	return nil
}

func SendSignInEmail(email string, subject string, body string) error {
	_, err := SignInClient.SendMail(
		context.Background(),
		&notificationpb.MailRequest{
			To:      email,
			Subject: subject,
			Body:    body,
		},
	)
	return err
}