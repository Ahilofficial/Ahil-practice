package repository

import (
	"backend_institutions/database"
	"backend_institutions/model"
	"errors"
	"time"

	"gorm.io/gorm"
)



func CreateDepartment(department *model.Department) error {
	return database.DB.Create(department).Error
}

//fetch department

func FetchDepartment() ([]model.Department, error) {
	var department []model.Department
	err := database.DB.Where("deleted_at IS NULL").
		Preload("Faculties", "deleted_at IS NULL").
		Preload("Faculties.Students", "deleted_at IS NULL").
		Preload("Faculties.Students.Fees", "deleted_at IS NULL").
		Find(&department).Error
	return department, err
}

//fetch paginated list

func FetchDepartmentPaginated(page, limit int) ([]model.Department, int64, error) {
	var department []model.Department
	var total int64

	err := database.DB.Model(&model.Department{}).Where("deleted_at IS NULL").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = database.DB.Where("deleted_at IS NULL").
		Preload("Faculties", "deleted_at IS NULL").
		Preload("Faculties.Students", "deleted_at IS NULL").
		Preload("Faculties.Students.Fees", "deleted_at IS NULL").
		Offset(offset).Limit(limit).Find(&department).Error
	return department, total, err
}

func FetchDepartmentById(id uint) (model.Department, error) {
	var department model.Department
	err := database.DB.Where("id = ? AND deleted_at IS NULL", id).
		Preload("Faculties", "deleted_at IS NULL").
		Preload("Faculties.Students", "deleted_at IS NULL").
		Preload("Faculties.Students.Fees", "deleted_at IS NULL").
		First(&department).Error
	return department, err
}

func FetchDepartmentDeleted() ([]model.Department, error) {
	var department []model.Department
	err := database.DB.Unscoped().Where("deleted_at IS NOT NULL").
		Preload("Faculties", "deleted_at IS NULL").
		Preload("Faculties.Students", "deleted_at IS NULL").
		Preload("Faculties.Students.Fees", "deleted_at IS NULL").
		Find(&department).Error
	return department, err
}

func GetActiveDepartment() (model.Department, error) {
	var department model.Department

	err := database.DB.
		Where("isactive = ?", true).
		First(&department).Error

	if err != nil {
		return model.Department{}, err
	}
	return department, nil
}

func GetInactiveDepartment() (model.Department, error) {
	var department model.Department

	err := database.DB.
		Where("isactive = ?", false).
		First(&department).Error

	if err != nil {
		return model.Department{}, err
	}
	return department, nil
}
func DeleteDepartment(id uint) error {
	var existing model.Department
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

	return database.DB.Model(&model.Department{}).
		Where("id = ?", id).
		Update("isactive", false).
		Update("deleted_at", time.Now()).
		Error
}

func UpdateDepartmentById(department *model.Department) error {
	return database.DB.Save(department).Error
}
