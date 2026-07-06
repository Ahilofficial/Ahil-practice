package repository

import (
	"backend_institutions/database"
	"backend_institutions/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

func CreateFees(fees *model.Fees) error {
	return database.DB.Create(fees).Error
}

func FetchFees() ([]model.Fees, error) {
	var fees []model.Fees
	err := database.DB.Where("deleted_at IS NULL").Find(&fees).Error
	return fees, err
}

func FetchFeesPaginated(page, limit int) ([]model.Fees, int64, error) {
	var fees []model.Fees
	var total int64

	err := database.DB.Model(&model.Fees{}).Where("deleted_at IS NULL").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = database.DB.Where("deleted_at IS NULL").Offset(offset).Limit(limit).Find(&fees).Error
	return fees, total, err
}

func FetchFeesById(id uint) (model.Fees, error) {
	var fees model.Fees
	err := database.DB.Where("id = ? AND isactive = ?", id, true).First(&fees).Error
	return fees, err
}

// func  FetchActiveFees() ([]model.Fees, error) {
// 	return r.FetchFees()
// }

func FetchInactiveFees() ([]model.Fees, error) {
	var fees []model.Fees

	err := database.DB.
		Where("isactive = ?", false).
		Find(&fees).Error

	if err != nil {
		return nil, err
	}

	return fees, nil
}

// func  FetchFeesDeleted() ([]model.Fees, error) {
// 	return r.FetchInactiveFees()
// }

func DeleteFees(id uint) error {
	var existing model.Fees
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

	return database.DB.Model(&model.Fees{}).
		Where("id = ?", id).
		Update("isactive", false).
		Update("deleted_at", time.Now()).
		Error
}

func UpdateFeesById(fees *model.Fees) error {
	return database.DB.Save(fees).Error
}
