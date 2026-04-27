package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

type UserAPIToken struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;index"`
	Name      string `gorm:"size:128;not null"`
	Token     string `gorm:"size:128;not null;uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserAPIToken) TableName() string {
	return "user_api_tokens"
}
