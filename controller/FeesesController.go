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



func GetInactiveFeesController(c fiber.Ctx) error {
	fees, err := services.GetInactiveFeesService()
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Inactive fees fetched successfully",
		fees,
	)
}

func CreateFeesController(c fiber.Ctx) error {
	var body dto.CreateFeesDTO

	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	fees := model.Fees{
		PaymentMode: body.PaymentMode,
		Amount:      body.Amount,
		StudentID:   body.StudentID,
	}

	if err := services.CreateFeesService(&fees); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(
		c,
		"Fees created successfully",
		fees,
	)
}

func GetAllFeesController(c fiber.Ctx) error {
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

	fees, total, err := services.GetFeesServicePaginated(page, limit)
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return helper.Success(
		c,
		"Fees fetched successfully",
		fiber.Map{
			"items":       fees,
			"total_count": total,
			"page":        page,
			"limit":       limit,
			"total_pages": totalPages,
		},
	)
}

func GetFeesByIDController(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid fees id")
	}

	fees, err := services.GetFeesServiceById(uint(id))
	if err != nil {
		return helper.Error(c, 404, err.Error())
	}

	return helper.Success(
		c,
		"Fees fetched successfully",
		fees,
	)
}

// func GetActiveFeesController(c fiber.Ctx) error {
// 	fees, err := services.GetActiveFeesService()
// 	if err != nil {
// 		return helper.Error(c, 404, err.Error())
// 	}

// 	return helper.Success(
// 		c,
// 		"Active fees fetched successfully",
// 		fees,
// 	)
// }



// func GetDeletedFeesController(c fiber.Ctx) error {
// 	fees, err := services.GetFeesServiceDeleted()
// 	if err != nil {
// 		return helper.Error(c, 500, err.Error())
// 	}

// 	return helper.Success(
// 		c,
// 		"Deleted fees fetched successfully",
// 		fees,
// 	)
// }

func UpdateFeesController(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid fees id")
	}

	var fees model.Fees

	if err := c.Bind().Body(&fees); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	if err := services.UpdateFeesService(
		uint(id),
		fees.PaymentMode,
		fees.Amount,
	); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	updated, err := services.GetFeesServiceById(uint(id))
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"Fees updated successfully",
		updated,
	)
}

func DeleteFeesController(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid fees id")
	}

	if err := services.DeleteFeesService(uint(id)); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(
		c,
		"Fees deleted successfully",
		nil,
	)
}

func FetchAllFeesController(c fiber.Ctx) error {
	fees, err := services.GetFeesService()
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(
		c,
		"All fees fetched successfully",
		fees,
	)
}
