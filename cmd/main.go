package main

import (
	"backend_institutions/internal/database"
	"backend_institutions/internal/model"
	"backend_institutions/internal/seeds"
	"log"
	"os"

	"backend_institutions/internal/wire"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, relying on environment variables or defaults")
	}

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

	router, err := wire.InitializeRouter()
	if err != nil {
		log.Fatal("Failed to initialize dependency injection: ", err)
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8090"
	}
	log.Printf("Server starting on :%s", port)

	log.Fatal(router.App.Listen(":" + port))
}
