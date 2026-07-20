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
	for roleName, permissions := range group {

		_, err := db.Exec(`
		INSERT INTO roles(name)
		SELECT ?
		WHERE NOT EXISTS (
			SELECT 1 FROM roles WHERE name = ?
		)
	`, roleName, roleName)

		if err != nil {
			log.Println(err)
			continue
		}

		var roleID int
		db.QueryRow("SELECT id FROM roles WHERE name = ?", roleName).Scan(&roleID)

		for _, permName := range permissions {

			_, err := db.Exec(`
			INSERT INTO permissions(name)
			SELECT ?
			WHERE NOT EXISTS (
				SELECT 1 FROM permissions WHERE name = ?
			)
		`, permName, permName)

			if err != nil {
				log.Println(err)
				continue
			}

			var permID int
			db.QueryRow("SELECT id FROM permissions WHERE name = ?", permName).Scan(&permID)

			_, err = db.Exec(`
			INSERT INTO role_permissions(role_id, permission_id)
			SELECT ?, ?
			WHERE NOT EXISTS (
				SELECT 1
				FROM role_permissions
				WHERE role_id = ?
				AND permission_id = ?
			)
		`, roleID, permID, roleID, permID)

			if err != nil {
				log.Println(err)
			}
		}
	}
}
