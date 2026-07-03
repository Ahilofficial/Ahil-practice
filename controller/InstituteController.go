package controller

import (
	"backend_institutions/dto"
	"backend_institutions/helper"
	"backend_institutions/model"
	"backend_institutions/services"
	"github.com/gofiber/fiber/v3"
	"math"
	"strconv"
)



func  CreateInstituteController(c fiber.Ctx) error {
	var body dto.CreateInstitutionDTO

	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	institute := model.Institutions{
		Name:            body.Name,
		InstitutionCode: body.InstitutionCode,
		State:           body.State,
	}

	if err := services.CreateInsituteService(&institute); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(
		c,
		"Institution created successfully",
		institute,
	)
}

func  GetAllInstitutesController(c fiber.Ctx) error {
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

	institutes, total, err := services.GetInstituteServicePaginated(page, limit)
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return helper.Success(
		c,
		"Institutes fetched successfully",
		fiber.Map{
			"items":       institutes,
			"total_count": total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	)
}

func  GetInstituteByIDController(c fiber.Ctx) error {
	idstr := c.Params("id")
	id, err := strconv.ParseUint(idstr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "Invalid institute ID")
	}

	institute, err := services.GetInstituteServiceById(uint(id))
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Institute fetched successfully",
		institute,
	)
}

func  GetDeletedInstitutesController(c fiber.Ctx) error {
	institutes, err := services.GetInstituteServiceDeleted()
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"Deleted institutes fetched successfully",
		institutes,
	)
}
func  GetActiveInstituteController(c fiber.Ctx) error {
	institute, err := services.GetActiveInstitute()
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Active institution fetched successfully",
		institute,
	)
}

func  GetInactiveInstituteController(c fiber.Ctx) error {
	institute, err := services.GetInactiveInstitute()
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Inactive institution fetched successfully",
		institute,
	)
}

func  UpdateInstituteController(c fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid id")
	}

	var institute model.Institutions

	if err := c.Bind().Body(&institute); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	if err := services.UpdateInstitutionService(
		uint(id),
		institute.Name,
		institute.InstitutionCode,
		institute.State,
	); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	// Fetch updated record to return in response
	updated, err := services.GetInstituteServiceById(uint(id))
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"Institution updated successfully",
		updated,
	)
}

func  DeleteInstituteController(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid institute id")
	}

	if err := services.DeleteInstitutionService(uint(id)); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(
		c,
		"Institution deleted successfully",
		nil,
	)
}

func  FetchAllInstitutesController(c fiber.Ctx) error {
	institutes, err := services.GetInstituteService()
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"All institutes fetched successfully",
		institutes,
	)
}
