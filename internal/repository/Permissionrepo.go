package repository

import (
	"backend_institutions/internal/model"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) GetAllPermissions() ([]model.Permission, error) {
	var perms []model.Permission
	err := r.db.Raw("SELECT id, name FROM permissions ORDER BY id ASC").Scan(&perms).Error
	return perms, err
}
