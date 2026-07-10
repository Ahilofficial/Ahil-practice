package repository

import (
	"backend_institutions/internal/model"
	"errors"
	"time"
	// "backend_institutions/internal/dto"
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
		`INSERT INTO students (name, email, gender, faculty_id, created_at, updated_at, is_active)
		SELECT ?, ?, ?, id, ?, ?, ? FROM faculties
		WHERE id = ? AND deleted_at IS NULL AND is_active = true
		  AND NOT EXISTS (
			  SELECT 1 FROM students 
			  WHERE email = ? AND deleted_at IS NULL
		  )`,
		student.Name, student.Email, student.Gender, now, now, true,
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
	if err == nil {
		student.ID = uint(id)
	}
	return nil
}

func (r *StudentRepository) FetchStudent() ([]model.Student, error) {
	var studs []model.Student
	err := r.db.Raw("SELECT * FROM students WHERE deleted_at IS NULL").Scan(&studs).Error
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

func (r *StudentRepository) FetchStudentPaginated(page, limit int) ([]model.Student, int64, error) {
	var total int64
	err := r.db.Raw("SELECT COUNT(*) FROM students WHERE deleted_at IS NULL").Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var studs []model.Student
	err = r.db.Raw("SELECT * FROM students WHERE deleted_at IS NULL LIMIT ? OFFSET ?", limit, offset).Scan(&studs).Error
	return studs, total, err
}

func (r *StudentRepository) FetchStudentById(id uint) (model.Student, error) {
	var studs []model.Student
	err := r.db.Raw("SELECT * FROM students WHERE id = ? AND deleted_at IS NULL LIMIT 1", id).Scan(&studs).Error
	if err != nil {
		return model.Student{}, err
	}
	if len(studs) == 0 {
		return model.Student{}, gorm.ErrRecordNotFound
	}
	return studs[0], nil
}

func (r *StudentRepository) FetchStudentDeleted() ([]model.Student, error) {
	var studs []model.Student
	err := r.db.Raw("SELECT * FROM students WHERE deleted_at IS NOT NULL").Scan(&studs).Error
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
