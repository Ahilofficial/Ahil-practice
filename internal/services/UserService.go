package services

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/grpc"
	"backend_institutions/internal/model"
	"backend_institutions/internal/repository"
	"backend_institutions/internal/utils"
	"errors"
	"fmt"
	"log"

	// "os"

	// "os/user"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v3"
	// "github.com/gofiber/fiber/v3"
	// "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserService struct {
	userrepo *repository.UserRepository
	sessionService *SessionService
}

func NewUserService(userrepo *repository.UserRepository, sessionService *SessionService,) *UserService {
	return &UserService{
		userrepo: userrepo,
		sessionService: sessionService,
		
	}
}

func (s *UserService) SignUp(dto *dto.SignUpDTO) (model.User, error) {
	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return model.User{}, err
	}

	token := utils.SignUpToken()
	if token == "" {
		return model.User{}, errors.New("failed to generate verification token")
	}
	var expiry = time.Now().Add(24 * time.Hour)

	user := model.User{
		Name:              dto.Name,
		Email:             dto.Email,
		Phone:             dto.Phone,
		Password:          hashedPassword,
		IsActive:          false,
		IsVerified:        false,
		VerificationToken: token,
		TokenExpiresAt:    expiry,
	}

	err = s.userrepo.CreateUser(&user)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			if strings.Contains(mysqlErr.Message, "email") {
				return model.User{}, errors.New("email already exists")
			}
			if strings.Contains(mysqlErr.Message, "phone") {
				return model.User{}, errors.New("phone number already exists")
			}
		}
		return model.User{}, err
	}

	verifyURL := fmt.Sprintf("http://localhost:8090/auth/verify?token=%s", token)
	subject := "Verify your email - Backend Institutions"
	body := fmt.Sprintf(`<h1>Hello %s,</h1>
<p>Thank you for signing up. Please verify your email by clicking the link below:</p>
<p><a href="%s" style="background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; display: inline-block; border-radius: 5px;">Verify Email Address</a></p>
<p>Or copy and paste this link in your browser:<br/><a href="%s">%s</a></p>
<p>This link will expire in 24 hours.</p>
<p>If you did not create this account, please ignore this email.</p>`, user.Name, verifyURL, verifyURL, verifyURL)

	defaultRole, roleErr := s.userrepo.FindRoleByName("user")
	if roleErr == nil && defaultRole.ID != 0 {
		_ = s.userrepo.AssignRoleToUser(user.ID, defaultRole.ID)
	}


	


	// Send verification email asynchronously
	go func(email, subject, body string) {
		if sendErr := grpc.SendEmail(email, subject, body, "signup"); sendErr != nil {
			log.Printf("Failed to send verification email via gRPC: %v\n", sendErr)
		}
	}(user.Email, subject, body)

	return user, nil
}

func (s *UserService) VerifyEmail(token string) error {
	user, err := s.userrepo.FindByVerificationToken(token)
	if err != nil {
		return errors.New("invalid verification link")
	}

	if user.IsVerified {
		return errors.New("email already verified")
	}

	if time.Now().After(user.TokenExpiresAt) {
		return errors.New("verification link expired")
	}

	user.IsVerified = true
	user.IsActive = true
	user.VerificationToken = ""
	user.TokenExpiresAt = time.Time{}

	return s.userrepo.UpdateUser(&user)
}


