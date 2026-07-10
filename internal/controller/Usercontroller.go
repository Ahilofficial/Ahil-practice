package controller

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/helper"
	"backend_institutions/internal/services"

	"github.com/gofiber/fiber/v3"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (cl *UserController) SignUpController(c fiber.Ctx) error {
	var body dto.SignUpDTO
	body.Sanitize()
	
	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "invalid request body: "+err.Error())
	}

	if err := body.Validate(); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	user, err := cl.userService.SignUp(&body)
	if err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(c, "Signed up successfully", dto.ToUserResponseDTO(&user))
}

func (cl *UserController) SignInController(c fiber.Ctx) error {
	var body dto.SignInDTO
	body.Sanitize()
	
	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "invalid request body: "+err.Error())
	}

	if err := body.Validate(); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	token, err := cl.userService.SignIn(&body)
	if err != nil {
		return helper.Error(c, 401, err.Error())
	}

	return helper.Success(c, "Signed in successfully", dto.AuthResponseDTO{Token: token})
}

func (cl *UserController) AssignRoleController(c fiber.Ctx) error {
	var body dto.AssignRoleDTO
	body.Sanitize()
	
	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "invalid request body: "+err.Error())
	}

	if err := body.Validate(); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	if err := cl.userService.AssignRole(body.UserID, body.Role); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(c, "Role assigned successfully", nil)
}
