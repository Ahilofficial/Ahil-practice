package repository

import (
	"backend_institutions/database"
	"backend_institutions/model"
	"errors"
	"time"

	"gorm.io/gorm"
)


func  CreateFaculty(faculty *model.Faculty) error {
	return database.DB.Create(faculty).Error
}

func  FetchFaculty() ([]model.Faculty, error) {
	var faculty []model.Faculty
	err := database.DB.Where("deleted_at IS NULL").
		Preload("Students", "deleted_at IS NULL").
		Preload("Students.Fees", "deleted_at IS NULL").
		Find(&faculty).Error
	return faculty, err
}

func  FetchFacultyPaginated(page, limit int) ([]model.Faculty, int64, error) {
	var faculty []model.Faculty
	var total int64

	err := database.DB.Model(&model.Faculty{}).Where("deleted_at IS NULL").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = database.DB.Where("deleted_at IS NULL").
		Preload("Students", "deleted_at IS NULL").
		Preload("Students.Fees", "deleted_at IS NULL").
		Offset(offset).Limit(limit).Find(&faculty).Error
	return faculty, total, err
}

func  FetchFacultyById(id uint) (model.Faculty, error) {
	var faculty model.Faculty
	err := database.DB.Where("id = ? AND isactive = ?", id, true).
		Preload("Students", "deleted_at IS NULL").
		Preload("Students.Fees", "deleted_at IS NULL").
		First(&faculty).Error
	return faculty, err
}

func  GetActiveFacutly() (model.Faculty, error) {
	var faculty model.Faculty

	err := database.DB.
		Where("isactive = ?", true).
		First(&faculty).Error

	if err != nil {
		return model.Faculty{}, err
	}
	return faculty, nil
}

func  GetInactiveFaculty() (model.Faculty, error) {
	var faculty model.Faculty

	err := database.DB.
		Where("isactive = ?", false).
		First(&faculty).Error

	if err != nil {
		return model.Faculty{}, err
	}
	return faculty, nil
}

func  FetchFacultyDeleted() ([]model.Faculty, error) {
	var faculty []model.Faculty
	err := database.DB.Unscoped().Where("deleted_at IS NOT NULL").
		Preload("Students", "deleted_at IS NULL").
		Preload("Students.Fees", "deleted_at IS NULL").
		Find(&faculty).Error
	return faculty, err
}

func  DeleteFaculty(id uint) error {
	var existing model.Faculty
	err := database.DB.Where("id = ?", id).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("record not found")
		}
		return err
	}

	if !existing.IsActive {
		return errors.New("already deleted")
	}

	return database.DB.Model(&model.Faculty{}).
		Where("id = ?", id).
		Update("isactive", false).
		Update("deleted_at", time.Now()).
		Error
}

func  UpdateFacultyById(faculty *model.Faculty) error {
	return database.DB.Save(faculty).Error
}
