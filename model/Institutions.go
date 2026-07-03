package model

import "time"

type Institutions struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string     `json:"name"`
	InstitutionCode string     `gorm:"column:institution_code;" json:"institution_code"`
	State           string     `json:"state"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"-"`
	IsActive        bool       `gorm:"column:isactive;default:true" json:"isactive"`

	Departments []Department `gorm:"foreignKey:InstitutionID;references:ID" json:"departments"`
}
