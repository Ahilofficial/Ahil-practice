package repository

import (
	"backend_institutions/database"
	"backend_institutions/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

func CreateInstitution(institute *model.Institutions) error {
	var existing model.Institutions

	err := database.DB.
		Where("name = ? OR institution_code = ?", institute.Name, institute.InstitutionCode).
		First(&existing).Error

	if err == nil {
		return errors.New("institution already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return database.DB.Create(institute).Error
}

func FetchInstitution() ([]model.Institutions, error) {
	var institution []model.Institutions
	err := database.DB.Where("deleted_at IS NULL").
		Preload("Departments", "deleted_at IS NULL").
		Preload("Departments.Faculties", "deleted_at IS NULL").
		Preload("Departments.Faculties.Students", "deleted_at IS NULL").
		Preload("Departments.Faculties.Students.Fees", "deleted_at IS NULL").
		Find(&institution).Error
	return institution, err
}

func FetchInstitutionPaginated(page, limit int) ([]model.Institutions, int64, error) {
	var institution []model.Institutions
	var total int64

	err := database.DB.Model(&model.Institutions{}).Where("deleted_at IS NULL").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = database.DB.Where("deleted_at IS NULL").
		Preload("Departments", "deleted_at IS NULL").
		Preload("Departments.Faculties", "deleted_at IS NULL").
		Preload("Departments.Faculties.Students", "deleted_at IS NULL").
		Preload("Departments.Faculties.Students.Fees", "deleted_at IS NULL").
		Offset(offset).Limit(limit).Find(&institution).Error
	return institution, total, err
}

func FetchInstitutionById(id uint) (model.Institutions, error) {
	var institution model.Institutions
	err := database.DB.Where("id = ? AND deleted_at IS NULL", id).
		Preload("Departments", "deleted_at IS NULL").
		Preload("Departments.Faculties", "deleted_at IS NULL").
		Preload("Departments.Faculties.Students", "deleted_at IS NULL").
		Preload("Departments.Faculties.Students.Fees", "deleted_at IS NULL").
		First(&institution).Error
	return institution, err
}

func GetActiveInstitute() (model.Institutions, error) {
	var institution model.Institutions

	err := database.DB.
		Where("isactive = ?", true).
		First(&institution).Error

	if err != nil {
		return model.Institutions{}, err
	}
	return institution, nil
}

func GetInactiveInstitute() (model.Institutions, error) {
	var institute model.Institutions

	err := database.DB.
		Where("isactive = ?", false).
		First(&institute).Error

	if err != nil {
		return model.Institutions{}, err
	}
	return institute, nil
}

func FetchInstitutionDeleted() ([]model.Institutions, error) {
	var institution []model.Institutions
	err := database.DB.Unscoped().Where("deleted_at IS NOT NULL").
		Preload("Departments", "deleted_at IS NULL").
		Preload("Departments.Faculties", "deleted_at IS NULL").
		Preload("Departments.Faculties.Students", "deleted_at IS NULL").
		Preload("Departments.Faculties.Students.Fees", "deleted_at IS NULL").
		Find(&institution).Error
	return institution, err
}

func DeleteInstitution(id uint) error {
	var existing model.Institutions
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

	return database.DB.Model(&model.Institutions{}).
		Where("id = ?", id).
		Update("isactive", false).
		Update("deleted_at", time.Now()).
		Error
}

func UpdateInstitution(institute *model.Institutions) error {
	return database.DB.Save(institute).Error
}
