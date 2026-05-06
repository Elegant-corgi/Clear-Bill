package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

const (
	CredentialTypeAKSK  = "aksk"
	CredentialTypeToken = "token"

	CredentialStatusActive  = "active"
	CredentialStatusRotated = "rotated"
	CredentialStatusRevoked  = "revoked"
)

type UserCredential struct {
	ID         uint           `gorm:"primaryKey"`
	UserID     uint           `gorm:"not null;index"`
	Type       string         `gorm:"size:16;not null;index"`
	Name       string         `gorm:"size:128;not null"`
	AccessKey  *string        `gorm:"size:128;uniqueIndex"`
	SecretKey  *string        `gorm:"size:128"`
	Token      *string        `gorm:"size:128;uniqueIndex"`
	Status     string         `gorm:"size:16;not null;index"`
	ParentID   *uint          `gorm:"index"`
	ExpiresAt  *time.Time     `gorm:"index"`
	RotatedAt  *time.Time     `gorm:"index"`
	LastUsedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (UserCredential) TableName() string {
	return "user_credentials"
}
