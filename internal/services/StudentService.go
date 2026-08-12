package services

import (
	"errors"

	"backend_institutions/internal/dto"
	"backend_institutions/internal/model"
	"backend_institutions/internal/repository"
)

type StudentService struct {
	studentRepo *repository.StudentRepository
	facultyRepo *repository.FacultyRepository
	userRepo    *repository.UserRepository
}

func NewStudentService(
	studentRepo *repository.StudentRepository,
	facultyRepo *repository.FacultyRepository,
	userRepo *repository.UserRepository,
) *StudentService {
	return &StudentService{
		studentRepo: studentRepo,
		facultyRepo: facultyRepo,
		userRepo:    userRepo,
	}
}

func (s *StudentService) checkInstitutionAccess(
	userID uint,
	institutionID uint,
) error {

	hasAccess, err := s.userRepo.HasInstitutionAccess(
		userID,
		institutionID,
	)
	if err != nil {
		return err
	}

	if !hasAccess {
		return errors.New(
			"user does not have access to this institution",
		)
	}

	return nil
}

func (s *StudentService) CreateStudentService(
	userID uint,
	student *model.Student,
) (*model.Student, error) {

	institutionID, err := s.facultyRepo.GetInstitutionByFacultyID(
		student.FacultyID,
	)
	if err != nil {
		return nil, err
	}

	if err := s.checkInstitutionAccess(
		userID,
		institutionID,
	); err != nil {
		return nil, err
	}

	if err := s.userRepo.ValidateUser(student.UserID); err != nil {
		return nil, err
	}

	exists, err := s.studentRepo.ExistsByUserID(student.UserID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("user is already a student")
	}

	if err := s.studentRepo.CreateStudent(student); err != nil {
		return nil, err
	}

	if student.UserID != 0 {
		if err := s.userRepo.UpdateUserStudentID(student.UserID, student.ID); err != nil {
			return nil, err
		}
	}

	return student, nil
}
func (s *StudentService) GetStudentService() ([]model.Student, error) {
	return s.studentRepo.FetchStudent()
}

func (s *StudentService) GetStudentServicePaginated(
	search string,
	page int,
	limit int,
) ([]model.Student, int64, error) {

	return s.studentRepo.FetchStudentPaginated(
		search,
		page,
		limit,
	)
}

func (s *StudentService) GetStudentServiceById(
	userID uint,
	id uint,
) (model.Student, error) {

	student, err := s.studentRepo.FetchStudentById(id)
	if err != nil {
		return model.Student{}, err
	}

	institutionID, err := s.facultyRepo.GetInstitutionByFacultyID(
		student.FacultyID,
	)
	if err != nil {
		return model.Student{}, err
	}

	if err := s.checkInstitutionAccess(
		userID,
		institutionID,
	); err != nil {
		return model.Student{}, err
	}

	return student, nil
}

func (s *StudentService) GetStudentServiceDeleted() ([]model.Student, error) {
	return s.studentRepo.FetchStudentDeleted()
}

func (s *StudentService) DeleteStudentService(
	userID uint,
	id uint,
) error {

	student, err := s.studentRepo.FetchStudentById(id)
	if err != nil {
		return err
	}

	institutionID, err := s.facultyRepo.GetInstitutionByFacultyID(
		student.FacultyID,
	)
	if err != nil {
		return err
	}

	if err := s.checkInstitutionAccess(
		userID,
		institutionID,
	); err != nil {
		return err
	}

	return s.studentRepo.DeleteStudent(id)
}

func (s *StudentService) UpdateStudentService(
	userID uint,
	id uint,
	req *dto.UpdateStudentDTO,
) error {

	student, err := s.studentRepo.FetchStudentById(id)
	if err != nil {
		return err
	}

	institutionID, err := s.facultyRepo.GetInstitutionByFacultyID(
		student.FacultyID,
	)
	if err != nil {
		return err
	}

	if err := s.checkInstitutionAccess(
		userID,
		institutionID,
	); err != nil {
		return err
	}

	student.Name = req.Name
	student.Email = req.Email
	student.Gender = req.Gender

	return s.studentRepo.UpdateStudentById(&student)
}

func (s *StudentService) GetActiveStudentService() (model.Student, error) {
	return s.studentRepo.GetActiveStudent()
}

func (s *StudentService) GetInactiveStudentService() (model.Student, error) {
	return s.studentRepo.GetInactiveStudent()
}

func (s *StudentService) FetchStudentsByPaymentMonth(
	month string,
) ([]model.Student, error) {

	return s.studentRepo.FetchStudentsByPaymentMonth(month)
}

func (s *StudentService) FetchPaidStudents() ([]model.Student, error) {
	return s.studentRepo.FetchPaidStudents()
}

func (s *StudentService) FetchNotPaidStudents() ([]model.Student, error) {
	return s.studentRepo.FetchNotPaidStudents()
}

func (s *StudentService) GetInstitutionID(
	studentID uint,
) (uint, error) {

	return s.studentRepo.GetInstitutionIDByStudent(studentID)
}
