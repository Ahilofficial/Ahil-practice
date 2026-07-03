package model

import (
	"time"

	"gorm.io/gorm"
)

type Department struct {
	ID             uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	DepartmentName string `gorm:"column:department_name;not null" json:"department_name"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	IsActive  bool           `gorm:"column:isactive;default:true" json:"isactive"`

	InstitutionID uint `json:"institution_id"`
	// Institution   Institutions `gorm:"foreignKey:InstitutionID;references:ID" json:"-"`

	Faculties []Faculty `gorm:"foreignKey:DepartmentID;references:ID" json:"faculties"`
}
