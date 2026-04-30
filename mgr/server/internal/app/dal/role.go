package dal

import (
	"context"
	"sync"
	"time"

	"clearbill/mgr/server/internal/app/dal/dbmodel"
	"gorm.io/gorm"
)

type RoleDAL struct {
	DB     *gorm.DB
	mu     sync.Mutex
	nextID uint
	items  map[uint]dbmodel.Role
}

func NewRoleDAL(db *gorm.DB) *RoleDAL {
	dal := &RoleDAL{
		DB:     db,
		nextID: 1,
		items:  make(map[uint]dbmodel.Role),
	}
	if db == nil {
		dal.seedBuiltinMemory()
	}
	return dal
}

func (d *RoleDAL) seedBuiltinMemory() {
	now := time.Now()
	for _, item := range builtinRoles() {
		item.ID = d.nextID
		item.CreatedAt = now
		item.UpdatedAt = now
		d.items[item.ID] = item
		d.nextID++
	}
}

func (d *RoleDAL) Create(ctx context.Context, role *dbmodel.Role) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()

		for _, item := range d.items {
			if item.Code == role.Code {
				return gorm.ErrDuplicatedKey
			}
		}
		now := time.Now()
		role.ID = d.nextID
		role.CreatedAt = now
		role.UpdatedAt = now
		d.items[role.ID] = *role
		d.nextID++
		return nil
	}

	return d.DB.WithContext(ctx).Create(role).Error
}

func (d *RoleDAL) List(ctx context.Context, keyword, scope string, tenantID *uint, page, pageSize int) ([]dbmodel.Role, PageQuery, int64, error) {
	query := NewPageQuery(page, pageSize)

	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()

		filtered := make([]dbmodel.Role, 0, len(d.items))
		for _, role := range d.items {
			if scope != "" && role.Scope != scope {
				continue
			}
			if tenantID != nil {
				if role.Scope != dbmodel.RoleScopeTenant || role.TenantID == nil || *role.TenantID != *tenantID {
					continue
				}
			}
			if keyword == "" || containsInsensitive(role.Code, keyword) || containsInsensitive(role.Name, keyword) {
				filtered = append(filtered, role)
			}
		}

		total := int64(len(filtered))
		start := query.Offset()
		if start >= len(filtered) {
			return []dbmodel.Role{}, query, total, nil
		}

		end := start + query.PageSize
		if end > len(filtered) {
			end = len(filtered)
		}

		return filtered[start:end], query, total, nil
	}

	var roles []dbmodel.Role
	tx := d.DB.WithContext(ctx).Order("builtin desc, id asc")
	if keyword != "" {
		like := "%" + keyword + "%"
		tx = tx.Where("code LIKE ? OR name LIKE ?", like, like)
	}
	if scope != "" {
		tx = tx.Where("scope = ?", scope)
	}
	if tenantID != nil {
		tx = tx.Where("scope = ? AND tenant_id = ?", dbmodel.RoleScopeTenant, *tenantID)
	}
	query, total, err := FindPage(tx, &dbmodel.Role{}, page, pageSize, &roles)
	if err != nil {
		return nil, query, 0, err
	}
	return roles, query, total, nil
}

func (d *RoleDAL) GetByID(ctx context.Context, id uint) (*dbmodel.Role, error) {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()

		item, ok := d.items[id]
		if !ok {
			return nil, gorm.ErrRecordNotFound
		}
		result := item
		return &result, nil
	}

	var role dbmodel.Role
	if err := d.DB.WithContext(ctx).First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (d *RoleDAL) GetByCode(ctx context.Context, code string) (*dbmodel.Role, error) {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()

		for _, role := range d.items {
			if role.Code == code {
				result := role
				return &result, nil
			}
		}
		return nil, gorm.ErrRecordNotFound
	}

	var role dbmodel.Role
	if err := d.DB.WithContext(ctx).Where("code = ?", code).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (d *RoleDAL) Update(ctx context.Context, role *dbmodel.Role) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()

		if _, ok := d.items[role.ID]; !ok {
			return gorm.ErrRecordNotFound
		}
		role.UpdatedAt = time.Now()
		d.items[role.ID] = *role
		return nil
	}

	return d.DB.WithContext(ctx).Save(role).Error
}

func (d *RoleDAL) Delete(ctx context.Context, id uint) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()

		if _, ok := d.items[id]; !ok {
			return gorm.ErrRecordNotFound
		}
		delete(d.items, id)
		return nil
	}

	return d.DB.WithContext(ctx).Delete(&dbmodel.Role{}, id).Error
}

func builtinRoles() []dbmodel.Role {
	return []dbmodel.Role{
		{
			Code:    dbmodel.RoleSysadmin,
			Name:    "System Administrator",
			Scope:   dbmodel.RoleScopeSystem,
			Builtin: true,
		},
		{
			Code:    dbmodel.RoleTenantAdmin,
			Name:    "Tenant Administrator",
			Scope:   dbmodel.RoleScopeTenant,
			Builtin: true,
		},
		{
			Code:    dbmodel.RoleUser,
			Name:    "User",
			Scope:   dbmodel.RoleScopeTenant,
			Builtin: true,
		},
	}
}
