package model

import (
	"time"

	"gorm.io/gorm"
)

type Faculty struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Gender      string    `gorm:"type:varchar(255)" json:"gender"`
	JoiningDate time.Time `gorm:"type:date" json:"joining_date"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
	IsActive  bool           `json:"isactive" gorm:"default:true"`

	DepartmentID uint `json:"department_id"`

	Students []Student `gorm:"foreignKey:FacultyID;references:ID" json:"students"`
}
