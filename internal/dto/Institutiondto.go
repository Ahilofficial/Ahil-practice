package dto

import (
	"backend_institutions/internal/model"
	"errors"
	"strings"
	"time"
)

type CreateInstitutionDTO struct {
	Name            string `json:"name"`
	InstitutionCode string `json:"institution_code"`
	State           string `json:"state"`
}

func (dto *CreateInstitutionDTO) Sanitize() {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.InstitutionCode = strings.TrimSpace(dto.InstitutionCode)
	dto.State = strings.TrimSpace(dto.State)
}

func (dto *CreateInstitutionDTO) Validate() error {


	if dto.Name == "" {
		return errors.New("name is required")
	}
	if dto.InstitutionCode == "" {
		return errors.New("institution code is required")
	}
	if dto.State == "" {
		return errors.New("state is required")
	}
	return nil
}

type UpdateInstitutionDTO struct {
	Name            string `json:"name"`
	InstitutionCode string `json:"institution_code"`
	State           string `json:"state"`
}

func (dto *UpdateInstitutionDTO) Sanitize() {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.InstitutionCode = strings.TrimSpace(dto.InstitutionCode)
	dto.State = strings.TrimSpace(dto.State)
}

func (dto *UpdateInstitutionDTO) Validate() error {
	

	if dto.Name == "" {
		return errors.New("name is required")
	}
	if dto.InstitutionCode == "" {
		return errors.New("institution code is required")
	}
	if dto.State == "" {
		return errors.New("state is required")
	}
	return nil
}

type InstitutionResponseDTO struct {
	ID              uint                    `json:"id"`
	Name            string                  `json:"name"`
	InstitutionCode string                  `json:"institution_code"`
	State           string                  `json:"state"`
	IsActive        bool                    `json:"isactive"`
	Departments     []DepartmentResponseDTO `json:"departments"`
}

func ToInstitutionResponseDTO(inst *model.Institutions) InstitutionResponseDTO {
	departments := make([]DepartmentResponseDTO, len(inst.Departments))

	for i, dept := range inst.Departments {
		departments[i] = ToDepartmentResponseDTO(&dept)
	}

	return InstitutionResponseDTO{
		ID:              inst.ID,
		Name:            inst.Name,
		InstitutionCode: inst.InstitutionCode,
		State:           inst.State,
		IsActive:        inst.IsActive,
		Departments:     departments,
	}
}

func ToInstitutionResponseListDTO(insts []model.Institutions) []InstitutionResponseDTO {
	list := make([]InstitutionResponseDTO, len(insts))

	for i := range insts {
		list[i] = ToInstitutionResponseDTO(&insts[i])
	}

	return list
}

type InstitutionFlatRow struct {
	InstID          uint
	InstName        string
	InstitutionCode string
	InstState       string
	InstActive      bool

	DeptID         *uint
	DepartmentName *string
	DeptActive     *bool

	FacID          *uint
	FacName        *string
	FacGender      *string
	FacJoiningDate *time.Time
	FacActive      *bool

	StudID         *uint
	StudName       *string
	StudEmail      *string
	StudGender     *string
	StudActive     *bool

	FeeID          *uint
	FeePaymentMode *string
	FeeAmount      *float64
	FeeActive      *bool
}