package dbmodel

import "time"

type RolePermission struct {
	ID           uint   `gorm:"primaryKey"`
	RoleID       uint   `gorm:"not null;uniqueIndex:uk_role_permission;index"`
	PermissionID string `gorm:"size:128;not null;uniqueIndex:uk_role_permission;index"`
	CreatedAt    time.Time
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
