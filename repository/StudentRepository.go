package repository

import (
	"backend_institutions/database"
	"backend_institutions/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

func CreateStudent(student *model.Student) error {
	return database.DB.Create(student).Error
}

func FetchStudent() ([]model.Student, error) {
	var student []model.Student
	err := database.DB.Where("deleted_at IS NULL").
		Preload("Fees", "deleted_at IS NULL").
		Find(&student).Error
	return student, err
}

func GetActiveStudent() (model.Student, error) {
	var student model.Student

	err := database.DB.
		Where("isactive = ?", true).
		First(&student).Error

	if err != nil {
		return model.Student{}, err
	}

	return student, nil
}

func GetInactiveStudent() (model.Student, error) {
	var student model.Student

	err := database.DB.
		Where("isactive = ?", false).
		First(&student).Error

	if err != nil {
		return model.Student{}, err
	}

	return student, nil
}

func FetchStudentPaginated(page, limit int) ([]model.Student, int64, error) {
	var student []model.Student
	var total int64

	err := database.DB.Model(&model.Student{}).Where("deleted_at IS NULL").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = database.DB.Where("deleted_at IS NULL").
		Preload("Fees", "deleted_at IS NULL").
		Offset(offset).Limit(limit).Find(&student).Error
	return student, total, err
}

func FetchStudentById(id uint) (model.Student, error) {
	var student model.Student
	err := database.DB.Where("id = ? AND deleted_at IS NULL", id).
		Preload("Fees", "deleted_at IS NULL").
		First(&student).Error
	return student, err
}

func FetchStudentDeleted() ([]model.Student, error) {
	var student []model.Student
	err := database.DB.Unscoped().Where("deleted_at IS NOT NULL").
		Preload("Fees", "deleted_at IS NULL").
		Find(&student).Error
	return student, err
}

func DeleteStudent(id uint) error {
	var existing model.Student
	err := database.DB.Unscoped().Where("id = ?", id).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("record not found")
		}
		return err
	}

	if !existing.IsActive {
		return errors.New("already deleted")
	}

	return database.DB.Model(&model.Student{}).
		Where("id = ?", id).
		Update("isactive", false).
		Update("deleted_at", time.Now()).
		Error
}

func UpdateStudentById(student *model.Student) error {
	return database.DB.Save(student).Error
}