func (s *UserService) SignIn(dto *dto.SignInDTO, c fiber.Ctx) (string, string,uint, string, error) {
	user, err := s.userrepo.FindByEmail(dto.Email)
	

	
	if err != nil {
		return "", "", 0,"", errors.New("invalid email or password")
	}

	if !user.IsActive {
		return "", "", 0,"", errors.New("account is inactive")
	}

	if !user.IsVerified {
		return "", "", 0,"", errors.New("please verify your email before signing in")
	}

	err = utils.ComparePassword(user.Password, dto.Password)
	if err != nil {
		return "", "", 0,"",errors.New("invalid email or password")
	}
	
	go func(email, name string) {
		subject := "New Sign-In"

		body := fmt.Sprintf(
			"Hello %s,\n\nYour account has just been signed in successfully.\n\nIf this wasn't you, please contact support immediately.",
			name,
		)

		if sendErr := grpc.SendEmail(email, subject, body, "signin"); sendErr != nil {
			log.Printf("Failed to send sign-in email via gRPC: %v\n", sendErr)
		}
	}(user.Email, user.Name)
	sessionID := uuid.New().String()
	userAgent := c.Get("User-Agent")

	

	accessToken, err := utils.GenerateAccessToken(user.ID,sessionID )
	if err != nil {
		return "", "",0, "", err
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID, sessionID)
	if err != nil {
		return "", "",0, "", err
	}
	_, err = s.sessionService.CreateSession(user.ID, userAgent,sessionID, accessToken,refreshToken )
	if err != nil {
	return "", "", 0, "", err
}
	
	

	
	return accessToken, refreshToken, user.ID, sessionID, nil
}

func (s *UserService) AssignRole(userID uint, roleID uint) error {
	return s.userrepo.AssignRoleToUser(userID, roleID)
}

func (s *UserService) DeleteUserService(id uint) error {
	return s.userrepo.DeleteUser(id)
}
func (s *UserService) ForgotPasswordService(mail dto.ForgotPasswordDTO) (model.User, error) {
	fetchemail, err := s.userrepo.ForgotPasswordRepo(mail)
	if err != nil {
		return model.User{}, err
	}
	subject := "Forgot Password mail"
	token := utils.ReseTToken()
	resetURL := fmt.Sprintf("http://localhost:8090/auth/reset-password?token=%s", token)
	fetchemail.ResetPasswordToken = token
	fetchemail.ResetTokenExpiresAt = time.Now().Add(15 * time.Minute)
	
	
	body := fmt.Sprintf(`
	<h2>Hello %s,</h2>

	<p>We received a request to reset your password for your <strong>Backend Institutions</strong> account.</p>

	<p>Click the button below to reset your password:</p>

	<p>
		<a href="%s"
		   style="
				background-color:#4CAF50;
				color:white;
				padding:12px 24px;
				text-decoration:none;
				border-radius:5px;
				display:inline-block;">
			Reset Password
		</a>
	</p>

	<p>If the button doesn't work, copy and paste this link into your browser:</p>

	<p>%s</p>

	<p><strong>This link will expire in 15 minutes.</strong></p>

	<p>If you did not request a password reset, you can safely ignore this email. Your password will remain unchanged.</p>

	<br>

	<p>Regards,</p>
	<p><strong>Backend Institutions Team</strong></p>
`, fetchemail.Name, resetURL, resetURL)


	go func(email, subject, body, resetURL string) {
		if sendErr := grpc.SendEmail(email, subject, body, "forgot-password"); sendErr != nil {
			log.Printf("Failed to send verification email via gRPC: %v\n", sendErr)
		}
	}(fetchemail.Email, subject, body, resetURL)

	err = s.userrepo.UpdateResetToken(fetchemail)
	
	

	return fetchemail, err

}

func (s *UserService) ResetPasswordService(token string,reset dto.ResetPassword) error {
	
	
	user, err := s.userrepo.FetchUsertoken(token)
	
	
	
	if err != nil {
		return err
	}
	err = utils.ComparePassword(user.Password, reset.CurrentPassword)
	if err != nil {
    return errors.New("current password is incorrect")
}
	hashedPassword, err := utils.HashPassword(reset.NewPassword)
	if err != nil {
    return err
}
fmt.Println("User ID:", user.ID)
fmt.Println("New Password:", reset.NewPassword)
fmt.Println("Hashed Password:", hashedPassword)
err = s.userrepo.UpdatePassword(user.ID, hashedPassword)

	err = s.userrepo.UpdatePassword(user.ID, hashedPassword)
	if err != nil {
		fmt.Println("Cant able to update the password")
	}
	return err

}


func (s *UserService) Logout(dto *dto.LogoutDTO) error {
	return s.userrepo.Logout(dto)
}

func (s *UserService) CheckEmailPresent(signin dto.SignInDTO, mail dto.ForgotPasswordDTO) error {

    if signin.Email != mail.Email {
        return errors.New("email does not match")
    }

    return nil
}