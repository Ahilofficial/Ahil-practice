package dto

import (
	"backend_institutions/internal/model"
	"errors"
	"strings"
	"github.com/jinzhu/copier"
)

type CreateFeesDTO struct {
	PaymentMode string  `json:"payment_mode"`
	TotalAmount float64 `json:"total_amount"`
	StudentID   uint    `json:"student_id"`
}

func (dto *CreateFeesDTO) Sanitize() {
	dto.PaymentMode = strings.TrimSpace(strings.ToLower(dto.PaymentMode))
}

func (dto *CreateFeesDTO) Validate() error {
	if dto.PaymentMode == "" {
		return errors.New("payment mode is required")
	}
	if dto.TotalAmount <= 0 {
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

type PaymentResponseDTO struct {
	ID          uint    `json:"id"`
	Month       string  `json:"month"`
	AmountPaid  float64 `json:"amount_paid"`
	PaymentMode string  `json:"payment_mode"`
}
type FeesResponseDTO struct {
	ID             uint                 `json:"id"`
	TotalAmount    float64              `json:"total_amount"`
	TotalPaid      float64              `json:"total_paid"`
	PendingAmount  float64              `json:"pending_amount"`
	StudentID      uint                 `json:"student_id"`
	IsActive       bool                 `json:"is_active"`
	Payments       []PaymentResponseDTO `json:"payments"`
}


func ToFeesResponseDTO(fees *model.Fees) FeesResponseDTO {
	var dto FeesResponseDTO
	copier.Copy(&dto, fees)
	return dto
}

func ToFeesResponseListDTO(fees []model.Fees) []FeesResponseDTO {
	list := make([]FeesResponseDTO, len(fees))
	for i, f := range fees {
		list[i] = ToFeesResponseDTO(&f)
	}
	return list
}

type CreatePaymentDTO struct {
	Month       string  `json:"month"`
	AmountPaid  float64 `json:"amount_paid"`
	PaymentMode string  `json:"payment_mode"`
	FeeID        uint    `json:"fee_id"`
}

func (dto *CreatePaymentDTO) Validate() error {
	if strings.TrimSpace(dto.Month) == "" {
		return errors.New("month is required")
	}

	if dto.AmountPaid <= 0 {
		return errors.New("amount paid must be greater than zero")
	}

	if dto.FeeID == 0 {
		return errors.New("fee id is required")
	}

	return nil
}

func (dto *CreatePaymentDTO) Sanitize() {
	dto.Month = strings.TrimSpace(dto.Month)
}