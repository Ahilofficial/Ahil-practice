package main

import (
	"backend_institutions/logger_service/internals/grpc"
	"backend_institutions/logger_service/internals/repository"
	// "backend_institutions/logger_service/internals/routes"
	"backend_institutions/logger_service/internals/services"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	
	loggerRepo := repository.NewLoggerRepo()
	loggerService := services.NewLoggerService(loggerRepo)

	
	go grpc.StartGRPCServer(loggerService, "50051")

	
	app := fiber.New()
	// routes.RegisterRoutes(app)
	log.Println("Logger HTTP server starting on :8081")
	log.Fatal(app.Listen(":8081"))
}