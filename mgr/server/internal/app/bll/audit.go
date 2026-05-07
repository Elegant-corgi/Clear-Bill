package bll

import (
	"context"
	"errors"
	"strings"
	"time"

	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/dal"
	"clearbill/mgr/server/internal/app/dal/dbmodel"
)

type AuditService struct {
	AuditDAL *dal.AuditDAL
}

func NewAuditService(auditDAL *dal.AuditDAL) *AuditService {
	return &AuditService{AuditDAL: auditDAL}
}

func (s *AuditService) Record(ctx context.Context, user, operation, resource, result string, occurredAt time.Time) error {
	user = strings.TrimSpace(user)
	if user == "" {
		user = "anonymous"
	}
	operation = strings.TrimSpace(operation)
	if operation == "" {
		operation = "unknown"
	}
	resource = strings.TrimSpace(resource)
	if resource == "" {
		resource = "/"
	}
	result = strings.TrimSpace(result)
	if result == "" {
		result = "success"
	}
	if occurredAt.IsZero() {
		occurredAt = time.Now().UTC()
	}

	return s.AuditDAL.Create(ctx, &dbmodel.AuditLog{
		User:       user,
		Operation:  operation,
		OccurredAt: occurredAt.UTC(),
		Resource:   resource,
		Result:     result,
	})
}

func (s *AuditService) ListAuditLogs(ctx context.Context, req *vo.ListAuditLogReq) (*vo.PageResult[vo.AuditLog], error) {
	from, to, err := parseAuditTimeRange(req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}

	pageReq := req.PageReq.Normalize()
	items, query, total, err := s.AuditDAL.List(ctx, req.User, req.Operation, from, to, pageReq.Page, pageReq.PageSize)
	if err != nil {
		return nil, err
	}

	result := make([]vo.AuditLog, 0, len(items))
	for _, item := range items {
		result = append(result, vo.AuditLog{
			ID:         item.ID,
			User:       item.User,
			Operation:  item.Operation,
			OccurredAt: item.OccurredAt,
			Resource:   item.Resource,
			Result:     item.Result,
		})
	}

	return &vo.PageResult[vo.AuditLog]{
		List:     result,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

func parseAuditTimeRange(startValue, endValue string) (*time.Time, *time.Time, error) {
	startValue = strings.TrimSpace(startValue)
	endValue = strings.TrimSpace(endValue)

	start, err := parseAuditTime(startValue)
	if err != nil {
		return nil, nil, err
	}
	end, err := parseAuditTime(endValue)
	if err != nil {
		return nil, nil, err
	}

	if start != nil && end != nil && start.After(*end) {
		return nil, nil, errors.New("startTime must be before endTime")
	}

	return start, end, nil
}

func parseAuditTime(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			utc := parsed.UTC()
			return &utc, nil
		}
	}

	return nil, errors.New("invalid time format")
}
