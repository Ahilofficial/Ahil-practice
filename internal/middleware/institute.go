package middleware

import (
	"backend_institutions/internal/database"
	"backend_institutions/internal/helper"

	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func InstitutionAccess() fiber.Handler {
	return func(c fiber.Ctx) error {

		userID := c.Locals("user_id").(uint)

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return helper.Error(c, 400, "Invalid ID")
		}

		var institutionID uint

		resource := strings.Split(strings.Trim(c.Route().Path, "/"), "/")[0]

		switch resource {

		case "students":
			err = database.DB.Raw(`
				SELECT d.institution_id
				FROM students s
				JOIN faculties f ON s.faculty_id = f.id
				JOIN departments d ON f.department_id = d.id
				WHERE s.id = ?
			`, id).Scan(&institutionID).Error

		case "faculties":
			err = database.DB.Raw(`
				SELECT d.institution_id
				FROM faculties f
				JOIN departments d ON f.department_id = d.id
				WHERE f.id = ?
			`, id).Scan(&institutionID).Error

		case "departments":
			err = database.DB.Raw(`
				SELECT institution_id
				FROM departments
				WHERE id = ?
			`, id).Scan(&institutionID).Error

		case "institutes":
			institutionID = uint(id)

		default:
			return helper.Error(c, 400, "Unsupported resource")
		}

		if err != nil {
			return helper.Error(c, 500, err.Error())
		}

		if !helper.HasInstitutionAccess(userID, institutionID) {
			return helper.Error(c, 403, "You don't have access to this institution")
		}

		return c.Next()
		
	}
	
}
