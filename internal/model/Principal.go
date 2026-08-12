package model

import (
	"time"

	"gorm.io/gorm"
)

type Principal struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Gender      string         `gorm:"type:varchar(255)" json:"gender"`
	JoiningDate time.Time      `gorm:"type:date" json:"joining_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-"`
	IsActive    bool           `gorm:"default:true" json:"isactive"`

	DepartmentID uint `json:"department_id"`
	UserID       uint `json:"user_id"`

	Department *Department `gorm:"foreignKey:DepartmentID;references:ID" json:"department,omitempty"`
	User       *User       `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}
