package model

import (
	"time"

	"gorm.io/gorm"
)

type Faculty struct {
	FacultyID    uint `gorm:"column:id;primaryKey;autoIncrement"`
	Name         string
	Gender       string
	Joining_Date time.Time
	DepartmentID uint
	Department   Department `gorm:"foreignKey:DepartmentID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt

	Students []Student `gorm:"foreignKey:FacultyID"`
}
