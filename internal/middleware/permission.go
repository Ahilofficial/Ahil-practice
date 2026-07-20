package middleware

import (
	
	"backend_institutions/internal/database"
	"backend_institutions/internal/helper"

	"github.com/gofiber/fiber/v3"
)

func RequirePermission(permission string) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := c.Locals("user_id")
		if userID == nil {
			return helper.Error(c, 401, "Unauthorized")
		}

		var count int64
		const query = `
			SELECT COUNT(*)
			FROM user_roles ur
			JOIN roles r ON ur.role_id = r.id
			LEFT JOIN role_permissions rp ON rp.role_id = r.id
			LEFT JOIN permissions p ON rp.permission_id = p.id
			WHERE ur.user_id = ?
			  AND (r.name = ? OR p.name = ?)
		`

		err := database.DB.
			Raw(query, userID, "admin", permission).
			Scan(&count).Error

		if err != nil || count == 0 {
			return helper.Error(c, 403, "Forbidden: you do not have permission to perform this action")
		}

		return c.Next()
	}
}