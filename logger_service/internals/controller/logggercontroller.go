package controller

import (
	"backend_institutions/logger_service/internals/model"
	"backend_institutions/logger_service/internals/services"

	"github.com/gofiber/fiber/v3"
)

type LoggerController struct {
	loggerService *services.LoggerService
}

func NewLoggerController(loggerService *services.LoggerService) *LoggerController {
	return &LoggerController{
		loggerService: loggerService,
	}
}

func (c *LoggerController) WriteLog(ctx fiber.Ctx) error {

	var log model.Log

	// Read JSON request body
	if err := ctx.Bind().Body(&log); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Save log
	if err := c.loggerService.SaveLog(log); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to write log",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Log written successfully",
	})
}