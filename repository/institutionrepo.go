package repository

import (
"backend_institutions/database"
"backend_institutions/model"
)

type InstitutionRepository struct{}

func NewInstituteRepository() *InstitutionRepository {
return &InstitutionRepository{}
}

func (r *InstitutionRepository) CreateInstitution(institute *model.Institutions) error {
return database.DB.Create(institute).Error
}

func (r *InstitutionRepository) FetchInstitution() ([]model.Institutions, error) {
var institution []model.Institutions
err := database.DB.Find(&institution).Error
return institution, err
}

func (r *InstitutionRepository) FetchInstitutionById(id uint) (model.Institutions, error) {
var institution model.Institutions
err := database.DB.First(&institution, id).Error
return institution, err
}

func (r *InstitutionRepository) DeleteInstitution(id uint) error {
return database.DB.Delete(&model.Institutions{}, id).Error
}

func (r *InstitutionRepository) UpdateInstitution(institute *model.Institutions) error {
return database.DB.Save(institute).Error
}
