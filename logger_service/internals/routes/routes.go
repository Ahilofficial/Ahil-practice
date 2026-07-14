package routes

import (
	"backend_institutions/logger_service/internals/controller"
	"backend_institutions/logger_service/internals/repository"
	"backend_institutions/logger_service/internals/services"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app *fiber.App) {

	// Dependency Injection
	loggerRepo := repository.NewLoggerRepo()
	loggerService := services.NewLoggerService(loggerRepo)
	loggerController := controller.NewLoggerController(loggerService)

	// Routes
	app.Post("/logs", loggerController.WriteLog)
}