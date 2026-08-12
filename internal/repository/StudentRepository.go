package repository

import (
	"errors"
	"time"

	"backend_institutions/internal/model"

	"gorm.io/gorm"
)

type StudentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{
		db: db,
	}
}

// --------------------------------------------------
// Create Student
// --------------------------------------------------

func (r *StudentRepository) CreateStudent(
	student *model.Student,
) error {

	db, err := r.db.DB()
	if err != nil {
		return err
	}

	now := time.Now()

	res, err := db.Exec(
		`INSERT INTO students
			(name, email, gender, faculty_id, user_id, created_at, updated_at, is_active)
		SELECT ?, ?, ?, id, ?, ?, ?, ?
		FROM faculties
		WHERE id = ?
		  AND deleted_at IS NULL
		  AND is_active = true
		  AND NOT EXISTS (
			  SELECT 1
			  FROM students
			  WHERE email = ?
			    AND deleted_at IS NULL
		  )`,
		student.Name,
		student.Email,
		student.Gender,
		student.UserID,
		now,
		now,
		true,
		student.FacultyID,
		student.Email,
	)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New(
			"student email already registered, or assigned faculty is inactive/invalid",
		)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	student.ID = uint(id)
	student.CreatedAt = now
	student.UpdatedAt = now
	student.IsActive = true

	return nil
}

// --------------------------------------------------
// Check Student By User ID
// --------------------------------------------------

func (r *StudentRepository) ExistsByUserID(
	userID uint,
) (bool, error) {

	var exists bool

	result := r.db.Raw(`
		SELECT EXISTS(
			SELECT 1
			FROM students
			WHERE user_id = ?
			  AND deleted_at IS NULL
		)
	`, userID).Scan(&exists)

	if result.Error != nil {
		return false, result.Error
	}

	return exists, nil
}

// --------------------------------------------------
// Get All Students
// --------------------------------------------------

func (r *StudentRepository) FetchStudent() ([]model.Student, error) {

	var students []model.Student

	err := r.db.
		Where("deleted_at IS NULL").
		Preload("Faculty").
		Preload("Fees").
		Preload("Fees.Payments").
		Find(&students).Error

	if err != nil {
		return nil, err
	}

	return students, nil
}

// --------------------------------------------------
// Get Paginated Students
// --------------------------------------------------

func (r *StudentRepository) FetchStudentPaginated(
	search string,
	page int,
	limit int,
) ([]model.Student, int64, error) {

	var (
		students []model.Student
		total    int64
	)

	query := r.db.
		Model(&model.Student{}).
		Where("students.deleted_at IS NULL")

	if search != "" {

		searchPattern := "%" + search + "%"

		query = query.Where(`
			(
				students.name LIKE ?
				OR students.email LIKE ?
				OR students.gender LIKE ?
			)
		`,
			searchPattern,
			searchPattern,
			searchPattern,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	err := query.
		Preload("Faculty").
		Preload("Fees").
		Preload("Fees.Payments").
		Limit(limit).
		Offset(offset).
		Find(&students).Error

	if err != nil {
		return nil, 0, err
	}

	return students, total, nil
}

// --------------------------------------------------
// Get Student By ID
// --------------------------------------------------

func (r *StudentRepository) FetchStudentById(
	id uint,
) (model.Student, error) {

	var student model.Student

	err := r.db.
		Where("id = ? AND deleted_at IS NULL", id).
		Preload("Faculty").
		Preload("Fees").
		Preload("Fees.Payments").
		First(&student).Error

	if err != nil {
		return model.Student{}, err
	}

	return student, nil
}

// --------------------------------------------------
// Get Deleted Students
// --------------------------------------------------

func (r *StudentRepository) FetchStudentDeleted() ([]model.Student, error) {

	var students []model.Student

	err := r.db.
		Unscoped().
		Where("deleted_at IS NOT NULL").
		Find(&students).Error

	if err != nil {
		return nil, err
	}

	return students, nil
}

// --------------------------------------------------
// Get Active Student
// --------------------------------------------------

func (r *StudentRepository) GetActiveStudent() (model.Student, error) {

	var student model.Student

	err := r.db.
		Where("is_active = ? AND deleted_at IS NULL", true).
		Preload("Faculty").
		Preload("Fees").
		Preload("Fees.Payments").
		First(&student).Error

	if err != nil {
		return model.Student{}, err
	}

	return student, nil
}

// --------------------------------------------------
// Get Inactive Student
// --------------------------------------------------

func (r *StudentRepository) GetInactiveStudent() (model.Student, error) {

	var student model.Student

	err := r.db.
		Where("is_active = ? AND deleted_at IS NULL", false).
		Preload("Faculty").
		Preload("Fees").
		Preload("Fees.Payments").
		First(&student).Error

	if err != nil {
		return model.Student{}, err
	}

	return student, nil
}

// --------------------------------------------------
// Delete Student
// --------------------------------------------------

func (r *StudentRepository) DeleteStudent(id uint) error {

	db, err := r.db.DB()
	if err != nil {
		return err
	}

	now := time.Now()

	res, err := db.Exec(
		`UPDATE students
		 SET is_active = ?,
		     deleted_at = ?
		 WHERE id = ?
		   AND is_active = ?
		   AND deleted_at IS NULL`,
		false,
		now,
		id,
		true,
	)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New(
			"record not found or already deleted",
		)
	}

	return nil
}

// --------------------------------------------------
// Update Student
// --------------------------------------------------

func (r *StudentRepository) UpdateStudentById(
	student *model.Student,
) error {

	db, err := r.db.DB()
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`UPDATE students
		 SET name = ?,
		     email = ?,
		     gender = ?,
		     updated_at = ?
		 WHERE id = ?
		   AND deleted_at IS NULL`,
		student.Name,
		student.Email,
		student.Gender,
		time.Now(),
		student.ID,
	)

	return err
}

