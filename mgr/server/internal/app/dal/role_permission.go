package dal

import (
	"context"
	"sort"
	"sync"
	"time"

	"clearbill/mgr/server/internal/app/dal/dbmodel"
	"gorm.io/gorm"
)

type RolePermissionDAL struct {
	DB    *gorm.DB
	mu    sync.Mutex
	items map[uint]map[string]time.Time
}

func NewRolePermissionDAL(db *gorm.DB) *RolePermissionDAL {
	return &RolePermissionDAL{
		DB:    db,
		items: make(map[uint]map[string]time.Time),
	}
}

func (d *RolePermissionDAL) Replace(ctx context.Context, roleID uint, permissionIDs []string) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()

		next := make(map[string]time.Time, len(permissionIDs))
		now := time.Now()
		for _, permissionID := range permissionIDs {
			next[permissionID] = now
		}
		d.items[roleID] = next
		return nil
	}

	return d.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&dbmodel.RolePermission{}).Error; err != nil {
			return err
		}
		if len(permissionIDs) == 0 {
			return nil
		}

		items := make([]dbmodel.RolePermission, 0, len(permissionIDs))
		for _, permissionID := range permissionIDs {
			items = append(items, dbmodel.RolePermission{
				RoleID:       roleID,
				PermissionID: permissionID,
			})
		}
		return tx.Create(&items).Error
	})
}

func (d *RolePermissionDAL) ListByRoleIDs(ctx context.Context, roleIDs []uint) (map[uint][]string, error) {
	result := make(map[uint][]string, len(roleIDs))
	if len(roleIDs) == 0 {
		return result, nil
	}

	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()

		for _, roleID := range roleIDs {
			roleItems := d.items[roleID]
			if len(roleItems) == 0 {
				result[roleID] = []string{}
				continue
			}

			permissionIDs := make([]string, 0, len(roleItems))
			for permissionID := range roleItems {
				permissionIDs = append(permissionIDs, permissionID)
			}
			sort.Strings(permissionIDs)
			result[roleID] = permissionIDs
		}
		return result, nil
	}

	var items []dbmodel.RolePermission
	if err := d.DB.WithContext(ctx).
		Where("role_id IN ?", roleIDs).
		Order("permission_id asc").
		Find(&items).Error; err != nil {
		return nil, err
	}

	for _, roleID := range roleIDs {
		result[roleID] = []string{}
	}
	for _, item := range items {
		result[item.RoleID] = append(result[item.RoleID], item.PermissionID)
	}
	return result, nil
}

func (d *RolePermissionDAL) DeleteByRoleID(ctx context.Context, roleID uint) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()

		delete(d.items, roleID)
		return nil
	}

	return d.DB.WithContext(ctx).Where("role_id = ?", roleID).Delete(&dbmodel.RolePermission{}).Error
}
