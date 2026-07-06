package seeds

import (
	"backend_institutions/database"
	"backend_institutions/model"
	"log"
)

func RunSeeders() {
	log.Println("Running database seeders...")
	seedPermissions()
	seedRoles()
	assignPermissionsToRoles()
	log.Println("Database seeding completed.")
}

func seedPermissions() {
	permissions := []string{
		"create_institutes", "view_institutes", "update_institutes", "delete_institutes",
		"create_departments", "view_departments", "update_departments", "delete_departments",
		"create_faculties", "view_faculties", "update_faculties", "delete_faculties",
		"create_students", "view_students", "update_students", "delete_students",
		"create_fees", "view_fees", "update_fees", "delete_fees",
		"assign_roles",
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

func seedRoles() {
	roles := []string{"admin", "principal", "faculty", "student", "user"}
	for _, rName := range roles {
		var role model.Role
		err := database.DB.Where("name = ?", rName).First(&role).Error 
		if err != nil {
			role = model.Role{Name: rName}
			if createErr := database.DB.Create(&role).Error; createErr != nil {
				log.Printf("Failed to seed role %s: %v", rName, createErr)
			}
		}
	}
}

func assignPermissionsToRoles() {
	principalPerms := []string{
		"create_departments", "view_departments", "update_departments", "delete_departments",
		"create_faculties", "view_faculties", "update_faculties", "delete_faculties",
		"create_students", "view_students", "update_students", "delete_students",
	}

	facultyPerms := []string{
		"create_faculties", "view_faculties", "update_faculties", "delete_faculties",
		"create_students", "view_students", "update_students", "delete_students",
	}

	studentPerms := []string{
		"create_students", "view_students", "update_students", "delete_students",
		"create_fees", "view_fees", "update_fees", "delete_fees",
	}

	allPerms := []string{
		"create_institutes", "view_institutes", "update_institutes", "delete_institutes",
		"create_departments", "view_departments", "update_departments", "delete_departments",
		"create_faculties", "view_faculties", "update_faculties", "delete_faculties",
		"create_students", "view_students", "update_students", "delete_students",
		"create_fees", "view_fees", "update_fees", "delete_fees",
		"assign_roles",
	}
	
	mapRolePerms("principal", principalPerms)
	mapRolePerms("faculty", facultyPerms)
	mapRolePerms("student", studentPerms)
	mapRolePerms("admin", allPerms)
	mapRolePerms("user", []string{})
}

func mapRolePerms(roleName string, permNames []string) {
	var role model.Role
	if err := database.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		log.Printf("Role %s not found for mapping: %v", roleName, err)
		return
	}

	var perms []model.Permission
	if err := database.DB.Where("name IN ?", permNames).Find(&perms).Error; err != nil {
		log.Printf("Failed to retrieve permissions for role %s: %v", roleName, err)
		return
	}

	if err := database.DB.Model(&role).Association("Permissions").Replace(perms); err != nil {
		log.Printf("Failed to map permissions for role %s: %v", roleName, err)
	}
}