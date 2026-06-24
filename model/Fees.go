package model

import (
	"time"

	"gorm.io/gorm"
)

type Fees struct {
	FeesID       uint `gorm:"column:id;primaryKey;autoIncrement"`
	StudentID    uint
	Student      Student `gorm:"foreignKey:StudentID"`
	Payment_mode string
	Amount       uint
	Payment_date time.Time
	Due_date     time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}
