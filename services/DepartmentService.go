package services

import (
	"backend_institutions/model"
	"backend_institutions/repository"
)

func AddDepartmentService(department *model.Department) error {
	return repository.CreateDepartment(department)
}

func GetDepartmentService() ([]model.Department, error) {
	return repository.FetchDepartment()
}

func GetDepartmentServicePaginated(page, limit int) ([]model.Department, int64, error) {
	return repository.FetchDepartmentPaginated(page, limit)
}

func GetDepartmentByIDService(id uint) (model.Department, error) {
	return repository.FetchDepartmentById(id)
}

func GetDepartmentServiceDeleted() ([]model.Department, error) {
	return repository.FetchDepartmentDeleted()
}

func DeleteDepartment(id uint) error {
	return repository.DeleteDepartment(id)
}
func GetActiveDepartmentService() (model.Department, error) {
	return repository.GetActiveDepartment()
}

func GetInactiveDepartmentService() (model.Department, error) {
	return repository.GetInactiveDepartment()
}
func UpdateDepartmentService(id uint, department_name string) error {

	department, err := repository.FetchDepartmentById(id)
	if err != nil {
		return err
	}
	department.DepartmentName = department_name
	return repository.UpdateDepartmentById(&department)
}
