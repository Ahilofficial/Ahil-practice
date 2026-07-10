package repository

import (
	"backend_institutions/internal/model"
	"errors"
	"time"
	// "backend_institutions/internal/dto"
	"gorm.io/gorm"
)

type InstitutionRepository struct {
	db *gorm.DB
}

func NewInstitutionRepository(db *gorm.DB) *InstitutionRepository {
	return &InstitutionRepository{
		db: db,
	}
}

func (r *InstitutionRepository) CreateInstitution(institute *model.Institutions) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	now := time.Now()
	res, err := db.Exec(
		`INSERT INTO institutions (name, institution_code, state, created_at, updated_at, is_active)
		SELECT ?, ?, ?, ?, ?, ? FROM DUAL
		WHERE NOT EXISTS (
			SELECT 1 FROM institutions 
			WHERE (name = ? OR institution_code = ?) AND deleted_at IS NULL
		)`,
		institute.Name, institute.InstitutionCode, institute.State, now, now, true,
		institute.Name, institute.InstitutionCode,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("institution name or code already exists")
	}
	id, err := res.LastInsertId()
	if err == nil {
		institute.ID = uint(id)
	}
	return nil
}

func (r *InstitutionRepository) FetchInstitution() ([]model.Institutions, error) {
var insts []model.Institutions
	err := r.db.Raw("SELECT * FROM institutions WHERE deleted_at IS NULL").Scan(&insts).Error
	return insts, err
}

func (r *InstitutionRepository) FetchInstitutionPaginated(page, limit int) ([]model.Institutions, int64, error) {
	var total int64
	err := r.db.Raw("SELECT COUNT(*) FROM institutions WHERE deleted_at IS NULL").Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var insts []model.Institutions
	err = r.db.Raw("SELECT * FROM institutions WHERE deleted_at IS NULL LIMIT ? OFFSET ?", limit, offset).Scan(&insts).Error
	return insts, total, err
}
	


func (r *InstitutionRepository) FetchInstitutionById(id uint) (model.Institutions, error) {
	var insts []model.Institutions
	err := r.db.Raw("SELECT * FROM institutions WHERE id = ? AND deleted_at IS NULL LIMIT 1", id).Scan(&insts).Error
	if err != nil {
		return model.Institutions{}, err
	}
	if len(insts) == 0 {
		return model.Institutions{}, gorm.ErrRecordNotFound
	}
	return insts[0], nil
}

func (r *InstitutionRepository) GetActiveInstitute() (model.Institutions, error) {
	var insts []model.Institutions
	err := r.db.Raw("SELECT * FROM institutions WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", true).Scan(&insts).Error
	if err != nil {
		return model.Institutions{}, err
	}
	if len(insts) == 0 {
		return model.Institutions{}, gorm.ErrRecordNotFound
	}
	return insts[0], nil
}

func (r *InstitutionRepository) GetInactiveInstitute() (model.Institutions, error) {
	var insts []model.Institutions
	err := r.db.Raw("SELECT * FROM institutions WHERE is_active = ? AND deleted_at IS NULL LIMIT 1", false).Scan(&insts).Error
	if err != nil {
		return model.Institutions{}, err
	}
	if len(insts) == 0 {
		return model.Institutions{}, gorm.ErrRecordNotFound
	}
	return insts[0], nil
}

func (r *InstitutionRepository) FetchInstitutionDeleted() ([]model.Institutions, error) {
	var insts []model.Institutions
	err := r.db.Raw("SELECT * FROM institutions WHERE deleted_at IS NOT NULL").Scan(&insts).Error
	return insts, err
}

func (r *InstitutionRepository) DeleteInstitution(id uint) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	res, err := db.Exec(
		"UPDATE institutions SET is_active = ?, deleted_at = ? WHERE id = ? AND is_active = ? AND deleted_at IS NULL",
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

func (r *InstitutionRepository) UpdateInstitution(institute *model.Institutions) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	_, err = db.Exec(
		"UPDATE institutions SET name = ?, institution_code = ?, state = ?, updated_at = ? WHERE id = ?",
		institute.Name, institute.InstitutionCode, institute.State, time.Now(), institute.ID,
	)
	return err
}
