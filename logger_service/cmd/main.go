package main

import (
	"log"
	"net"
	"os"

	"backend_institutions/internal/loggerpb"
	"backend_institutions/logger_service/internals/repository"
	"backend_institutions/logger_service/internals/services"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	loggerRepo := repository.NewLoggerRepo()
	loggerService := services.NewLoggerService(loggerRepo)

	grpcServer := grpc.NewServer()
	loggerpb.RegisterLoggerServiceServer(grpcServer, loggerService)

	port := os.Getenv("LOGGER_GRPC_PORT")
	if port == "" {
		port = "15051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port %s: %v", port, err)
	}

	log.Printf("gRPC Logger Server starting on :%s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}