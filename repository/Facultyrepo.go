package repository

import (
"backend_institutions/database"
"backend_institutions/model"
)

type FacultyRepository struct{}

func NewFacultyRepository() *FacultyRepository {
return &FacultyRepository{}
}

func (r *FacultyRepository) CreateFaculty(faculty *model.Faculty) error {
return database.DB.Create(faculty).Error
}

func (r *FacultyRepository) FetchFaculty() ([]model.Faculty, error) {
var faculty []model.Faculty
err := database.DB.Find(&faculty).Error
return faculty, err
}

func (r *FacultyRepository) FetchFacultyById(id uint) (model.Faculty, error) {
var faculty model.Faculty
err := database.DB.First(&faculty, id).Error
return faculty, err
}

func (r *FacultyRepository) DeleteFaculty(id uint) error {
return database.DB.Delete(&model.Faculty{}, id).Error
}

func (r *FacultyRepository) UpdateFacultyById(faculty *model.Faculty) error {
return database.DB.Save(faculty).Error
}
