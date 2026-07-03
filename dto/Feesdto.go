package dto

type CreateFeesDTO struct {
	PaymentMode string `json:"payment_mode"`
	Amount      uint   `json:"amount"`
	StudentID   uint   `json:"student_id"`
}
