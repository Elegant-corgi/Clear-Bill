package dal

import (
	"context"
	"strings"
	"sync"
	"time"

	"clearbill/mgr/server/internal/app/dal/dbmodel"
	"gorm.io/gorm"
)

type CredentialDAL struct {
	DB     *gorm.DB
	mu     sync.Mutex
	nextID uint
	items  map[uint]dbmodel.UserCredential
}

func NewCredentialDAL(db *gorm.DB) *CredentialDAL {
	return &CredentialDAL{
		DB:     db,
		nextID: 1,
		items:  make(map[uint]dbmodel.UserCredential),
	}
}

func (d *CredentialDAL) Create(ctx context.Context, credential *dbmodel.UserCredential) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		if d.duplicateLocked(credential) {
			return gorm.ErrDuplicatedKey
		}
		now := time.Now()
		credential.ID = d.nextID
		credential.CreatedAt = now
		credential.UpdatedAt = now
		d.items[credential.ID] = *credential
		d.nextID++
		return nil
	}
	return d.DB.WithContext(ctx).Create(credential).Error
}

func (d *CredentialDAL) Update(ctx context.Context, credential *dbmodel.UserCredential) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		if _, ok := d.items[credential.ID]; !ok {
			return gorm.ErrRecordNotFound
		}
		credential.UpdatedAt = time.Now()
		d.items[credential.ID] = *credential
		return nil
	}
	return d.DB.WithContext(ctx).Save(credential).Error
}

func (d *CredentialDAL) Delete(ctx context.Context, id uint) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		if _, ok := d.items[id]; !ok {
			return gorm.ErrRecordNotFound
		}
		delete(d.items, id)
		return nil
	}
	return d.DB.WithContext(ctx).Delete(&dbmodel.UserCredential{}, id).Error
}

func (d *CredentialDAL) TouchLastUsed(ctx context.Context, id uint) error {
	now := time.Now()
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		item, ok := d.items[id]
		if !ok {
			return gorm.ErrRecordNotFound
		}
		item.LastUsedAt = &now
		item.UpdatedAt = now
		d.items[id] = item
		return nil
	}
	return d.DB.WithContext(ctx).Model(&dbmodel.UserCredential{}).Where("id = ?", id).Updates(map[string]any{
		"last_used_at": now,
		"updated_at":   now,
	}).Error
}

func (d *CredentialDAL) GetByID(ctx context.Context, id uint) (*dbmodel.UserCredential, error) {
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
	var credential dbmodel.UserCredential
	if err := d.DB.WithContext(ctx).First(&credential, id).Error; err != nil {
		return nil, err
	}
	return &credential, nil
}

func (d *CredentialDAL) GetActiveByToken(ctx context.Context, token string) (*dbmodel.UserCredential, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, gorm.ErrRecordNotFound
	}

	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		for _, item := range d.items {
			if item.Type == dbmodel.CredentialTypeToken && item.Status == dbmodel.CredentialStatusActive && item.Token != nil && *item.Token == token {
				result := item
				return &result, nil
			}
		}
		return nil, gorm.ErrRecordNotFound
	}

	var credential dbmodel.UserCredential
	if err := d.DB.WithContext(ctx).
		Where("type = ? AND status = ? AND token = ?", dbmodel.CredentialTypeToken, dbmodel.CredentialStatusActive, token).
		First(&credential).Error; err != nil {
		return nil, err
	}
	return &credential, nil
}

func (d *CredentialDAL) GetActiveByAccessKey(ctx context.Context, accessKey string) (*dbmodel.UserCredential, error) {
	accessKey = strings.TrimSpace(accessKey)
	if accessKey == "" {
		return nil, gorm.ErrRecordNotFound
	}

	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		for _, item := range d.items {
			if item.Type == dbmodel.CredentialTypeAKSK && item.Status == dbmodel.CredentialStatusActive && item.AccessKey != nil && *item.AccessKey == accessKey {
				result := item
				return &result, nil
			}
		}
		return nil, gorm.ErrRecordNotFound
	}

	var credential dbmodel.UserCredential
	if err := d.DB.WithContext(ctx).
		Where("type = ? AND status = ? AND access_key = ?", dbmodel.CredentialTypeAKSK, dbmodel.CredentialStatusActive, accessKey).
		First(&credential).Error; err != nil {
		return nil, err
	}
	return &credential, nil
}

func (d *CredentialDAL) List(ctx context.Context, userID *uint, typ, status string, page, pageSize int) ([]dbmodel.UserCredential, PageQuery, int64, error) {
	query := NewPageQuery(page, pageSize)

	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()

		filtered := make([]dbmodel.UserCredential, 0, len(d.items))
		for _, item := range d.items {
			if userID != nil && item.UserID != *userID {
				continue
			}
			if typ != "" && item.Type != typ {
				continue
			}
			if status != "" && item.Status != status {
				continue
			}
			filtered = append(filtered, item)
		}

		total := int64(len(filtered))
		start := query.Offset()
		if start >= len(filtered) {
			return []dbmodel.UserCredential{}, query, total, nil
		}

		end := start + query.PageSize
		if end > len(filtered) {
			end = len(filtered)
		}

		return filtered[start:end], query, total, nil
	}

	var items []dbmodel.UserCredential
	tx := d.DB.WithContext(ctx).Order("id desc")
	if userID != nil {
		tx = tx.Where("user_id = ?", *userID)
	}
	if typ != "" {
		tx = tx.Where("type = ?", typ)
	}
	if status != "" {
		tx = tx.Where("status = ?", status)
	}
	query, total, err := FindPage(tx, &dbmodel.UserCredential{}, page, pageSize, &items)
	if err != nil {
		return nil, query, 0, err
	}
	return items, query, total, nil
}

func (d *CredentialDAL) duplicateLocked(candidate *dbmodel.UserCredential) bool {
	for _, item := range d.items {
		if sameOptionalString(item.AccessKey, candidate.AccessKey) {
			return true
		}
		if sameOptionalString(item.Token, candidate.Token) {
			return true
		}
	}
	return false
}

func sameOptionalString(a, b *string) bool {
	if a == nil || b == nil {
		return false
	}
	return strings.TrimSpace(*a) != "" && *a == *b
}
