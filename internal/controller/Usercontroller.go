package controller

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/helper"
	"backend_institutions/internal/services"
	"fmt"


	// "go/token"
	"strconv"

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


	accessToken, refreshToken,user_id, session_id, err := cl.userService.SignIn(&body, c)
	if err != nil {
		return helper.Error(c, 401, err.Error())
	}

	return helper.Success(c, "Signed in successfully", dto.AuthResponseDTO{
		UserID : user_id,
		SessionID: session_id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		
	})
}

func (c *UserController) VerifyEmail(ctx fiber.Ctx) error {
	token := ctx.Query("token")

	if token == "" {
		return helper.Error(ctx, fiber.StatusBadRequest, "Verification token is required")
	}

	err := c.userService.VerifyEmail(token)
	if err != nil {
		return helper.Error(ctx, fiber.StatusBadRequest, err.Error())
	}

	return helper.Success(ctx, "Email verified successfully", nil)
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

	if err := cl.userService.AssignRole(body.UserID, body.RoleID); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(c, "Role assigned successfully", nil)
}

func (cl *UserController) DeleteUserController(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return helper.Error(c, 400, "invalid user id")
	}

	if err := cl.userService.DeleteUserService(uint(id)); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(c, "User deleted successfully", nil)
}

func( cl *UserController)ForgotPassword(c fiber.Ctx)error{
	var forgotpassword dto.ForgotPasswordDTO
	err:=c.Bind().Body(&forgotpassword)
	if err!=nil{
		fmt.Println("Error:", err)

		return helper.Error(c, 400, "invalid request body: "+err.Error())
	}
	_, err = cl.userService.ForgotPasswordService(forgotpassword)
	if err != nil {
		return helper.Error(c, 400, err.Error())
	}
	// return helper.Success(c, "success", user)
	return c.SendString("Sended the forgot password link just check it")

}
func (cl *UserController) ResetPassword(c fiber.Ctx) error {

	token := c.Query("token")
	if token == "" {
		return helper.Error(c, 400, "reset token is required")
	}

	var body dto.ResetPassword

	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	err := cl.userService.ResetPasswordService(token, body)
	if err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(c, "Password reset successfully", nil)
}

func (cl *UserController) Logout(c fiber.Ctx) error {
	var body dto.LogoutDTO

	if err := c.Bind().Body(&body); err != nil {
		return helper.Error(c, 400, err.Error())
	}

	err := cl.userService.Logout(&body)
	if err != nil {
		return helper.Error(c, 400, err.Error())
	}

	return helper.Success(c,  "Logout successful", nil)
}