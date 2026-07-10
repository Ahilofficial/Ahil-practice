package model

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name   string `gorm:"type:varchar(255);not null" json:"name"`
	Email  string `gorm:"type:varchar(255);unique;not null" json:"email"`
	Gender string `gorm:"type:varchar(255)" json:"gender"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
	IsActive  bool           `json:"isactive" gorm:"default:true"`

	FacultyID uint `json:"faculty_id"`

	Fees []Fees `gorm:"foreignKey:StudentID;references:ID" json:"fees"`
}
