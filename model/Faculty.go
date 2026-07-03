package model

import (
	"time"

	"gorm.io/gorm"
)

type Faculty struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Gender      string    `json:"gender"`
	JoiningDate time.Time `gorm:"column:joining_date" json:"joining_date"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	IsActive  bool           `gorm:"column:isactive;default:true" json:"isactive"`

	DepartmentID uint `json:"department_id"`
	// Department   Department `gorm:"foreignKey:DepartmentID;references:ID" json:"-"`

	Students []Student `gorm:"foreignKey:FacultyID;references:ID" json:"students"`
}
