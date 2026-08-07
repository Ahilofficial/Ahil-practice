package repository

import (
	"backend_institutions/internal/model"
	"errors"
	// "fmt"
	"time"

	"gorm.io/gorm"
)

type StudentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}



func (r *StudentRepository) CreateStudent(student *model.Student) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}

	now := time.Now()

	res, err := db.Exec(
		`INSERT INTO students
			(name, email, gender, faculty_id, created_at, updated_at, is_active)
		SELECT ?, ?, ?, id, ?, ?, ?
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
		return errors.New("student email already registered, or assigned faculty is inactive/invalid")
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
func (r *StudentRepository) FetchStudent() ([]model.Student, error) {
	var studs []model.Student
	err := r.db.Raw("SELECT * FROM students WHERE deleted_at IS NULL").Scan(&studs).Error
	if err != nil {
		return nil, err
	}
	
	return studs, err
}

func (r *StudentRepository) GetActiveStudent() (model.Student, error) {
	var studs []model.Student
	err := r.db.Raw("SELECT * FROM students WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", true).Scan(&studs).Error
	if err != nil {
		return model.Student{}, err
	}
	if len(studs) == 0 {
		return model.Student{}, gorm.ErrRecordNotFound
	}
	
	
	return studs[0], nil
}

func (r *StudentRepository) GetInactiveStudent() (model.Student, error) {
	var studs []model.Student
	err := r.db.Raw("SELECT * FROM students WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", false).Scan(&studs).Error
	if err != nil {
		return model.Student{}, err
	}
	if len(studs) == 0 {
		return model.Student{}, gorm.ErrRecordNotFound
	}
	
	
	return studs[0], nil
}

func (r *StudentRepository) FetchStudentPaginated(search string, page, limit int) ([]model.Student, int64, error) {
	var (
		studs []model.Student
		total int64
	)

	query := r.db.Model(&model.Student{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(`
			deleted_at IS NULL AND (
				name LIKE ? OR
				email LIKE ? OR
				gender LIKE ?
			)
		`, searchPattern, searchPattern, searchPattern)
	} else {
		query = query.Where("deleted_at IS NULL")
	}

	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	
	err := query.
		Preload("Fees").
		Preload("Fees.Payments").
		Limit(limit).
		Offset(offset).
		Find(&studs).Error
	if err != nil {
		return nil, 0, err
	}

	return studs, total, nil
}


func (r *StudentRepository) FetchStudentById(id uint) (model.Student, error) {
	var student model.Student

	err := r.db.
		Preload("Fees").
		Preload("Fees.Payments").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&student).Error

	if err != nil {
		return model.Student{}, err
	}

	return student, nil
}

func (r *StudentRepository) FetchStudentDeleted() ([]model.Student, error) {
	var studs []model.Student
	err := r.db.Raw("SELECT * FROM students WHERE deleted_at IS NOT NULL").Scan(&studs).Error
	if err != nil {
		return nil, err
	}
	
	return studs, err
}

func (r *StudentRepository) DeleteStudent(id uint) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	res, err := db.Exec(
		"UPDATE students SET is_active = ?, deleted_at = ? WHERE id = ? AND is_active = ? AND deleted_at IS NULL",
		false, time.Now(), id, true,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("record not found or already deleted")
	}
	return nil
}

func (r *StudentRepository) UpdateStudentById(student *model.Student) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	_, err = db.Exec(
		"UPDATE students SET name = ?, email = ?, gender = ?, updated_at = ? WHERE id = ?",
		student.Name, student.Email, student.Gender, time.Now(), student.ID,
	)
	return err
}


func (r *StudentRepository) FetchStudentsByPaymentMonth(month string) ([]model.Student, error) {
	var students []model.Student

	err := r.db.
		Model(&model.Student{}).
		Joins("JOIN fees ON fees.student_id = students.id").
		Joins("JOIN payments ON payments.fee_id = fees.id").
		Where("LOWER(payments.month) = LOWER(?)", month).
		Preload("Fees").
		Preload("Fees.Payments").
		Distinct().
		Find(&students).Error

	if err != nil {
		return nil, err
	}

	return students, nil
}


func (r *StudentRepository) FetchPaidStudents() ([]model.Student, error) {
	var students []model.Student

	err := r.db.
		Model(&model.Student{}).
		Joins("JOIN fees ON fees.student_id = students.id").
		Where("fees.total_amount = fees.total_paid").
		Preload("Fees").
		Preload("Fees.Payments").
		Find(&students).Error

	if err != nil {
		return nil, err
	}

	return students, nil
}

func (r *StudentRepository) FetchNotPaidStudents() ([]model.Student, error) {
	var students []model.Student

	err := r.db.Raw(`
		SELECT DISTINCT students.*
		FROM students
		INNER JOIN fees
			ON fees.student_id = students.id
		WHERE students.deleted_at IS NULL
		AND fees.total_amount <> fees.total_paid
	`).Scan(&students).Error

	if err != nil {
		return nil, err
	}

	// Load associations
	if err := r.db.
		Preload("Fees").
		Preload("Fees.Payments").
		Find(&students).Error; err != nil {
		return nil, err
	}

	return students, nil
}

func (r *StudentRepository) GetInstitutionIDByStudent(studentID uint) (uint, error) {

	var institutionID uint

	err := r.db.Raw(`
		SELECT d.institution_id
		FROM students s
		JOIN faculties f ON s.faculty_id = f.id
		JOIN departments d ON f.department_id = d.id
		WHERE s.id = ?
	`, studentID).Scan(&institutionID).Error

	if err != nil {
		return 0, err
	}

	return institutionID, nil
}
