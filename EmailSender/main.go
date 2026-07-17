package main

import (
	"log"
	"net"
	"os"

	"backend_institutions/EmailSender/notificationpb"
	"backend_institutions/EmailSender/repository"
	"backend_institutions/EmailSender/service"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	emailRepo := repository.NewEmailRepository()

	emailService := service.NewEmailService(emailRepo)

	
	notificationService := service.NewNotificationService(emailService)
	

	grpcServer := grpc.NewServer()

	notificationpb.RegisterSendMailServer(
		grpcServer,
		notificationService,
	)

	port := os.Getenv("NOTIFICATION_GRPC_PORT")
	if port == "" {
		port = "15052"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Notification gRPC Server started on port %s", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}