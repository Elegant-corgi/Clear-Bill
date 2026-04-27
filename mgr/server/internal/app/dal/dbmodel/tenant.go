package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	ID           uint   `gorm:"primaryKey"`
	Code         string `gorm:"size:64;not null;uniqueIndex"`
	Name         string `gorm:"size:128;not null"`
	ContactName  string `gorm:"size:64"`
	ContactPhone string `gorm:"size:32"`
	Status       string `gorm:"size:32;not null;default:active"`
	Remark       string `gorm:"size:255"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (Tenant) TableName() string {
	return "tenants"
}
