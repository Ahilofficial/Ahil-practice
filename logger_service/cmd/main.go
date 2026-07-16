package main

import (
	"backend_institutions/logger_service/internals/grpc"
	"backend_institutions/logger_service/internals/repository"
	// "backend_institutions/logger_service/internals/routes"
	"backend_institutions/logger_service/internals/services"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
)

func main() {
	
	loggerRepo := repository.NewLoggerRepo()
	loggerService := services.NewLoggerService(loggerRepo)

	port := os.Getenv("LOGGER_GRPC_PORT")
	if port == "" {
		port = "15051"
	}
	
	go grpc.StartGRPCServer(loggerService, port)

	
	app := fiber.New()
	// routes.RegisterRoutes(app)
	log.Println("Logger HTTP server starting on :8081")
	log.Fatal(app.Listen(":8081"))
}