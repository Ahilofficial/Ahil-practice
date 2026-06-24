package services

import (
	"backend_institutions/model"
	"backend_institutions/repository"
)

type DepartmentService struct {
	departmentrepo *repository.DepartmentRepository
}

func NewDepartmentService() *DepartmentService {
	return &DepartmentService{
		departmentrepo: repository.NewDepartmentRepository(),
	}
}

func (s DepartmentService) AddDepartmentService(department *model.Department) error {
	return s.departmentrepo.CreateDepartment(department)
}

func (s DepartmentService) GetDepartmentService() ([]model.Department, error) {
	return s.departmentrepo.FetchDepartment()
}

func (s DepartmentService) GetDepartmentByIDService(id uint) (model.Department, error) {
	return s.departmentrepo.FetchDepartmentById(id)
}

func (s DepartmentService) DeleteDepartment(id uint) error {
	return s.departmentrepo.DeleteDepartment(id)
}

// var department model.Department
func (s DepartmentService) UpdateDepartmentService(id uint, department_name string) error {
	department, err := s.departmentrepo.FetchDepartmentById(id)
	if err != nil {
		return err
	}
	department.Department_Name = department_name
	return s.departmentrepo.UpdateDepartmentById(&department)
}
