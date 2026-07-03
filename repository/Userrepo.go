package repository

import (
	"backend_institutions/database"
	"backend_institutions/model"
)

func  CreateUser(user *model.User) error {
	return database.DB.Create(user).Error
}

func  FindByEmail(email string) (model.User, error) {
	var user model.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return user, err
}

func  FindByPhone(phone string) (model.User, error) {
	var user model.User
	err := database.DB.Where("phone = ?", phone).First(&user).Error
	return user, err
}

func  AssignRoleToUser(userID uint, roleName string) error {
	var user model.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return err
	}

	var role model.Role
	if err := database.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}

	return database.DB.Model(&user).Association("Roles").Replace(&role)
}