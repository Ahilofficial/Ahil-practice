package model

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name   string `gorm:"not null" json:"name"`
	Email  string `gorm:"unique;not null" json:"email"`
	Gender string `json:"gender"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	IsActive  bool           `gorm:"column:isactive;default:true" json:"isactive"`

	FacultyID uint `json:"faculty_id"`
	// Faculty   Faculty `gorm:"foreignKey:FacultyID;references:ID" json:"-"`

	Fees []Fees `gorm:"foreignKey:StudentID;references:ID" json:"fees"`
}
