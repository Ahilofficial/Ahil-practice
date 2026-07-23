package dto

import (
	"backend_institutions/internal/model"

	"github.com/jinzhu/copier"
)

type PermissionResponseDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func ToPermissionResponseDTO(perm *model.Permission) PermissionResponseDTO {
	var dto PermissionResponseDTO
	copier.Copy(&dto, perm)
	return dto
}

func ToPermissionResponseListDTO(perms []model.Permission) []PermissionResponseDTO {
	list := make([]PermissionResponseDTO, len(perms))
	for i, p := range perms {
		list[i] = ToPermissionResponseDTO(&p)
	}
	return list
}
