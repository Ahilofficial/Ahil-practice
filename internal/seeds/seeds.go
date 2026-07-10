package seeds

import (
	"backend_institutions/internal/constants"
	"backend_institutions/internal/database"
	"backend_institutions/internal/model"
	"log"
)

func RunSeeders() {
	log.Println("Running database seeders...")
	seedPermissions()

	log.Println("Database seeding completed.")
}

var group= constants.PermissionGroups

func seedPermissions() {
	for _, permissions := range group {
		for _, pName := range permissions {
			var perm model.Permission

			err := database.DB.Where("name = ?", pName).First(&perm).Error
			if err != nil {
				perm = model.Permission{Name: pName}
				database.DB.Create(&perm)
			}
		}
	}
}