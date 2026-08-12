package controller

import (
	"backend_institutions/internal/services"

	"github.com/gofiber/fiber/v3"
)

type MenuController struct {
	menuService *services.MenuService
}

func NewMenuController(service *services.MenuService) *MenuController {

	return &MenuController{
		menuService: service,
	}
}

func (c *MenuController) GetMenus(ctx fiber.Ctx) error {

	userID := ctx.Locals("user_id").(uint)

	menus, err := c.menuService.GetMenus(userID)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"menus": menus,
	})
}
