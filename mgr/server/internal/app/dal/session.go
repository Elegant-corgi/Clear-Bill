package dal

import (
	"context"
	"sync"
	"time"

	"clearbill/mgr/server/internal/app/dal/dbmodel"
	"gorm.io/gorm"
)

type SessionDAL struct {
	DB    *gorm.DB
	mu    sync.Mutex
	items map[string]dbmodel.UserSession
}

func NewSessionDAL(db *gorm.DB) *SessionDAL {
	return &SessionDAL{
		DB:    db,
		items: make(map[string]dbmodel.UserSession),
	}
}

func (d *SessionDAL) Create(ctx context.Context, session *dbmodel.UserSession) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		now := time.Now()
		session.CreatedAt = now
		session.UpdatedAt = now
		d.items[session.Token] = *session
		return nil
	}
	return d.DB.WithContext(ctx).Create(session).Error
}

func (d *SessionDAL) GetByToken(ctx context.Context, token string) (*dbmodel.UserSession, error) {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		session, ok := d.items[token]
		if !ok {
			return nil, gorm.ErrRecordNotFound
		}
		result := session
		return &result, nil
	}
	var session dbmodel.UserSession
	if err := d.DB.WithContext(ctx).Where("token = ?", token).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (d *SessionDAL) DeleteByToken(ctx context.Context, token string) error {
	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		delete(d.items, token)
		return nil
	}
	return d.DB.WithContext(ctx).Where("token = ?", token).Delete(&dbmodel.UserSession{}).Error
}
