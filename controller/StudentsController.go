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

func  CreateStudentControllers(c fiber.Ctx) error {
	var body dto.CreateStudentDTO

	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	student := model.Student{
		Name:      body.Name,
		Email:     body.Email,
		Gender:    body.Gender,
		FacultyID: body.FacultyID,
	}

	if err := services.CreateStudentService(&student); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(
		c,
		"Student created successfully",
		student,
	)
}
func  GetActiveStudentController(c fiber.Ctx) error {
	student, err := services.GetActiveStudentService()
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Active student fetched successfully",
		student,
	)
}

func  GetInactiveStudentController(c fiber.Ctx) error {
	student, err := services.GetInactiveStudentService()
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Inactive student fetched successfully",
		student,
	)
}

func  GetAllStudentsControllers(c fiber.Ctx) error {
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

	students, total, err := services.GetStudentServicePaginated(page, limit)
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return helper.Success(
		c,
		"Students fetched successfully",
		fiber.Map{
			"items":       students,
			"total_count": total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	)
}

func  GetStudentByIDControllers(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid student id")
	}

	student, err := services.GetStudentServiceById(uint(id))
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Student fetched successfully",
		student,
	)
}

func  GetDeletedStudentsController(c fiber.Ctx) error {
	students, err := services.GetStudentServiceDeleted()
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"Deleted students fetched successfully",
		students,
	)
}

func  UpdateStudentControllers(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid student id")
	}

	var student model.Student

	if err := c.Bind().Body(&student); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	if err := services.UpdateStudentService(
		uint(id),
		student.Name,
		student.Email,
		student.Gender,
	); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	updated, err := services.GetStudentServiceById(uint(id))
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"Student updated successfully",
		updated,
	)
}

func  DeleteStudentControllers(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid student id")
	}

	if err := services.DeleteStudentService(uint(id)); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(
		c,
		"Student deleted successfully",
		nil,
	)
}

func  FetchAllStudentsControllers(c fiber.Ctx) error {
	students, err := services.GetStudentService()
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"All students fetched successfully",
		students,
	)
}
