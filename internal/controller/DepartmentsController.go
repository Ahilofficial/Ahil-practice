package controller

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/helper"
	"backend_institutions/internal/services"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type DepartmentController struct {
	departmentService *services.DepartmentService
}

func NewDepartmentController(departmentService *services.DepartmentService) *DepartmentController {
	return &DepartmentController{departmentService: departmentService}
}

func (cl *DepartmentController) GetActiveDepartmentController(c fiber.Ctx) error {
	
	department, err := cl.departmentService.GetActiveDepartmentService()
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Active department fetched successfully",
		dto.ToDepartmentResponseDTO(&department),
	)
}

func (cl *DepartmentController) GetInactiveDepartmentController(c fiber.Ctx) error {
	department, err := cl.departmentService.GetInactiveDepartmentService()
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Inactive department fetched successfully",
		dto.ToDepartmentResponseDTO(&department),
	)
}

func (cl *DepartmentController) CreateDepartmentController(c fiber.Ctx) error {
	var body dto.CreateDepartmentDTO
	body.Sanitize()

	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	if err := body.Validate(); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	department, err := cl.departmentService.AddDepartmentService(&body)
	if err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(
		c,
		"Department created successfully",
		dto.ToDepartmentResponseDTO(&department),
	)
}

func (cl *DepartmentController) GetAllDepartmentsController(c fiber.Ctx) error {
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

	departments, total, err := cl.departmentService.GetDepartmentServicePaginated(page, limit)
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return helper.Success(
		c,
		"Departments fetched successfully",
		fiber.Map{
			"items":       dto.ToDepartmentResponseListDTO(departments),
			"total_count": total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	)
}

func (cl *DepartmentController) GetDepartmentByIDController(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid department id")
	}

	department, err := cl.departmentService.GetDepartmentByIDService(uint(id))
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Department fetched successfully",
		dto.ToDepartmentResponseDTO(&department),
	)
}

func (cl *DepartmentController) GetDeletedDepartmentsController(c fiber.Ctx) error {
	departments, err := cl.departmentService.GetDepartmentServiceDeleted()
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"Deleted departments fetched successfully",
		dto.ToDepartmentResponseListDTO(departments),
	)
}

func (cl *DepartmentController) UpdateDepartmentController(c fiber.Ctx) error {
	var body dto.UpdateDepartmentDTO
	body.Sanitize()
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid department id")
	}

	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	if err := body.Validate(); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	if err := cl.departmentService.UpdateDepartmentService(
		uint(id),
		&body,
	); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	updated, err := cl.departmentService.GetDepartmentByIDService(uint(id))
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"Department updated successfully",
		dto.ToDepartmentResponseDTO(&updated),
	)
}

func (cl *DepartmentController) DeleteDepartmentController(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid department id")
	}

	if err := cl.departmentService.DeleteDepartment(uint(id)); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(
		c,
		"Department deleted successfully",
		nil,
	)
}

func (cl *DepartmentController) FetchAllDepartmentsController(c fiber.Ctx) error {
	departments, err := cl.departmentService.GetDepartmentService()
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"All departments fetched successfully",
		dto.ToDepartmentResponseListDTO(departments),
	)
}
