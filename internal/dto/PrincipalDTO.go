package dto

import (
	"backend_institutions/internal/model"
	"errors"
	"strings"
	"time"

	"github.com/jinzhu/copier"
)

type CreatePrincipalDTO struct {
	Name         string    `json:"name"`
	Gender       string    `json:"gender"`
	JoiningDate  time.Time `json:"joining_date"`
	DepartmentID uint      `json:"department_id"`
	UserID       uint      `json:"user_id"`
}

func (dto *CreatePrincipalDTO) Sanitize() {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.Gender = strings.TrimSpace(strings.ToLower(dto.Gender))
}

func (dto *CreatePrincipalDTO) Validate() error {
	if dto.Name == "" {
		return errors.New("name is required")
	}
	if dto.Gender == "" {
		return errors.New("gender is required")
	}
	if dto.JoiningDate.IsZero() {
		return errors.New("joining date is required")
	}
	if dto.DepartmentID == 0 {
		return errors.New("department id is required")
	}
	return nil
}

type UpdatePrincipalDTO struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

func (dto *UpdatePrincipalDTO) Sanitize() {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.Gender = strings.TrimSpace(strings.ToLower(dto.Gender))
}

func (dto *UpdatePrincipalDTO) Validate() error {
	if dto.Name == "" {
		return errors.New("name is required")
	}
	if dto.Gender == "" {
		return errors.New("gender is required")
	}
	return nil
}

type PrincipalResponseDTO struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Gender       string    `json:"gender"`
	JoiningDate  time.Time `json:"joining_date"`
	DepartmentID uint      `json:"department_id"`
	IsActive     bool      `json:"isactive"`
}

func ToPrincipalResponseDTO(pr *model.Principal) PrincipalResponseDTO {
	var dto PrincipalResponseDTO
	copier.Copy(&dto, pr)
	return dto
}

func ToPrincipalResponseListDTO(prs []model.Principal) []PrincipalResponseDTO {
	list := make([]PrincipalResponseDTO, len(prs))
	for i := range prs {
		list[i] = ToPrincipalResponseDTO(&prs[i])
	}
	return list
}
