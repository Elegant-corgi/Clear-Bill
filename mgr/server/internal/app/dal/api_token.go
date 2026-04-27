package dal

import (
	"context"
	"sync"
	"time"

	"clearbill/mgr/server/internal/app/dal/dbmodel"
	"gorm.io/gorm"
)

type APITokenDAL struct {
	DB    *gorm.DB
	mu    sync.Mutex
	items map[string]dbmodel.UserAPIToken
}

func NewAPITokenDAL(db *gorm.DB) *APITokenDAL {
	return &APITokenDAL{
		DB:    db,
		items: make(map[string]dbmodel.UserAPIToken),
	}
}

func (d *APITokenDAL) Create(ctx context.Context, token *dbmodel.UserAPIToken) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		now := time.Now()
		token.CreatedAt = now
		token.UpdatedAt = now
		d.items[token.Token] = *token
		return nil
	}
	return d.DB.WithContext(ctx).Create(token).Error
}

func (d *APITokenDAL) GetByToken(ctx context.Context, token string) (*dbmodel.UserAPIToken, error) {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		item, ok := d.items[token]
		if !ok {
			return nil, gorm.ErrRecordNotFound
		}
		result := item
		return &result, nil
	}
	var item dbmodel.UserAPIToken
	if err := d.DB.WithContext(ctx).Where("token = ?", token).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
