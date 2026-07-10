package seeds

import (
	"backend_institutions/internal/database"
	"backend_institutions/internal/model"
	"log"
	"backend_institutions/internal/constants"
)

func RunSeeders() {
	log.Println("Running database seeders...")
	seedPermissions()
	
	log.Println("Database seeding completed.")
}



func seedPermissions() {
	permissions := []string{
		constants.PermissionCreateInstitutes, constants.PermissionViewInstitutes, constants.PermissionUpdateInstitutes, constants.PermissionDeleteInstitutes,
		constants.PermissionCreateDepartments, constants.PermissionViewDepartments, constants.PermissionUpdateDepartments, constants.PermissionDeleteDepartments,
		constants.PermissionCreateFaculties, constants.PermissionViewFaculties, constants.PermissionUpdateFaculties, constants.PermissionDeleteFaculties,
		constants.PermissionCreateStudents, constants.PermissionViewStudents, constants.PermissionUpdateStudents, constants.PermissionDeleteStudents,
		constants.PermissionCreateFees, constants.PermissionViewFees, constants.PermissionUpdateFees, constants.PermissionDeleteFees,
		constants.PermissionAssignRoles,
	}

	for _, pName := range permissions {
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