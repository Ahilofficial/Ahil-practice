package controller

import (
	"backend_institutions/internal/helper"
	"backend_institutions/internal/services"

	"github.com/gofiber/fiber/v3"
)

type PermissionController struct {
	permService *services.PermissionService
}

func NewPermissionController(permService *services.PermissionService) *PermissionController {
	return &PermissionController{permService: permService}
}

func (cl *PermissionController) GetAllPermissionsController(c fiber.Ctx) error {
	perms, err := cl.permService.GetAllPermissions()
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}
	return helper.Success(c, "Permissions retrieved successfully", perms)
}
