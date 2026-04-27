package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

type UserSession struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"size:128;not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}
