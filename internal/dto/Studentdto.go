package dto

import (
	"backend_institutions/internal/model"
	"errors"
	"strings"
	"github.com/jinzhu/copier"
)

type CreateStudentDTO struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	FacultyID uint   `json:"faculty_id"`
}

func (dto *CreateStudentDTO) Sanitize() {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.Email = strings.TrimSpace(strings.ToLower(dto.Email))
	dto.Gender = strings.TrimSpace(strings.ToLower(dto.Gender))
}

func (dto *CreateStudentDTO) Validate() error {


	if dto.Name == "" {
		return errors.New("name is required")
	}
	if dto.Email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(dto.Email) {
		return errors.New("invalid email format")
	}
	if dto.Gender == "" {
		return errors.New("gender is required")
	}
	if dto.FacultyID == 0 {
		return errors.New("faculty id is required")
	}
	return nil
}

type UpdateStudentDTO struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Gender string `json:"gender"`
}



func (dto *UpdateStudentDTO) Validate() error {
	dto.Sanitize()
	

	if dto.Name == "" {
		return errors.New("name is required")
	}
	if dto.Email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(dto.Email) {
		return errors.New("invalid email format")
	}
	if dto.Gender == "" {
		return errors.New("gender is required")
	}
	return nil
}

func (dto *UpdateStudentDTO) Sanitize() {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.Email = strings.TrimSpace(strings.ToLower(dto.Email))
	dto.Gender = strings.TrimSpace(strings.ToLower(dto.Gender))
}
type StudentResponseDTO struct {
	ID         uint             `json:"id"`
	Name       string           `json:"name"`
	Email      string           `json:"email"`
	Gender     string           `json:"gender"`
	FacultyID  uint             `json:"faculty_id"`
	IsActive   bool             `json:"isactive"`
	Fees        []FeesResponseDTO `json:"fees"`
}



func ToStudentResponseDTO(stud *model.Student) StudentResponseDTO {
	var dto StudentResponseDTO
	copier.Copy(&dto, stud)

	dto.Fees = make([]FeesResponseDTO, len(stud.Fees))
	for i := range stud.Fees {
		dto.Fees[i] = ToFeesResponseDTO(&stud.Fees[i])
	}

	return dto
}

func ToStudentResponseListDTO(studs []model.Student) []StudentResponseDTO {
	list := make([]StudentResponseDTO, len(studs))

	for i := range studs {
		list[i] = ToStudentResponseDTO(&studs[i])
	}

	return list
}



