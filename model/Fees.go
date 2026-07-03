package model

import (
	"time"

	"gorm.io/gorm"
)

type Fees struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	PaymentMode string         `gorm:"column:payment_mode" json:"payment_mode"`
	Amount      uint           `json:"amount"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	IsActive    bool           `gorm:"column:isactive;default:true" json:"isactive"`
	StudentID   uint           `json:"student_id"`

	Student Student `gorm:"foreignKey:StudentID;references:ID" json:"-"`
}
