package repository

import (
"backend_institutions/database"
"backend_institutions/model"
)

type FeesRepository struct{}

func NewFeesRepository() *FeesRepository {
return &FeesRepository{}
}

func (r *FeesRepository) CreateFees(fees *model.Fees) error {
return database.DB.Create(fees).Error
}

func (r *FeesRepository) FetchFees() ([]model.Fees, error) {
var fees []model.Fees
err := database.DB.Find(&fees).Error
return fees, err
}

func (r *FeesRepository) FetchFeesById(id uint) (model.Fees, error) {
var fees model.Fees
err := database.DB.First(&fees, id).Error
return fees, err
}

func (r *FeesRepository) DeleteFees(id uint) error {
return database.DB.Delete(&model.Fees{}, id).Error
}

func (r *FeesRepository) UpdateFeesById(fees *model.Fees) error {
return database.DB.Save(fees).Error
}
