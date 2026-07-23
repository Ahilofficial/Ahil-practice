package repository

import (
	"backend_institutions/internal/model"
	"errors"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) CreateRole(role *model.Role) error {
	return r.db.Exec(
		"INSERT INTO roles (name) VALUES (?)",
		role.Name,
	).Error
}

func (r *RoleRepository) GetRoleByID(id uint) (model.Role, error) {
	var role model.Role
	err := r.db.Raw("SELECT id, name, created_at, updated_at FROM roles WHERE id = ? LIMIT 1", id).Scan(&role).Error
	if err != nil {
		return role, err
	}
	if role.ID == 0 {
		return role, gorm.ErrRecordNotFound
	}
	return role, nil
}



func (r *RoleRepository) AssignPermissionsToRole(roleID uint, permissionIDs []uint, permissionNames []string) error {
	role, err := r.GetRoleByID(roleID)
	if err != nil || role.ID == 0 {
		return errors.New("role not found")
	}

	var targetIDs []uint
	if len(permissionIDs) > 0 {
		targetIDs = append(targetIDs, permissionIDs...)
	}

	if len(permissionNames) > 0 {
		var nameIDs []uint
		err := r.db.Raw("SELECT id FROM permissions WHERE name IN ?", permissionNames).Scan(&nameIDs).Error
		if err == nil {
			targetIDs = append(targetIDs, nameIDs...)
		}
	}

	if len(targetIDs) == 0 {
		return errors.New("no valid permissions found to assign")
	}

	for _, pid := range targetIDs {
		_ = r.db.Exec("INSERT IGNORE INTO role_permissions (role_id, permission_id) VALUES (?, ?)", roleID, pid)
	}

	return nil
}

func (r *RoleRepository) GetRolePermissions(roleID uint) ([]model.Permission, error) {
	var perms []model.Permission
	query := `
		SELECT p.id, p.name 
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = ?
	`
	err := r.db.Raw(query, roleID).Scan(&perms).Error
	return perms, err
}

func (r *RoleRepository) RemovePermissionFromRole(roleID uint, permissionID uint) error {
	res := r.db.Exec("DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?", roleID, permissionID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("role permission mapping not found")
	}
	return nil
}
