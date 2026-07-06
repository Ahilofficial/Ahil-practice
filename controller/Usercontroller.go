package controller

import (
	"backend_institutions/dto"
	"backend_institutions/helper"
	"backend_institutions/services"

	"github.com/gofiber/fiber/v3"
)

func SignUpController(c fiber.Ctx) error {
	var body dto.SignUpDTO
	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "invalid request body: "+err.Error())
	}

	user, err := services.SignUp(body.Name, body.Email, body.Phone, body.Password, body.Role)
	if err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(c, "Signed up successfully", user)
}

func SignInController(c fiber.Ctx) error {
	var body dto.SignInDTO
	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "invalid request body: "+err.Error())
	}

	token, err := services.SignIn(body.Email, body.Password)
	if err != nil {
		return helper.Error(c, 401, err.Error())
	}

	return helper.Success(c, "Signed in successfully", dto.AuthResponseDTO{Token: token})
}

func AssignRoleController(c fiber.Ctx) error {
	var body dto.AssignRoleDTO
	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, "invalid request body: "+err.Error())
	}

	if err := services.AssignRole(body.UserID, body.Role); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(c, "Role assigned successfully", nil)
}
