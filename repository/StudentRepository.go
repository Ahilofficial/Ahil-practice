package repository

import (
"backend_institutions/database"
"backend_institutions/model"
)

type StudentRepository struct{}

func NewStudentRepository() *StudentRepository {
return &StudentRepository{}
}

func (r *StudentRepository) CreateStudent(student *model.Student) error {
return database.DB.Create(student).Error
}

func (r *StudentRepository) FetchStudent() ([]model.Student, error) {
var student []model.Student
err := database.DB.Find(&student).Error
return student, err
}

func (r *StudentRepository) FetchStudentById(id uint) (model.Student, error) {
var student model.Student
err := database.DB.First(&student, id).Error
return student, err
}

func (r *StudentRepository) DeleteStudent(id uint) error {
return database.DB.Delete(&model.Student{}, id).Error
}

func (r *StudentRepository) UpdateStudentById(student *model.Student) error {
return database.DB.Save(student).Error
}
