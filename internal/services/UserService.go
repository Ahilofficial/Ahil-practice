package services

import (
	"backend_institutions/internal/constants"
	"backend_institutions/internal/dto"
	"backend_institutions/internal/grpc"
	"backend_institutions/internal/model"
	"backend_institutions/internal/repository"
	"backend_institutions/internal/utils"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type UserService struct {
	userrepo *repository.UserRepository
}

func NewUserService(userrepo *repository.UserRepository) *UserService {
	return &UserService{
		userrepo: userrepo,
	}
}

func (s *UserService) SignUp(dto *dto.SignUpDTO) (model.User, error) {
	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Phone:    dto.Phone,
		Password: hashedPassword,
		IsActive: true,
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

	// Assign default "user" role to newly signed up user
	err = s.userrepo.AssignRoleToUser(user.ID, constants.UserRole)
	if err != nil {
		return model.User{}, err
	}

	// Send welcome email asynchronously to not block user registration response
	go func(email, name string) {
		subject := "Welcome to Backend Institutions!"
		body := fmt.Sprintf("<h1>Hello %s,</h1><p>Thank you for registering on our platform. Your account is now active!</p>", name)
		if sendErr := grpc.SendEmail(email, subject, body); sendErr != nil {
			log.Printf("Failed to send welcome email via gRPC: %v\n", sendErr)
		}
	}(user.Email, user.Name)

	return user, nil
}

func (s *UserService) SignIn(dto *dto.SignInDTO) (string, error) {
	user, err := s.userrepo.FindByEmail(dto.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !user.IsActive {
		return "", errors.New("account is inactive")
	}

	err = utils.ComparePassword(user.Password, dto.Password)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	return utils.GenerateToken(user.ID)
}

func (s *UserService) AssignRole(userID uint, roleName string) error {
	return s.userrepo.AssignRoleToUser(userID, roleName)
}

func (s *UserService) DeleteUserService(id uint) error {
	return s.userrepo.DeleteUser(id)
}
