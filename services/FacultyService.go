package services

import (
"backend_institutions/model"
"backend_institutions/repository"
)

type FacultyService struct {
facultyrepo *repository.FacultyRepository
}

func NewFacultyService() *FacultyService {
return &FacultyService{
facultyrepo: repository.NewFacultyRepository(),
}
}

func (s *FacultyService) CreateFacultyService(faculty *model.Faculty) error {
return s.facultyrepo.CreateFaculty(faculty)
}

func (s *FacultyService) GetFacultyService() ([]model.Faculty, error) {
return s.facultyrepo.FetchFaculty()
}

func (s *FacultyService) GetFacultyServiceById(id uint) (model.Faculty, error) {
return s.facultyrepo.FetchFacultyById(id)
}

func (s *FacultyService) DeleteFacultyService(id uint) error {
return s.facultyrepo.DeleteFaculty(id)
}

func (s *FacultyService) UpdateFacultyService(id uint, name string, gender string) error {
faculty, err := s.facultyrepo.FetchFacultyById(id)
if err != nil {
return err
}

faculty.Name = name
faculty.Gender = gender

return s.facultyrepo.UpdateFacultyById(&faculty)
}
