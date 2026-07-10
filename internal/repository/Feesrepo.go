package repository

import (
	"backend_institutions/internal/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type FeesRepository struct {
	db *gorm.DB
}

func NewFeesRepository(db *gorm.DB) *FeesRepository {
	return &FeesRepository{db: db}
}

func (r *FeesRepository) CreateFees(fees *model.Fees) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	now := time.Now()
	res, err := db.Exec(
		"INSERT INTO fees (payment_mode, amount, student_id, created_at, updated_at, is_active) VALUES (?, ?, ?, ?, ?, ?)",
		fees.PaymentMode, fees.Amount, fees.StudentID, now, now, true,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		fees.ID = uint(id)
	}
	return nil
}

func (r *FeesRepository) FetchFees() ([]model.Fees, error) {
	var fees []model.Fees
	err := r.db.Raw("SELECT * FROM fees WHERE deleted_at IS NULL").Scan(&fees).Error
	return fees, err
}

func (r *FeesRepository) FetchFeesPaginated(page, limit int) ([]model.Fees, int64, error) {
	var fees []model.Fees
	var total int64

	err := r.db.Raw("SELECT COUNT(*) FROM fees WHERE deleted_at IS NULL").Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = r.db.Raw("SELECT * FROM fees WHERE deleted_at IS NULL LIMIT ? OFFSET ?", limit, offset).Scan(&fees).Error
	return fees, total, err
}

func (r *FeesRepository) FetchFeesById(id uint) (model.Fees, error) {
	var fees model.Fees
	err := r.db.Raw("SELECT * FROM fees WHERE id = ? AND deleted_at IS NULL LIMIT 1", id).Scan(&fees).Error
	if err != nil {
		return fees, err
	}
	if fees.ID == 0 {
		return fees, gorm.ErrRecordNotFound
	}
	return fees, nil
}

func (r *FeesRepository) DeleteFees(id uint) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	res, err := db.Exec(
		"UPDATE fees SET is_active = ?, deleted_at = ? WHERE id = ? AND is_active = ? AND deleted_at IS NULL",
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

func (r *FeesRepository) FetchInactiveFees() ([]model.Fees, error) {
	var fees []model.Fees
	err := r.db.Raw("SELECT * FROM fees WHERE is_active = ? AND deleted_at IS NULL", false).Scan(&fees).Error
	return fees, err
}

func (r *FeesRepository) UpdateFeesById(fees *model.Fees) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	_, err = db.Exec(
		"UPDATE fees SET payment_mode = ?, amount = ?, updated_at = ? WHERE id = ?",
		fees.PaymentMode, fees.Amount, time.Now(), fees.ID,
	)
	return err
}
