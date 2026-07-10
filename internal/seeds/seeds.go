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

var institute = constants.InstitutePermissions
var department = constants.DepartmentPermissions
var faculty = constants.FacultyPermissions
var student = constants.StudentPermissions
var other = constants.OtherPermissions

func seedPermissions() {

	permissionsGroup := [][]string{

		institute,
		department,
		faculty,
		student,
		other,
	}
	for _, group := range permissionsGroup {
		for _, pName := range group {
			var perm model.Permission
			err := database.DB.Where("name = ?", pName).First(&perm).Error
			if err != nil {
				perm = model.Permission{
					Name: pName,
				}
				if createErr := database.DB.Create(&perm).Error; createErr != nil {
					log.Printf("Failed to seed permission %s: %v", pName, createErr)
				}
			}
		}
	}
}

// 	for _, pName := range permissions {
// 		var perm model.Permission
// 		err := database.DB.Where("name = ?", pName).First(&perm).Error
// 		if err != nil {
// 			perm = model.Permission{
// 				Name: pName,
// 			}
// 			if createErr := database.DB.Create(&perm).Error; createErr != nil {
// 				log.Printf("Failed to seed permission %s: %v", pName, createErr)
// 			}
// 		}
// 	}
// }
