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

    notificationService := service.NewNotificationService(emailRepo)

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
        log.Fatal(err)
    }

    log.Printf("Notification gRPC Server started on %s", port)

    grpcServer.Serve(lis)
}