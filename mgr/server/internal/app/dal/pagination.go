package dal

import (
	"clearbill/mgr/server/api/vo"
	"gorm.io/gorm"
)

type PageQuery struct {
	Page     int
	PageSize int
}

func NewPageQuery(page, pageSize int) PageQuery {
	normalized := vo.PageReq{
		Page:     page,
		PageSize: pageSize,
	}.Normalize()

	return PageQuery{
		Page:     normalized.Page,
		PageSize: normalized.PageSize,
	}
}

func (q PageQuery) Offset() int {
	return (q.Page - 1) * q.PageSize
}

func FindPage[T any](tx *gorm.DB, model any, page, pageSize int, dest *[]T) (PageQuery, int64, error) {
	query := NewPageQuery(page, pageSize)

	var total int64
	if err := tx.Session(&gorm.Session{}).Model(model).Count(&total).Error; err != nil {
		return query, 0, err
	}

	if err := tx.Offset(query.Offset()).Limit(query.PageSize).Find(dest).Error; err != nil {
		return query, 0, err
	}

	return query, total, nil
}
