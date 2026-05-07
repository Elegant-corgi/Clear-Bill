package dal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"clearbill/mgr/server/internal/app/dal/dbmodel"
	"gorm.io/gorm"
)

type AuditDAL struct {
	DB         *gorm.DB
	mu         sync.Mutex
	nextID     uint
	items      map[uint]dbmodel.AuditLog
	lastHash   string
	lastSeenID uint
}

func NewAuditDAL(db *gorm.DB) *AuditDAL {
	return &AuditDAL{
		DB:     db,
		nextID: 1,
		items:  make(map[uint]dbmodel.AuditLog),
	}
}

func (d *AuditDAL) Create(ctx context.Context, log *dbmodel.AuditLog) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if log == nil {
		return errors.New("audit log is required")
	}
	if log.OccurredAt.IsZero() {
		log.OccurredAt = time.Now().UTC()
	}

	if d.DB == nil {
		log.ID = d.nextID
		log.PrevHash = d.lastHash
		log.EntryHash = buildAuditEntryHash(log.PrevHash, log.User, log.Operation, log.OccurredAt, log.Resource, log.Result)
		log.CreatedAt = log.OccurredAt
		d.items[log.ID] = *log
		d.lastHash = log.EntryHash
		d.lastSeenID = log.ID
		d.nextID++
		return nil
	}

	prevHash, err := d.latestHash(ctx)
	if err != nil {
		return err
	}
	log.PrevHash = prevHash
	log.EntryHash = buildAuditEntryHash(log.PrevHash, log.User, log.Operation, log.OccurredAt, log.Resource, log.Result)
	if err := d.DB.WithContext(ctx).Create(log).Error; err != nil {
		return err
	}
	d.lastHash = log.EntryHash
	d.lastSeenID = log.ID
	return nil
}

func (d *AuditDAL) List(ctx context.Context, user, operation string, from, to *time.Time, page, pageSize int) ([]dbmodel.AuditLog, PageQuery, int64, error) {
	query := NewPageQuery(page, pageSize)

	if d.DB == nil {
		d.mu.Lock()
		defer d.mu.Unlock()

		filtered := make([]dbmodel.AuditLog, 0, len(d.items))
		for _, item := range d.items {
			if user != "" && !strings.EqualFold(item.User, user) {
				continue
			}
			if operation != "" && !strings.EqualFold(item.Operation, operation) {
				continue
			}
			if from != nil && item.OccurredAt.Before(from.UTC()) {
				continue
			}
			if to != nil && item.OccurredAt.After(to.UTC()) {
				continue
			}
			filtered = append(filtered, item)
		}

		sort.Slice(filtered, func(i, j int) bool {
			if !filtered[i].OccurredAt.Equal(filtered[j].OccurredAt) {
				return filtered[i].OccurredAt.After(filtered[j].OccurredAt)
			}
			return filtered[i].ID > filtered[j].ID
		})

		total := int64(len(filtered))
		start := query.Offset()
		if start >= len(filtered) {
			return []dbmodel.AuditLog{}, query, total, nil
		}

		end := start + query.PageSize
		if end > len(filtered) {
			end = len(filtered)
		}

		return filtered[start:end], query, total, nil
	}

	tx := d.DB.WithContext(ctx).Order("id desc")
	if user != "" {
		tx = tx.Where("user = ?", user)
	}
	if operation != "" {
		tx = tx.Where("operation = ?", operation)
	}
	if from != nil {
		tx = tx.Where("occurred_at >= ?", from.UTC())
	}
	if to != nil {
		tx = tx.Where("occurred_at <= ?", to.UTC())
	}

	var items []dbmodel.AuditLog
	query, total, err := FindPage(tx, &dbmodel.AuditLog{}, page, pageSize, &items)
	if err != nil {
		return nil, query, 0, err
	}
	return items, query, total, nil
}

func (d *AuditDAL) latestHash(ctx context.Context) (string, error) {
	if d.DB == nil {
		return d.lastHash, nil
	}

	var item dbmodel.AuditLog
	err := d.DB.WithContext(ctx).Order("id desc").Limit(1).Take(&item).Error
	switch {
	case err == nil:
		return item.EntryHash, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return "", nil
	default:
		return "", err
	}
}

func buildAuditEntryHash(prevHash, user, operation string, occurredAt time.Time, resource, result string) string {
	payload := fmt.Sprintf("%s|%s|%s|%s|%s|%s",
		prevHash,
		user,
		operation,
		occurredAt.UTC().Format(time.RFC3339Nano),
		resource,
		result,
	)
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}
