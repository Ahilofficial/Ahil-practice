package main

import (
	"backend_institutions/database"
	"backend_institutions/model"
	"backend_institutions/routes"
	"backend_institutions/seeds"
	"log"
	"os"
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

	routes.SetUpRoutes(app)
	
	port:=os.Getenv("APP_PORT")
	if port==""{
		port="8090"
	}
	log.Printf("Server starting on :%s", port)

	log.Fatal(app.Listen(":" + port))
}
