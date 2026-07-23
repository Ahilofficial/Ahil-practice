package model

type Session struct {
    ID        uint   `gorm:"primaryKey"`
    UserID    uint
	SessionID string
	AccessToken string
	RefreshToken string
    IsActive  bool
    Platform  string
}