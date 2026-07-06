package services

import (
	"backend_institutions/model"
	"backend_institutions/repository"
	"backend_institutions/utils"
	"errors"
	"regexp"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9._%+-]{0,62}[a-zA-Z0-9])?@(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,63}$`)
	phoneRegex = regexp.MustCompile(`^[0-9]{10}$`)
)

func SignUp(name, email, phone, password, role string) (model.User, error) {
	if name == "" || email == "" || phone == "" || password == "" {
		return model.User{}, errors.New("all fields (name, email, phone, password) are required")
	}

	if !emailRegex.MatchString(email) {
		return model.User{}, errors.New("invalid email format")
	}

	if !phoneRegex.MatchString(phone) {
		return model.User{}, errors.New("phone number must be exactly 10 digits")
	}

	// Check if user already exists by email
	existingUser, err := repository.FindByEmail(email)
	if err == nil && existingUser.ID > 0 {
		return model.User{}, errors.New("email already exists")
	}

	// Check if user already exists by phone
	existingUserByPhone, err := repository.FindByPhone(phone)
	if err == nil && existingUserByPhone.ID > 0 {
		return model.User{}, errors.New("phone number already exists")
	}

	// Hashing password

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		Name:     name,
		Email:    email,
		Phone:    phone,
		Password: hashedPassword,
		IsActive: true,
	}

	if role == "" {
		role = "user"
	} else if role != "user" {
		return model.User{}, errors.New("can't able to assign role other than user during signup")
	}

	err = repository.CreateUser(&user)
	if err != nil {
		return model.User{}, err
	}

	err = repository.AssignRoleToUser(user.ID, role)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func SignIn(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("email and password are required")
	}

	user, err := repository.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !user.IsActive {
		return "", errors.New("account is inactive")
	}

	// Compare password using utils utility
	err = utils.ComparePassword(user.Password, password)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate a JWT Token
	return utils.GenerateToken(user.ID)
}

func AssignRole(userID uint, roleName string) error {
	if roleName != "user" && roleName != "principal" && roleName != "faculty" && roleName != "student" && roleName != "admin" {
		return errors.New("invalid role name")
	}
	return repository.AssignRoleToUser(userID, roleName)
}
