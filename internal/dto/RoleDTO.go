package dto

import (
	"backend_institutions/internal/model"
	"errors"
	"strings"
	"time"

	"github.com/jinzhu/copier"
)

type CreateRoleDTO struct {
	Name string `json:"name"`
}

type UpdateRoleDTO struct {
	Name string `json:"name"`
}

type AssignPermissionsDTO struct {
	PermissionIDs   []uint   `json:"permission_ids,omitempty"`
	PermissionNames []string `json:"permission_names,omitempty"`
}

type RoleResponseDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (dto *CreateRoleDTO) Sanitize() {
	dto.Name = strings.TrimSpace(strings.ToLower(dto.Name))
}

func (dto *CreateRoleDTO) Validate() error {
	if dto.Name == "" {
		return errors.New("role name is required")
	}
	return nil
}

func (dto *UpdateRoleDTO) Sanitize() {
	dto.Name = strings.TrimSpace(strings.ToLower(dto.Name))
}

func (dto *UpdateRoleDTO) Validate() error {
	if dto.Name == "" {
		return errors.New("role name is required")
	}
	return nil
}

func (dto *AssignPermissionsDTO) Validate() error {
	if len(dto.PermissionIDs) == 0 && len(dto.PermissionNames) == 0 {
		return errors.New("either permission_ids or permission_names must be provided")
	}
	return nil
}

func ToRoleResponseDTO(role *model.Role) RoleResponseDTO {
	var dto RoleResponseDTO
	copier.Copy(&dto, role)
	return dto
}

func ToRoleResponseListDTO(roles []model.Role) []RoleResponseDTO {
	list := make([]RoleResponseDTO, len(roles))
	for i, r := range roles {
		list[i] = ToRoleResponseDTO(&r)
	}
	return list
}
