package dal

import (
	"context"
	"sync"
	"time"

	"clearbill/mgr/server/internal/app/dal/dbmodel"
	"clearbill/mgr/server/pkg/passwordx"
	"gorm.io/gorm"
)

type UserDAL struct {
	DB     *gorm.DB
	mu     sync.Mutex
	nextID uint
	items  map[uint]dbmodel.User
}

func NewUserDAL(db *gorm.DB) *UserDAL {
	dal := &UserDAL{
		DB:     db,
		nextID: 1,
		items:  make(map[uint]dbmodel.User),
	}
	if db == nil {
		_ = dal.seedSysadminMemory()
	}
	return dal
}

func (d *UserDAL) seedSysadminMemory() error {
	hash, err := passwordx.HashPassword(dbmodel.DefaultUserPassword)
	if err != nil {
		return err
	}
	now := time.Now()
	d.items[d.nextID] = dbmodel.User{
		ID:           d.nextID,
		Username:     "sysadmin",
		DisplayName:  "System Administrator",
		PasswordHash: hash,
		Role:         dbmodel.RoleSysadmin,
		Status:       dbmodel.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	d.nextID++
	return nil
}

func (d *UserDAL) Create(ctx context.Context, user *dbmodel.User) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		for _, item := range d.items {
			if item.Username == user.Username {
				return gorm.ErrDuplicatedKey
			}
		}
		now := time.Now()
		user.ID = d.nextID
		user.CreatedAt = now
		user.UpdatedAt = now
		d.items[user.ID] = *user
		d.nextID++
		return nil
	}
	return d.DB.WithContext(ctx).Create(user).Error
}

func (d *UserDAL) List(ctx context.Context, keyword string, tenantID *uint) ([]dbmodel.User, error) {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		result := make([]dbmodel.User, 0, len(d.items))
		for _, user := range d.items {
			if tenantID != nil {
				if user.TenantID == nil || *user.TenantID != *tenantID {
					continue
				}
			}
			if keyword == "" || containsInsensitive(user.Username, keyword) || containsInsensitive(user.DisplayName, keyword) {
				result = append(result, user)
			}
		}
		return result, nil
	}

	var users []dbmodel.User
	tx := d.DB.WithContext(ctx).Order("id desc")
	if tenantID != nil {
		tx = tx.Where("tenant_id = ?", *tenantID)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		tx = tx.Where("username LIKE ? OR display_name LIKE ?", like, like)
	}
	if err := tx.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (d *UserDAL) GetByID(ctx context.Context, id uint) (*dbmodel.User, error) {
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
	var user dbmodel.User
	if err := d.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDAL) GetByUsername(ctx context.Context, username string) (*dbmodel.User, error) {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		for _, user := range d.items {
			if user.Username == username {
				result := user
				return &result, nil
			}
		}
		return nil, gorm.ErrRecordNotFound
	}
	var user dbmodel.User
	if err := d.DB.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDAL) Update(ctx context.Context, user *dbmodel.User) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		if _, ok := d.items[user.ID]; !ok {
			return gorm.ErrRecordNotFound
		}
		user.UpdatedAt = time.Now()
		d.items[user.ID] = *user
		return nil
	}
	return d.DB.WithContext(ctx).Save(user).Error
}

func (d *UserDAL) Delete(ctx context.Context, id uint) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		if _, ok := d.items[id]; !ok {
			return gorm.ErrRecordNotFound
		}
		delete(d.items, id)
		return nil
	}
	return d.DB.WithContext(ctx).Delete(&dbmodel.User{}, id).Error
}
