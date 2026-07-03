package services

import (
	"backend_institutions/model"
	"backend_institutions/repository"
)



func  CreateFacultyService(faculty *model.Faculty) error {
	return repository.CreateFaculty(faculty)
}

func  GetFacultyService() ([]model.Faculty, error) {
	return repository.FetchFaculty()
}

func  GetFacultyServicePaginated(page, limit int) ([]model.Faculty, int64, error) {
	return repository.FetchFacultyPaginated(page, limit)
}

func  GetFacultyServiceById(id uint) (model.Faculty, error) {
	return repository.FetchFacultyById(id)
}

func  GetFacultyServiceDeleted() ([]model.Faculty, error) {
	return repository.FetchFacultyDeleted()
}

func  GetActiveFacultyService() (model.Faculty, error) {
	return repository.GetActiveFacutly()
}

func  GetInactiveFacultyService() (model.Faculty, error) {
	return repository.GetInactiveFaculty()
}

func  DeleteFacultyService(id uint) error {
	return repository.DeleteFaculty(id)
}

func  UpdateFacultyService(id uint, name string, gender string) error {
	faculty, err := repository.FetchFacultyById(id)
	if err != nil {
		return err
	}
	faculty.Name = name
	faculty.Gender = gender
	return repository.UpdateFacultyById(&faculty)
}
