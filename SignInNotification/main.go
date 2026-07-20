package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"backend_institutions/EmailSender/notificationpb"
	"backend_institutions/SignInNotification/repository"
	"backend_institutions/SignInNotification/service"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	
	err := godotenv.Load("../.env")
	if err!=nil{
		fmt.Println("Error in loading the envirornment variable")
	}

	emailRepo := repository.NewEmailRepository()
	notificationService := service.NewNotificationService(emailRepo)

	grpcServer := grpc.NewServer()
	notificationpb.RegisterSendMailServer(grpcServer, notificationService)

	port := os.Getenv("SIGNIN_GRPC_PORT")
	if port == "" {
		port = "15053"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("SignIn Notification gRPC Server started on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
