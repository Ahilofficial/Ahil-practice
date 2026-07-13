package services

import (
	"backend_institutions/internal/dto"
	"backend_institutions/internal/model"
	"backend_institutions/internal/repository"
	"context"
)

type StudentService struct {
	studentrepo *repository.StudentRepository
}

func NewStudentService(studentrepo *repository.StudentRepository) *StudentService {
	return &StudentService{studentrepo: studentrepo}
}

func (s *StudentService) CreateStudentService(dto *dto.CreateStudentDTO) (model.Student, error) {
	student := model.Student{
		Name:      dto.Name,
		Email:     dto.Email,
		Gender:    dto.Gender,
		FacultyID: dto.FacultyID,
	}

	err := s.studentrepo.CreateStudent(&student)
	return student, err
}

func (s *StudentService) GetStudentService(ctx context.Context) ([]dto.StudentFlatRow, error) {
	return s.studentrepo.FetchStudent(ctx)
}
func (s *StudentService) GetStudentServicePaginated(page, limit int) ([]model.Student, int64, error) {
	return s.studentrepo.FetchStudentPaginated(page, limit)
}

func (s *StudentService) GetStudentServiceById(id uint) (model.Student, error) {
	return s.studentrepo.FetchStudentById(id)
}

func (s *StudentService) GetStudentServiceDeleted() ([]model.Student, error) {
	return s.studentrepo.FetchStudentDeleted()
}

func (s *StudentService) DeleteStudentService(id uint) error {
	return s.studentrepo.DeleteStudent(id)
}

func (s *StudentService) UpdateStudentService(id uint, dto *dto.UpdateStudentDTO) error {
	student, err := s.studentrepo.FetchStudentById(id)
	if err != nil {
		return err
	}
	student.Name = dto.Name
	student.Email = dto.Email
	student.Gender = dto.Gender
	return s.studentrepo.UpdateStudentById(&student)
}

func (s *StudentService) GetActiveStudentService() (model.Student, error) {
	return s.studentrepo.GetActiveStudent()
}

func (s *StudentService) GetInactiveStudentService() (model.Student, error) {
	return s.studentrepo.GetInactiveStudent()
}
