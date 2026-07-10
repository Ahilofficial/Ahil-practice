package services

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/model"
	"backend_institutions/internal/repository"
)

type FacultyService struct {
	facultyrepo *repository.FacultyRepository
}

func NewFacultyService(facultyrepo *repository.FacultyRepository) *FacultyService {
	return &FacultyService{facultyrepo: facultyrepo}
}

func (s *FacultyService) CreateFacultyService(dto *dto.CreateFacultyDTO) (model.Faculty, error) {
	faculty := model.Faculty{
		Name:         dto.Name,
		Gender:       dto.Gender,
		JoiningDate:  dto.JoiningDate,
		DepartmentID: dto.DepartmentID,
	}

	err := s.facultyrepo.CreateFaculty(&faculty)
	return faculty, err
}

func (s *FacultyService) GetFacultyService() ([]model.Faculty, error) {
	return s.facultyrepo.FetchFaculty()
}

func (s *FacultyService) GetFacultyServicePaginated(page, limit int) ([]model.Faculty, int64, error) {
	return s.facultyrepo.FetchFacultyPaginated(page, limit)
}

func (s *FacultyService) GetFacultyServiceById(id uint) (model.Faculty, error) {
	return s.facultyrepo.FetchFacultyById(id)
}

func (s *FacultyService) GetFacultyServiceDeleted() ([]model.Faculty, error) {
	return s.facultyrepo.FetchFacultyDeleted()
}

func (s *FacultyService) GetActiveFacultyService() (model.Faculty, error) {
	return s.facultyrepo.GetActiveFacutly()
}

func (s *FacultyService) GetInactiveFacultyService() (model.Faculty, error) {
	return s.facultyrepo.GetInactiveFaculty()
}

func (s *FacultyService) DeleteFacultyService(id uint) error {
	return s.facultyrepo.DeleteFaculty(id)
}

func (s *FacultyService) UpdateFacultyService(id uint, dto *dto.UpdateFacultyDTO) error {
	faculty, err := s.facultyrepo.FetchFacultyById(id)
	if err != nil {
		return err
	}
	faculty.Name = dto.Name
	faculty.Gender = dto.Gender
	return s.facultyrepo.UpdateFacultyById(&faculty)
}
