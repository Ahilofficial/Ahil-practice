package controller

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/helper"
	"backend_institutions/internal/services"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type FeesController struct {
	feesService *services.FeesService
}

func NewFeesController(feesService *services.FeesService) *FeesController {
	return &FeesController{feesService: feesService}
}

func (cl *FeesController) GetInactiveFeesController(c fiber.Ctx) error {
	fees, err := cl.feesService.GetInactiveFeesService()
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Inactive fees fetched successfully",
		dto.ToFeesResponseListDTO(fees),
	)
}

func (cl *FeesController) CreateFeesController(c fiber.Ctx) error {
	var body dto.CreateFeesDTO
	body.Sanitize()

	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	if err := body.Validate(); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	fees, err := cl.feesService.CreateFeesService(&body)
	if err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(
		c,
		"Fees created successfully",
		dto.ToFeesResponseDTO(&fees),
	)
}

func (cl *FeesController) GetAllFeesController(c fiber.Ctx) error {
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

	fees, total, err := cl.feesService.GetFeesServicePaginated(page, limit)
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return helper.Success(
		c,
		"Fees fetched successfully",
		fiber.Map{
			"items":       dto.ToFeesResponseListDTO(fees),
			"total_count": total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	)
}

func (cl *FeesController) GetFeesByIDController(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid fees id")
	}

	fees, err := cl.feesService.GetFeesServiceById(uint(id))
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Fees fetched successfully",
		dto.ToFeesResponseDTO(&fees),
	)
}

func (cl *FeesController) UpdateFeesController(c fiber.Ctx) error {
	var body dto.UpdateFeesDTO
	body.Sanitize()

	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid fees id")
	}

	
	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	if err := body.Validate(); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	if err := cl.feesService.UpdateFeesService(
		uint(id),
		&body,
	); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	updated, err := cl.feesService.GetFeesServiceById(uint(id))
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"Fees updated successfully",
		dto.ToFeesResponseDTO(&updated),
	)
}

func (cl *FeesController) DeleteFeesController(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid fees id")
	}

	if err := cl.feesService.DeleteFeesService(uint(id)); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(
		c,
		"Fees deleted successfully",
		nil,
	)
}

func (cl *FeesController) FetchAllFeesController(c fiber.Ctx) error {
	fees, err := cl.feesService.GetFeesService()
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"All fees fetched successfully",
		dto.ToFeesResponseListDTO(fees),
	)
}
