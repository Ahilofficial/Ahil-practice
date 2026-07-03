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



func  GetActiveDepartmentController(c fiber.Ctx) error {
	department, err := services.GetActiveDepartmentService()
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Active department fetched successfully",
		department,
	)
}

func  GetInactiveDepartmentController(c fiber.Ctx) error {
	department, err := services.GetInactiveDepartmentService()
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Inactive department fetched successfully",
		department,
	)
}

func  CreateDepartmentController(c fiber.Ctx) error {
	var body dto.CreateDepartmentDTO

	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	department := model.Department{
		DepartmentName: body.DepartmentName,
		InstitutionID:  body.InstitutionID,
	}

	if err := services.AddDepartmentService(&department); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(
		c,
		"Department created successfully",
		department,
	)
}

func  GetAllDepartmentsController(c fiber.Ctx) error {
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

	departments, total, err := services.GetDepartmentServicePaginated(page, limit)
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return helper.Success(
		c,
		"Departments fetched successfully",
		fiber.Map{
			"items":       departments,
			"total_count": total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	)
}

func  GetDepartmentByIDController(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid department id")
	}

	department, err := services.GetDepartmentByIDService(uint(id))
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Department fetched successfully",
		department,
	)
}

func  GetDeletedDepartmentsController(c fiber.Ctx) error {
	departments, err := services.GetDepartmentServiceDeleted()
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"Deleted departments fetched successfully",
		departments,
	)
}

func  UpdateDepartmentController(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid department id")
	}

	var department model.Department

	if err := c.Bind().Body(&department); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	if err := services.UpdateDepartmentService(
		uint(id),
		department.DepartmentName,
	); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	updated, err := services.GetDepartmentByIDService(uint(id))
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"Department updated successfully",
		updated,
	)
}

func  DeleteDepartmentController(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid department id")
	}

	if err := services.DeleteDepartment(uint(id)); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(
		c,
		"Department deleted successfully",
		nil,
	)
}

func  FetchAllDepartmentsController(c fiber.Ctx) error {
	departments, err := services.GetDepartmentService()
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"All departments fetched successfully",
		departments,
	)
}
