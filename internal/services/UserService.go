package services

import (
	"backend_institutions/internal/constants"
	"backend_institutions/internal/dto"
	"backend_institutions/internal/model"
	"backend_institutions/internal/repository"
	"backend_institutions/internal/utils"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
)

type UserService struct {
	userrepo *repository.UserRepository
}

func NewUserService(userrepo *repository.UserRepository) *UserService {
	return &UserService{userrepo: userrepo}
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

	role := dto.Role
	if role == "" {
		role = constants.UserRole
	} else if role != constants.UserRole {
		return model.User{}, errors.New("can't able to assign role other than user during signup")
	}

	err = s.userrepo.CreateUser(&user)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			msg := strings.ToLower(mysqlErr.Message)
			if strings.Contains(msg, "email") {
				return model.User{}, errors.New("email already exists")
			}
			if strings.Contains(msg, "phone") {
				return model.User{}, errors.New("phone number already exists")
			}
			return model.User{}, errors.New("email or phone number already exists")
		}
		return model.User{}, err
	}

	err = s.userrepo.AssignRoleToUser(user.ID, role)
	if err != nil {
		return model.User{}, err
	}

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
