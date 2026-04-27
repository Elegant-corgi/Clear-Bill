package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

const (
	DefaultUserPassword = "bill123;"
	RoleSysadmin        = "sysadmin"
	RoleTenantAdmin     = "tenant_admin"
	RoleUser            = "user"
	RoleScopeSystem     = "system"
	RoleScopeTenant     = "tenant"
	StatusActive        = "active"
	StatusDisabled      = "disabled"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"size:64;not null;uniqueIndex"`
	DisplayName  string `gorm:"size:128;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	Role         string `gorm:"size:32;not null;index"`
	TenantID     *uint  `gorm:"index"`
	Status       string `gorm:"size:32;not null;default:active"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
