package model

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Month       string  `json:"month"`
	AmountPaid  float64 `gorm:"type:decimal(10,2)" json:"amount_paid"`
	PaymentMode string  `json:"payment_mode"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`

	FeeID uint `json:"fee_id"`
	Fee   Fees `gorm:"foreignKey:FeeID;references:ID" json:"-"`
}