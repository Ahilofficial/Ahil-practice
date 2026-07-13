package seeds

import (
	"backend_institutions/internal/constants"
	"backend_institutions/internal/database"
	"log"
)

func RunSeeders() {
	log.Println("Running database seeders...")
	seedPermissions()

	log.Println("Database seeding completed.")
}

var group = constants.PermissionGroups

func seedPermissions() {
	db, err := database.DB.DB()
	if err != nil {
		log.Printf("Failed to get database connection: %v", err)
		return
	}
	for _, permissions := range group {
		for _, pName := range permissions {
			var count int
			err := db.QueryRow("SELECT COUNT(*) FROM permissions WHERE name = ?", pName).Scan(&count)
			if err != nil {
				log.Printf("Failed to check permission %s: %v", pName, err)
				continue
			}
			if count == 0 {
				_, err = db.Exec("INSERT INTO permissions (name) VALUES (?)", pName)
				if err != nil {
					log.Printf("Failed to insert permission %s: %v", pName, err)
				}
			}
		}
	}
}