package model

import (
	"time"

	"gorm.io/gorm"
)

type Fees struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	PaymentMode   string         `gorm:"type:varchar(255)" json:"payment_mode"`
	TotalAmount   float64        `gorm:"type:decimal(10,2)" json:"amount"`
	TotalPaid     float64        `gorm:"type:decimal(10,2)" json:"total_paid"`
	PendingAmount float64        `gorm:"type:decimal(10,2)" json:"pending_amount"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-"`
	IsActive      bool           `gorm:"default:true" json:"isactive"`

	StudentID uint `json:"student_id"`

	Student  *Student  `gorm:"foreignKey:StudentID;references:ID" json:"student,omitempty"`
	Payments []Payment `gorm:"foreignKey:FeeID;references:ID" json:"payments,omitempty"`
}
