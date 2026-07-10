package repository

import (
	// "backend_institutions/internal/dto"
	"backend_institutions/internal/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) CreateDepartment(department *model.Department) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	now := time.Now()
	res, err := db.Exec(
		`INSERT INTO departments (department_name, institution_id, created_at, updated_at, is_active)
		SELECT ?, id, ?, ?, ? FROM institutions
		WHERE id = ? AND deleted_at IS NULL AND is_active = true
		  AND NOT EXISTS (
			  SELECT 1 FROM departments 
			  WHERE department_name = ? AND institution_id = ? AND deleted_at IS NULL
		  )`,
		department.DepartmentName, now, now, true,
		department.InstitutionID,
		department.DepartmentName, department.InstitutionID,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("department name already exists in this institution, or parent institution is inactive/invalid")
	}
	id, err := res.LastInsertId()
	if err == nil {
		department.ID = uint(id)
	}
	return nil
}

func (r *DepartmentRepository) FetchDepartment() ([]model.Department, error) {
	var depts []model.Department
	err := r.db.Raw("SELECT * FROM departments WHERE deleted_at IS NULL").Scan(&depts).Error
	return depts, err
}

func (r *DepartmentRepository) FetchDepartmentPaginated(page, limit int) ([]model.Department, int64, error) {
	var total int64
	err := r.db.Raw("SELECT COUNT(*) FROM departments WHERE deleted_at IS NULL").Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var depts []model.Department
	err = r.db.Raw("SELECT * FROM departments WHERE deleted_at IS NULL LIMIT ? OFFSET ?", limit, offset).Scan(&depts).Error
	return depts, total, err
}

func (r *DepartmentRepository) FetchDepartmentById(id uint) (model.Department, error) {
	var depts []model.Department
	err := r.db.Raw("SELECT * FROM departments WHERE id = ? AND deleted_at IS NULL LIMIT 1", id).Scan(&depts).Error
	if err != nil {
		return model.Department{}, err
	}
	if len(depts) == 0 {
		return model.Department{}, gorm.ErrRecordNotFound
	}
	return depts[0], nil
}

func (r *DepartmentRepository) FetchDepartmentDeleted() ([]model.Department, error) {
	var depts []model.Department
	err := r.db.Raw("SELECT * FROM departments WHERE deleted_at IS NOT NULL").Scan(&depts).Error
	return depts, err
}

func (r *DepartmentRepository) GetActiveDepartment() (model.Department, error) {
	var depts []model.Department
	err := r.db.Raw("SELECT * FROM departments WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", true).Scan(&depts).Error
	if err != nil {
		return model.Department{}, err
	}
	if len(depts) == 0 {
		return model.Department{}, gorm.ErrRecordNotFound
	}
	return depts[0], nil
}

func (r *DepartmentRepository) GetInactiveDepartment() (model.Department, error) {
	var depts []model.Department
	err := r.db.Raw("SELECT * FROM departments WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", false).Scan(&depts).Error
	if err != nil {
		return model.Department{}, err
	}
	if len(depts) == 0 {
		return model.Department{}, gorm.ErrRecordNotFound
	}
	return depts[0], nil
}

func (r *DepartmentRepository) DeleteDepartment(id uint) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	res, err := db.Exec(
		"UPDATE departments SET is_active = ?, deleted_at = ? WHERE id = ? AND is_active = ? AND deleted_at IS NULL",
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

func (r *DepartmentRepository) UpdateDepartmentById(department *model.Department) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	_, err = db.Exec(
		"UPDATE departments SET department_name = ?, updated_at = ? WHERE id = ?",
		department.DepartmentName, time.Now(), department.ID,
	)
	return err
}
