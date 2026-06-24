package model

import (
	"time"

	"gorm.io/gorm"
)

type Institutions struct {
	InstitutionID    uint `gorm:"column:id;primaryKey;autoIncrement"`
	Name             string
	Institution_code string
	State            string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt

	Departments []Department `gorm:"foreignKey:InstitutionID"`
}
