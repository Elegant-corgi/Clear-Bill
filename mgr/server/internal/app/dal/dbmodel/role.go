package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        uint   `gorm:"primaryKey"`
	Code      string `gorm:"size:64;not null;uniqueIndex"`
	Name      string `gorm:"size:128;not null"`
	Scope     string `gorm:"size:16;not null;index"`
	TenantID  *uint  `gorm:"index"`
	Builtin   bool   `gorm:"not null;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Role) TableName() string {
	return "roles"
}
