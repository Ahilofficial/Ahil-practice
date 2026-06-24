package model

import (
	"time"

	"gorm.io/gorm"
)

type Department struct {
	DepartmentID    uint `gorm:"column:id;primaryKey;autoIncrement"`
	Department_Name string
	InstitutionID   uint
	Institution     Institutions `gorm:"foreignKey:InstitutionID"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt

	Faculty []Faculty `gorm:"foreignKey:DepartmentID"`
}
