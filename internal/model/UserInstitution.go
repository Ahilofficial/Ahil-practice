package model

type Institution_Admins struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	UserID        uint `gorm:"not null"`
	InstitutionID uint `gorm:"not null"`
}
