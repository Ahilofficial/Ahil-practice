package dto

import (
	"backend_institutions/internal/model"
	"errors"
	"strings"
	"time"
)

type CreateFacultyDTO struct {
	Name         string    `json:"name"`
	Gender       string    `json:"gender"`
	JoiningDate  time.Time `json:"joining_date"`
	DepartmentID uint      `json:"department_id"`
}

func (dto *CreateFacultyDTO) Sanitize() {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.Gender = strings.TrimSpace(strings.ToLower(dto.Gender))
}

func (dto *CreateFacultyDTO) Validate() error {
	

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

type UpdateFacultyDTO struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

func (dto *UpdateFacultyDTO) Sanitize() {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.Gender = strings.TrimSpace(strings.ToLower(dto.Gender))
}

func (dto *UpdateFacultyDTO) Validate() error {
	
	if dto.Name == "" {
		return errors.New("name is required")
	}
	if dto.Gender == "" {
		return errors.New("gender is required")
	}
	return nil
}

type FacultyResponseDTO struct {
	ID           uint                 `json:"id"`
	Name         string               `json:"name"`
	Gender       string               `json:"gender"`
	JoiningDate  time.Time            `json:"joining_date"`
	DepartmentID uint                 `json:"department_id"`
	IsActive     bool                 `json:"isactive"`
	Students     []StudentResponseDTO `json:"students"`
}

func ToFacultyResponseDTO(fac *model.Faculty) FacultyResponseDTO {
	students := make([]StudentResponseDTO, len(fac.Students))

	for i, student := range fac.Students {
		students[i] = ToStudentResponseDTO(&student)
	}

	return FacultyResponseDTO{
		ID:           fac.ID,
		Name:         fac.Name,
		Gender:       fac.Gender,
		JoiningDate:  fac.JoiningDate,
		DepartmentID: fac.DepartmentID,
		IsActive:     fac.IsActive,
		Students:     students,
	}
}

func ToFacultyResponseListDTO(facs []model.Faculty) []FacultyResponseDTO {
	list := make([]FacultyResponseDTO, len(facs))

	for i := range facs {
		list[i] = ToFacultyResponseDTO(&facs[i])
	}

	return list
}


type FacultyFlatRow struct {
	FacID          uint
	FacName        string
	FacGender      string
	FacJoiningDate time.Time
	DepartmentID   uint
	FacActive      bool

	StudID     *uint
	StudName   *string
	StudEmail  *string
	StudGender *string
	StudActive *bool

	FeeID          *uint
	FeePaymentMode *string
	FeeAmount      *float64
	FeeActive      *bool
}