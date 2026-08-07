package repository

import (
	"backend_institutions/internal/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type FacultyRepository struct {
	db *gorm.DB
}

func NewFacultyRepository(db *gorm.DB) *FacultyRepository {
	return &FacultyRepository{
		db: db,
	}
}

func (r *FacultyRepository) CreateFaculty(faculty *model.Faculty) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}

	now := time.Now()

	res, err := db.Exec(
		`INSERT INTO faculties
			(name, gender, joining_date, department_id, created_at, updated_at, is_active)
		SELECT ?, ?, ?, id, ?, ?, ?
		FROM departments
		WHERE id = ?
		  AND deleted_at IS NULL
		  AND is_active = true
		  AND NOT EXISTS (
			  SELECT 1
			  FROM faculties
			  WHERE name = ?
			    AND department_id = ?
			    AND deleted_at IS NULL
		  )`,
		faculty.Name,
		faculty.Gender,
		faculty.JoiningDate,
		now,
		now,
		true,
		faculty.DepartmentID,
		faculty.Name,
		faculty.DepartmentID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("faculty name already exists in this department, or parent department is inactive/invalid")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	faculty.ID = uint(id)
	faculty.CreatedAt = now
	faculty.UpdatedAt = now
	faculty.IsActive = true

	return nil
}

func (r *FacultyRepository) FetchFaculty() ([]model.Faculty, error) {
	var facs []model.Faculty
	err := r.db.Raw("SELECT * FROM faculties WHERE deleted_at IS NULL").Scan(&facs).Error
	if err != nil {
		return nil, err
	}

	return facs, err
}

func (r *FacultyRepository) FetchFacultyPaginated(search string, page, limit int) ([]model.Faculty, int64, error) {
	var (
		facs  []model.Faculty
		total int64
	)

	query := r.db.Model(&model.Faculty{})

	// Search
	if search != "" {
		search = "%" + search + "%"
		query = query.Where(`
			name LIKE ? OR
			gender LIKE ? OR
			joining_date LIKE ?
		`, search, search, search)
	}

	// Count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	// Fetch with associations
	err := query.
		Preload("Students").
		Preload("Students.Fees").
		Preload("Students.Fees.Payments").
		Limit(limit).
		Offset(offset).
		Find(&facs).Error

	if err != nil {
		return nil, 0, err
	}

	return facs, total, nil
}

func (r *FacultyRepository) FetchFacultyById(id uint) (model.Faculty, error) {
	var fac model.Faculty
	err := r.db.
		Preload("Students").
		Preload("Students.Fees").
		Preload("Students.Fees.Payments").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&fac).Error
	if err != nil {
		return model.Faculty{}, err
	}

	return fac, nil
}

func (r *FacultyRepository) FetchFacultyDeleted() ([]model.Faculty, error) {
	var facs []model.Faculty
	err := r.db.Raw("SELECT * FROM faculties WHERE deleted_at IS NOT NULL").Scan(&facs).Error
	if err != nil {
		return nil, err
	}
	return facs, err
}

func (r *FacultyRepository) GetActiveFaculty() (model.Faculty, error) {
	var facs []model.Faculty
	err := r.db.Raw("SELECT * FROM faculties WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", true).Scan(&facs).Error
	if err != nil {
		return model.Faculty{}, err
	}
	if len(facs) == 0 {
		return model.Faculty{}, gorm.ErrRecordNotFound
	}

	return facs[0], nil
}

func (r *FacultyRepository) GetInactiveFaculty() (model.Faculty, error) {
	var facs []model.Faculty
	err := r.db.Raw("SELECT * FROM faculties WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", false).Scan(&facs).Error
	if err != nil {
		return model.Faculty{}, err
	}
	if len(facs) == 0 {
		return model.Faculty{}, gorm.ErrRecordNotFound
	}

	return facs[0], nil
}

func (r *FacultyRepository) DeleteFaculty(id uint) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	res, err := db.Exec(
		"UPDATE faculties SET is_active = ?, deleted_at = ? WHERE id = ? AND is_active = ? AND deleted_at IS NULL",
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

func (r *FacultyRepository) UpdateFacultyById(faculty *model.Faculty) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	_, err = db.Exec(
		"UPDATE faculties SET name = ?, gender = ?, updated_at = ? WHERE id = ?",
		faculty.Name, faculty.Gender, time.Now(), faculty.ID,
	)
	return err
}
