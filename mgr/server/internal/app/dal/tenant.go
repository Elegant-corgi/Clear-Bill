package dal

import (
	"context"
	"sync"
	"time"

	"clearbill/mgr/server/internal/app/dal/dbmodel"
	"gorm.io/gorm"
)

type TenantDAL struct {
	DB     *gorm.DB
	mu     sync.Mutex
	nextID uint
	items  map[uint]dbmodel.Tenant
}

func NewTenantDAL(db *gorm.DB) *TenantDAL {
	return &TenantDAL{
		DB:     db,
		nextID: 1,
		items:  make(map[uint]dbmodel.Tenant),
	}
}

func (d *TenantDAL) Create(ctx context.Context, tenant *dbmodel.Tenant) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()

		now := time.Now()
		tenant.ID = d.nextID
		tenant.CreatedAt = now
		tenant.UpdatedAt = now
		d.items[tenant.ID] = *tenant
		d.nextID++
		return nil
	}

	return d.DB.WithContext(ctx).Create(tenant).Error
}

func (d *TenantDAL) List(ctx context.Context, keyword string) ([]dbmodel.Tenant, error) {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()

		tenants := make([]dbmodel.Tenant, 0, len(d.items))
		for _, tenant := range d.items {
			if keyword == "" || containsInsensitive(tenant.Code, keyword) || containsInsensitive(tenant.Name, keyword) {
				tenants = append(tenants, tenant)
			}
		}

		return tenants, nil
	}

	var tenants []dbmodel.Tenant

	tx := d.DB.WithContext(ctx).Order("id desc")
	if keyword != "" {
		like := "%" + keyword + "%"
		tx = tx.Where("code LIKE ? OR name LIKE ?", like, like)
	}

	if err := tx.Find(&tenants).Error; err != nil {
		return nil, err
	}

	return tenants, nil
}

func (d *TenantDAL) GetByID(ctx context.Context, id uint) (*dbmodel.Tenant, error) {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()

		tenant, ok := d.items[id]
		if !ok {
			return nil, gorm.ErrRecordNotFound
		}

		result := tenant
		return &result, nil
	}

	var tenant dbmodel.Tenant
	if err := d.DB.WithContext(ctx).First(&tenant, id).Error; err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (d *TenantDAL) Update(ctx context.Context, tenant *dbmodel.Tenant) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()

		if _, ok := d.items[tenant.ID]; !ok {
			return gorm.ErrRecordNotFound
		}
		tenant.UpdatedAt = time.Now()
		d.items[tenant.ID] = *tenant
		return nil
	}

	return d.DB.WithContext(ctx).Save(tenant).Error
}

func (d *TenantDAL) Delete(ctx context.Context, id uint) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()

		if _, ok := d.items[id]; !ok {
			return gorm.ErrRecordNotFound
		}
		delete(d.items, id)
		return nil
	}

	return d.DB.WithContext(ctx).Delete(&dbmodel.Tenant{}, id).Error
}

func containsInsensitive(value, keyword string) bool {
	if keyword == "" {
		return true
	}

	valueRunes := []rune(value)
	keywordRunes := []rune(keyword)
	if len(keywordRunes) > len(valueRunes) {
		return false
	}

	for i := 0; i+len(keywordRunes) <= len(valueRunes); i++ {
		match := true
		for j := range keywordRunes {
			if toLowerASCII(valueRunes[i+j]) != toLowerASCII(keywordRunes[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}

	return false
}

func toLowerASCII(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + ('a' - 'A')
	}

	return r
}
