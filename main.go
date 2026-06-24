package main

import (
	"backend_institutions/database"
	"backend_institutions/model"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()
	database.Connect()
	err := database.DB.AutoMigrate(
		&model.Institutions{},
		&model.Department{},
		&model.Faculty{},
		&model.Student{},
		&model.Fees{},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(app.Listen(":8000"))
}
