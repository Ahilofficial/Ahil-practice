package dto

import (
	"backend_institutions/internal/model"
	"errors"
	"strings"
)

type CreateFeesDTO struct {
	PaymentMode string  `json:"payment_mode"`
	Amount      float64 `json:"amount"`
	StudentID   uint    `json:"student_id"`
}

func (dto *CreateFeesDTO) Sanitize() {
	dto.PaymentMode = strings.TrimSpace(strings.ToLower(dto.PaymentMode))
}

func (dto *CreateFeesDTO) Validate() error {
	if dto.PaymentMode == "" {
		return errors.New("payment mode is required")
	}
	if dto.Amount == 0 {
		return errors.New("amount is required and must be greater than 0")
	}
	if dto.StudentID == 0 {
		return errors.New("student id is required")
	}
	return nil
}

type UpdateFeesDTO struct {
	PaymentMode string  `json:"payment_mode"`
	Amount      float64 `json:"amount"`
}

func (dto *UpdateFeesDTO) Sanitize() {
	dto.PaymentMode = strings.TrimSpace(strings.ToLower(dto.PaymentMode))
}

func (dto *UpdateFeesDTO) Validate() error {
	if dto.PaymentMode == "" {
		return errors.New("payment mode is required")
	}
	if dto.Amount == 0 {
		return errors.New("amount is required and must be greater than 0")
	}
	return nil
}

type FeesResponseDTO struct {
	ID          uint    `json:"id"`
	PaymentMode string  `json:"payment_mode"`
	Amount      float64 `json:"amount"`
	StudentID   uint    `json:"student_id"`
	IsActive    bool    `json:"isactive"`
}

func ToFeesResponseDTO(fees *model.Fees) FeesResponseDTO {
	return FeesResponseDTO{
		ID:          fees.ID,
		PaymentMode: fees.PaymentMode,
		Amount:      fees.Amount,
		StudentID:   fees.StudentID,
		IsActive:    fees.IsActive,
	}
}

func ToFeesResponseListDTO(fees []model.Fees) []FeesResponseDTO {
	list := make([]FeesResponseDTO, len(fees))
	for i, f := range fees {
		list[i] = ToFeesResponseDTO(&f)
	}
	return list
}
