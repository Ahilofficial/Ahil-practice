package middleware

import (
	"backend_institutions/internal/database"
	"backend_institutions/internal/helper"
	"backend_institutions/internal/model"

	"github.com/gofiber/fiber/v3"
)

func RequirePermission(permission string) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := c.Locals("user_id")
		if userID == nil {
			return helper.Error(c, 401, "Unauthorized")
		}

		var user model.User
		err := database.DB.
			Preload("Roles.Permissions").
			First(&user, userID).Error

		if err != nil {
			return helper.Error(c, 403, "Forbidden")
		}

		for _, role := range user.Roles {
			if role.Name == "admin" {
				return c.Next()
			}
			for _, perm := range role.Permissions {
				if perm.Name == permission {
					return c.Next()
				}
			}
		}

		return helper.Error(c, 403, "Forbidden: you do not have permission to perform this action")
	}
}
