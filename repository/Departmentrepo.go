package repository

import (
"backend_institutions/database"
"backend_institutions/model"
)

type DepartmentRepository struct{}

func NewDepartmentRepository() *DepartmentRepository {
return &DepartmentRepository{}
}

func (r *DepartmentRepository) CreateDepartment(department *model.Department) error {
return database.DB.Create(department).Error
}

func (r *DepartmentRepository) FetchDepartment() ([]model.Department, error) {
var department []model.Department
err := database.DB.Find(&department).Error
return department, err
}

func (r *DepartmentRepository) FetchDepartmentById(id uint) (model.Department, error) {
var department model.Department
err := database.DB.First(&department, id).Error
return department, err
}

func (r *DepartmentRepository) DeleteDepartment(id uint) error {
return database.DB.Delete(&model.Department{}, id).Error
}

func (r *DepartmentRepository) UpdateDepartmentById(department *model.Department) error {
return database.DB.Save(department).Error
}