// --------------------------------------------------
// Students By Payment Month
// --------------------------------------------------

func (r *StudentRepository) FetchStudentsByPaymentMonth(
	month string,
) ([]model.Student, error) {

	var students []model.Student

	err := r.db.
		Model(&model.Student{}).
		Joins("JOIN fees ON fees.student_id = students.id").
		Joins("JOIN payments ON payments.fee_id = fees.id").
		Where(`
			students.deleted_at IS NULL
			AND LOWER(payments.month) = LOWER(?)
		`, month).
		Preload("Faculty").
		Preload("Fees").
		Preload("Fees.Payments").
		Distinct().
		Find(&students).Error

	if err != nil {
		return nil, err
	}

	return students, nil
}

// --------------------------------------------------
// Paid Students
// --------------------------------------------------

func (r *StudentRepository) FetchPaidStudents() ([]model.Student, error) {

	var students []model.Student

	err := r.db.
		Model(&model.Student{}).
		Joins("JOIN fees ON fees.student_id = students.id").
		Where(`
			students.deleted_at IS NULL
			AND fees.total_amount = fees.total_paid
		`).
		Preload("Faculty").
		Preload("Fees").
		Preload("Fees.Payments").
		Distinct().
		Find(&students).Error

	if err != nil {
		return nil, err
	}

	return students, nil
}

// --------------------------------------------------
// Not Paid Students
// --------------------------------------------------

func (r *StudentRepository) FetchNotPaidStudents() ([]model.Student, error) {

	var students []model.Student

	err := r.db.
		Model(&model.Student{}).
		Joins("JOIN fees ON fees.student_id = students.id").
		Where(`
			students.deleted_at IS NULL
			AND fees.total_amount <> fees.total_paid
		`).
		Preload("Faculty").
		Preload("Fees").
		Preload("Fees.Payments").
		Distinct().
		Find(&students).Error

	if err != nil {
		return nil, err
	}

	return students, nil
}

// --------------------------------------------------
// Get Institution ID By Student
// --------------------------------------------------

func (r *StudentRepository) GetInstitutionIDByStudent(
	studentID uint,
) (uint, error) {

	var institutionID uint

	err := r.db.Raw(`
		SELECT d.institution_id
		FROM students s
		JOIN faculties f
			ON s.faculty_id = f.id
		JOIN departments d
			ON f.department_id = d.id
		WHERE s.id = ?
		  AND s.deleted_at IS NULL
	`, studentID).Scan(&institutionID).Error

	if err != nil {
		return 0, err
	}

	return institutionID, nil
}

func (r *StudentRepository) GetInstitutionByStudentID(studentID uint) (uint, error) {
	return r.GetInstitutionIDByStudent(studentID)
}
