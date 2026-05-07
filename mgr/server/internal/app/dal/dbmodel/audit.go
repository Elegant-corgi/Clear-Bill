package dbmodel

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type AuditLog struct {
	ID         uint      `gorm:"primaryKey"`
	User       string    `gorm:"size:128;not null;index:idx_audit_user_operation_time,priority:1;<-:create"`
	Operation  string    `gorm:"size:128;not null;index:idx_audit_user_operation_time,priority:2;<-:create"`
	OccurredAt time.Time `gorm:"not null;index:idx_audit_user_operation_time,priority:3;<-:create"`
	Resource   string    `gorm:"size:255;not null;<-:create"`
	Result     string    `gorm:"size:64;not null;<-:create"`
	PrevHash   string    `gorm:"size:64;not null;index;<-:create"`
	EntryHash  string    `gorm:"size:64;not null;uniqueIndex;<-:create"`
	CreatedAt  time.Time
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

func (AuditLog) BeforeUpdate(*gorm.DB) error {
	return errors.New("audit logs are immutable")
}

func (AuditLog) BeforeDelete(*gorm.DB) error {
	return errors.New("audit logs are immutable")
}
