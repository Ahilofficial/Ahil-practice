package model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"unique;not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions,omitempty"`
}
