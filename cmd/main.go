package main

import (
	"backend_institutions/internal/database"
	"backend_institutions/internal/model"
	"backend_institutions/internal/routes"
	"backend_institutions/internal/seeds"
	"log"
	"os"

	"backend_institutions/internal/wire"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, relying on environment variables or defaults")
	}

	app := fiber.New()
	database.Connect()
	err = database.DB.AutoMigrate(
		&model.Institutions{},
		&model.Department{},
		&model.Faculty{},
		&model.Student{},
		&model.Fees{},
		&model.User{},
		&model.Role{},
		&model.Permission{},
	)
	if err != nil {
		log.Fatal(err)
	}

	seeds.RunSeeders()

	router := wire.InitializeRouter()
	if err != nil {
		log.Fatal(err)
	}

	router.Start()

	routes.SetUpRoutes(app)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8090"
	}
	log.Printf("Server starting on :%s", port)

	log.Fatal(app.Listen(":" + port))
}
