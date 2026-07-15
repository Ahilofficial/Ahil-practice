package middleware

import (
	"backend_institutions/internal/grpc"
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func RequestResponseLogger() fiber.Handler {
	return func(c fiber.Ctx) error {
		
		method := c.Method()
		endpoint := c.Path()
		reqBody := string(c.Body())

		
		err := c.Next()

		
		status := c.Response().StatusCode()
		respBody := string(c.Response().Body())

		
		serviceName := "General"
		cleanEndpoint := strings.ToLower(endpoint)
		if strings.HasPrefix(cleanEndpoint, "/institutes") {
			serviceName = "Institution"
		} else if strings.HasPrefix(cleanEndpoint, "/departments") {
			serviceName = "Department"
		} else if strings.HasPrefix(cleanEndpoint, "/faculties") {
			serviceName = "Faculty"
		} else if strings.HasPrefix(cleanEndpoint, "/students") {
			serviceName = "Student"
		} else if strings.HasPrefix(cleanEndpoint, "/fees") {
			serviceName = "Fees"
		} else if strings.HasPrefix(cleanEndpoint, "/signup") || strings.HasPrefix(cleanEndpoint, "/signin") {
			serviceName = "Auth"
			reqBody = "[REDACTED/SENSITIVE]"
		}

		go func() {
			sendLogErr := grpc.SendLog(
				serviceName,
				method,
				endpoint,
				reqBody,
				respBody,
				int32(status),
			)
			if sendLogErr != nil {
				log.Printf("Error sending request/response log to gRPC service: %v\n", sendLogErr)
			}
		}()

		return err
	}
}
