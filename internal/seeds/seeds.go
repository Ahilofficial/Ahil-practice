package seeds

import (
	"backend_institutions/internal/constants"
	"backend_institutions/internal/database"
	"backend_institutions/internal/utils"
	"log"
)

func RunSeeders() {
	log.Println("Running database seeders...")
	seedPermissions()
	seedSuperAdmin()
	log.Println("Database seeding completed.")
}

func seedPermissions() {
	db, err := database.DB.DB()
	if err != nil {
		log.Printf("Failed to get database connection: %v", err)
		return
	}

	for _, pName := range constants.AllPermissions {
		var count int

		err := db.QueryRow(
			"SELECT COUNT(*) FROM permissions WHERE name = ?",
			pName,
		).Scan(&count)

		if err != nil {
			log.Printf("Failed to check permission %s: %v", pName, err)
			continue
		}

		if count == 0 {
			_, err = db.Exec(
				"INSERT INTO permissions (name) VALUES (?)",
				pName,
			)

			if err != nil {
				log.Printf("Failed to insert permission %s: %v", pName, err)
			}
		}
	}
}
func seedSuperAdmin() {
	hashed, err := utils.HashPassword("123456")
	if err != nil {
		panic(err)
	}

	// Create role if it doesn't exist
	err = database.DB.Exec(`
		INSERT INTO roles (name, created_at, updated_at)
		SELECT ?, NOW(), NOW()
		WHERE NOT EXISTS (
			SELECT 1 FROM roles WHERE name = ?
		)
	`, "superadmin", "superadmin").Error
	if err != nil {
		panic(err)
	}

	// Create user if it doesn't exist
	err = database.DB.Exec(`
		INSERT INTO users (
			name,
			email,
			password,
			is_active,
			is_verified,
			created_at,
			updated_at
		)
		SELECT ?, ?, ?, ?, ?, NOW(), NOW()
		WHERE NOT EXISTS (
			SELECT 1 FROM users WHERE email = ?
		)
	`,
		"Super Admin",
		"ahilcicillin@gmail.com",
		hashed,
		true,
		true,
		"ahilcicillin@gmail.com",
	).Error
	if err != nil {
		panic(err)
	}

	// Assign role if it doesn't exist
	err = database.DB.Exec(`
		INSERT INTO user_roles (user_id, role_id)
		SELECT u.id, r.id
		FROM users u
		JOIN roles r ON r.name = ?
		WHERE u.email = ?
		AND NOT EXISTS (
			SELECT 1
			FROM user_roles ur
			WHERE ur.user_id = u.id
			  AND ur.role_id = r.id
		)
	`,
		"superadmin",
		"ahilcicillin@gmail.com",
	).Error
	if err != nil {
		panic(err)
	}
}