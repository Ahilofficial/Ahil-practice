package main

import (
	"backend_institutions/logger_service/internals/routes"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()
	routes.RegisterRoutes(app)
	log.Fatal(app.Listen(":8081"))
}
