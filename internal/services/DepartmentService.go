package services

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/model"
	"backend_institutions/internal/repository"
)

type DepartmentService struct {
	departmentrepo *repository.DepartmentRepository
}

func NewDepartmentService(departmentrepo *repository.DepartmentRepository) *DepartmentService {
	return &DepartmentService{departmentrepo: departmentrepo}
}

func (s *DepartmentService) AddDepartmentService(dto *dto.CreateDepartmentDTO) (model.Department, error) {
	department := model.Department{
		DepartmentName: dto.DepartmentName,
		InstitutionID:  dto.InstitutionID,
	}

	err := s.departmentrepo.CreateDepartment(&department)
	return department, err
}

func (s *DepartmentService) GetDepartmentService() ([]model.Department, error) {
	return s.departmentrepo.FetchDepartment()
}

func (s *DepartmentService) GetDepartmentServicePaginated(page, limit int) ([]model.Department, int64, error) {
	return s.departmentrepo.FetchDepartmentPaginated(page, limit)
}

func (s *DepartmentService) GetDepartmentByIDService(id uint) (model.Department, error) {
	return s.departmentrepo.FetchDepartmentById(id)
}

func (s *DepartmentService) GetDepartmentServiceDeleted() ([]model.Department, error) {
	return s.departmentrepo.FetchDepartmentDeleted()
}

func (s *DepartmentService) DeleteDepartment(id uint) error {
	return s.departmentrepo.DeleteDepartment(id)
}

func (s *DepartmentService) GetActiveDepartmentService() (model.Department, error) {
	return s.departmentrepo.GetActiveDepartment()
}

func (s *DepartmentService) GetInactiveDepartmentService() (model.Department, error) {
	return s.departmentrepo.GetInactiveDepartment()
}

func (s *DepartmentService) UpdateDepartmentService(id uint, dto *dto.UpdateDepartmentDTO) error {
	department, err := s.departmentrepo.FetchDepartmentById(id)
	if err != nil {
		return err
	}
	department.DepartmentName = dto.DepartmentName
	return s.departmentrepo.UpdateDepartmentById(&department)
}
