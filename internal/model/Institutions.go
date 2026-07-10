package model

import "time"

type Institutions struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string     `gorm:"type:varchar(255)" json:"name"`
	InstitutionCode string     `gorm:"type:varchar(255)" json:"institution_code"`
	State           string     `gorm:"type:varchar(255)" json:"state"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"-"`
	IsActive        bool       `json:"isactive" gorm:"default:true"`

	Departments []Department `gorm:"foreignKey:InstitutionID;references:ID" json:"departments"`
}
