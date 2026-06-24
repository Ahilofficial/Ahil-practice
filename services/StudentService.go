package services

import (
"backend_institutions/model"
"backend_institutions/repository"
)

type StudentService struct {
studentrepo *repository.StudentRepository
}

func NewStudentService() *StudentService {
return &StudentService{
studentrepo: repository.NewStudentRepository(),
}
}

func (s *StudentService) CreateStudentService(student *model.Student) error {
return s.studentrepo.CreateStudent(student)
}

func (s *StudentService) GetStudentService() ([]model.Student, error) {
return s.studentrepo.FetchStudent()
}

func (s *StudentService) GetStudentServiceById(id uint) (model.Student, error) {
return s.studentrepo.FetchStudentById(id)
}

func (s *StudentService) DeleteStudentService(id uint) error {
return s.studentrepo.DeleteStudent(id)
}

func (s *StudentService) UpdateStudentService(id uint, name string, email string, gender string) error {
student, err := s.studentrepo.FetchStudentById(id)
if err != nil {
return err
}

student.Name = name
student.Email = email
student.Gender = gender

return s.studentrepo.UpdateStudentById(&student)
}
