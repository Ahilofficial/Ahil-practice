package main

import (
	"log"
	"net"
	"os"

	"notification_service/proto/notificationpb"
	"notification_service/repository"
	"notification_service/service"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	
	err := godotenv.Load()
	if err != nil {
		err = godotenv.Load("../.env")
		if err != nil {
			log.Printf("Warning: could not load .env file from current or parent directory: %v", err)
		}
	}
	

	port := os.Getenv("NOTIFICATION_GRPC_PORT")
	if port == "" {
		port = os.Getenv("NOTIFICATION_PORT")
	}
	if port == "" {
		port = "15052"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	loggerRepo := repository.NewLoggerRepo()
	notificationService := service.NewNotificationService(loggerRepo)
	notificationpb.RegisterNotificationServiceServer(grpcServer, notificationService)

	log.Printf("Notification gRPC Server starting on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
