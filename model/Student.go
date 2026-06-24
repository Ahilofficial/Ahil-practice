package model

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	StudentID uint `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string
	Email     string
	Gender    string
	FacultyID uint
	Faculty   Faculty `gorm:"foreignKey:FacultyID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Fees []Fees `gorm:"foreignKey:StudentID"`
}
