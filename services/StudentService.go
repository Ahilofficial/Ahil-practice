package services

import (
	"backend_institutions/model"
	"backend_institutions/repository"
)


func CreateStudentService(student *model.Student) error {
	return repository.CreateStudent(student)
}

func GetStudentService() ([]model.Student, error) {
	return repository.FetchStudent()
}

func GetStudentServicePaginated(page, limit int) ([]model.Student, int64, error) {
	return repository.FetchStudentPaginated(page, limit)
}

func GetStudentServiceById(id uint) (model.Student, error) {
	return repository.FetchStudentById(id)
}

func GetStudentServiceDeleted() ([]model.Student, error) {
	return repository.FetchStudentDeleted()
}

func DeleteStudentService(id uint) error {
	return repository.DeleteStudent(id)
}

func UpdateStudentService(id uint, name string, email string, gender string) error {
	student, err := repository.FetchStudentById(id)
	if err != nil {
		return err
	}
	student.Name = name
	student.Email = email
	student.Gender = gender
	return repository.UpdateStudentById(&student)
}

func GetActiveStudentService() (model.Student, error) {
	return repository.GetActiveStudent()
}

func GetInactiveStudentService() (model.Student, error) {
	return repository.GetInactiveStudent()
}