package controller

import (
	"backend_institutions/dto"
	"backend_institutions/helper"
	"backend_institutions/model"
	"backend_institutions/services"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func CreateFacultyController(c fiber.Ctx) error {
	var body dto.CreateFacultyDTO

	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	faculty := model.Faculty{
		Name:         body.Name,
		Gender:       body.Gender,
		JoiningDate:  body.JoiningDate,
		DepartmentID: body.DepartmentID,
	}

	if err := services.CreateFacultyService(&faculty); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(
		c,
		"Faculty created successfully",
		faculty,
	)
}

func GetAllFacultiesController(c fiber.Ctx) error {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	faculties, total, err := services.GetFacultyServicePaginated(page, limit)
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return helper.Success(
		c,
		"Faculties fetched successfully",
		fiber.Map{
			"items":       faculties,
			"total_count": total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	)
}

func GetFacultyByIDController(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid faculty ID")
	}

	faculty, err := services.GetFacultyServiceById(uint(id))
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Faculty fetched successfully",
		faculty,
	)
}

func GetDeletedFacultiesController(c fiber.Ctx) error {
	faculties, err := services.GetFacultyServiceDeleted()
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"Deleted faculties fetched successfully",
		faculties,
	)
}

func GetActiveFacultyController(c fiber.Ctx) error {
	faculty, err := services.GetActiveFacultyService()
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Active faculty fetched successfully",
		faculty,
	)
}

func GetInactiveFacultyController(c fiber.Ctx) error {
	faculty, err := services.GetInactiveFacultyService()
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Inactive faculty fetched successfully",
		faculty,
	)
}

func UpdateFacultyController(c fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid id")
	}

	var faculty model.Faculty

	if err := c.Bind().Body(&faculty); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	if err := services.UpdateFacultyService(
		uint(id),
		faculty.Name,
		faculty.Gender,
	); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	updated, err := services.GetFacultyServiceById(uint(id))
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"Faculty updated successfully",
		updated,
	)
}

func DeleteFacultyController(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid faculty id")
	}

	if err := services.DeleteFacultyService(uint(id)); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(
		c,
		"Faculty deleted successfully",
		nil,
	)
}

func FetchAllFacultiesController(c fiber.Ctx) error {
	faculties, err := services.GetFacultyService()
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"All faculties fetched successfully",
		faculties,
	)
}
