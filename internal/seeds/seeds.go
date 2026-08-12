package seeds

import (
	"backend_institutions/internal/constants"
	"backend_institutions/internal/database"
	// "backend_institutions/internal/model"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func RunSeeders() {
	log.Println("Running database seeders...")

	seedPermissions()
	seedSuperAdminRole()
	seedSuperAdminUser()
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

func seedSuperAdminRole() {
	db, err := database.DB.DB()
	if err != nil {
		log.Printf("Failed to get database connection for super admin seed: %v", err)
		return
	}

	var roleID uint
	err = db.QueryRow("SELECT id FROM roles WHERE LOWER(name) IN ('super admin', 'super_admin', 'superadmin') LIMIT 1").Scan(&roleID)
	if err != nil || roleID == 0 {
		res, err := db.Exec("INSERT INTO roles (name, created_at, updated_at) VALUES ('Super Admin', NOW(), NOW())")
		if err != nil {
			log.Printf("Failed to insert Super Admin role: %v", err)
			return
		}
		id, _ := res.LastInsertId()
		roleID = uint(id)
		log.Printf("Seeded 'Super Admin' role with ID %d", roleID)
	}

	for _, pName := range constants.AllPermissions {
		var permID uint
		_ = db.QueryRow("SELECT id FROM permissions WHERE name = ?", pName).Scan(&permID)
		if permID > 0 {
			var count int
			_ = db.QueryRow("SELECT COUNT(*) FROM role_permissions WHERE role_id = ? AND permission_id = ?", roleID, permID).Scan(&count)
			if count == 0 {
				_, _ = db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", roleID, permID)
			}
		}
	}
}

func seedSuperAdminUser() {
	db, err := database.DB.DB()
	if err != nil {
		log.Printf("Failed to get database connection for super admin user seed: %v", err)
		return
	}

	superEmail := "ahilcicillin@gmail.com"
	superPassword := "123456"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(superPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash super admin password: %v", err)
		return
	}

	var userID uint
	err = db.QueryRow("SELECT id FROM users WHERE email = ? LIMIT 1", superEmail).Scan(&userID)
	if err != nil || userID == 0 {
		res, err := db.Exec(
			"INSERT INTO users (name, email, phone, password, is_active, is_verified, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())",
			"Super Admin",
			superEmail,
			"0000000000",
			string(hashedPassword),
			true,
			true,
		)
		if err != nil {
			log.Printf("Failed to insert Super Admin user: %v", err)
			return
		}
		id, _ := res.LastInsertId()
		userID = uint(id)
		log.Printf("Seeded Super Admin user '%s' with ID %d", superEmail, userID)
	} else {
		_, err = db.Exec(
			"UPDATE users SET password = ?, is_active = true, is_verified = true, updated_at = NOW() WHERE id = ?",
			string(hashedPassword),
			userID,
		)
		if err != nil {
			log.Printf("Failed to update Super Admin user credentials: %v", err)
		}
	}

	var roleID uint
	_ = db.QueryRow("SELECT id FROM roles WHERE LOWER(name) IN ('super admin', 'super_admin', 'superadmin') LIMIT 1").Scan(&roleID)

	if roleID > 0 && userID > 0 {
		var count int
		_ = db.QueryRow("SELECT COUNT(*) FROM user_roles WHERE user_id = ? AND role_id = ?", userID, roleID).Scan(&count)
		if count == 0 {
			_, err = db.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, roleID)
			if err != nil {
				log.Printf("Failed to assign Super Admin role to user: %v", err)
			} else {
				log.Printf("Assigned Super Admin role to user '%s'", superEmail)
			}
		}
	}
}
