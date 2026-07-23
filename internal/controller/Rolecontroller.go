package controller

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/helper"
	"backend_institutions/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type RoleController struct {
	roleService *services.RoleService
}

func NewRoleController(roleService *services.RoleService) *RoleController {
	return &RoleController{roleService: roleService}
}

func (cl *RoleController) CreateRoleController(c fiber.Ctx) error {
	var body dto.CreateRoleDTO
	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "Invalid request body")
	}

	body.Sanitize()
	if err := body.Validate(); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	role, err := cl.roleService.CreateRole(&body)
	if err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(c, "Role created successfully", role)
}



func (cl *RoleController) GetRoleByIDController(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id == 0 {
		return helper.Error(c, 400, "Invalid role ID")
	}

	role, err := cl.roleService.GetRoleByID(uint(id))
	if err != nil {
		return helper.Error(c, 404, "Role not found")
	}

	return helper.Success(c, "Role retrieved successfully", role)
}


func (cl *RoleController) AssignPermissionsController(c fiber.Ctx) error {
	idParam := c.Params("id")
	roleID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || roleID == 0 {
		return helper.Error(c, 400, "Invalid role ID")
	}

	var body dto.AssignPermissionsDTO
	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "Invalid request body")
	}

	if err := body.Validate(); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	if err := cl.roleService.AssignPermissionsToRole(uint(roleID), &body); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(c, "Permissions assigned to role successfully", nil)
}

func (cl *RoleController) GetRolePermissionsController(c fiber.Ctx) error {
	idParam := c.Params("id")
	roleID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || roleID == 0 {
		return helper.Error(c, 400, "Invalid role ID")
	}

	perms, err := cl.roleService.GetRolePermissions(uint(roleID))
	if err != nil {
		return helper.Error(c, 500, err.Error())
	}

	return helper.Success(c, "Role permissions retrieved successfully", perms)
}

func (cl *RoleController) RemovePermissionController(c fiber.Ctx) error {
	idParam := c.Params("id")
	roleID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || roleID == 0 {
		return helper.Error(c, 400, "Invalid role ID")
	}

	permParam := c.Params("permissionId")
	permID, err := strconv.ParseUint(permParam, 10, 32)
	if err != nil || permID == 0 {
		return helper.Error(c, 400, "Invalid permission ID")
	}

	if err := cl.roleService.RemovePermissionFromRole(uint(roleID), uint(permID)); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(c, "Permission removed from role successfully", nil)
}
