package model

import (
	"time"

	"gorm.io/gorm"
)

type Department struct {
	ID             uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	DepartmentName string `gorm:"type:varchar(255)" json:"department_name"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
	IsActive  bool           `json:"isactive" gorm:"default:true"`

	InstitutionID uint `json:"institution_id"`

	Faculties []Faculty `gorm:"foreignKey:DepartmentID;references:ID" json:"-"`
}
