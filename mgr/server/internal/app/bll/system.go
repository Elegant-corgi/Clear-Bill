package bll

import (
	"context"
	"time"

	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/dal"
	"clearbill/mgr/server/internal/version"
)

type SystemService struct {
	SystemDAL *dal.SystemDAL
}

func NewSystemService(systemDAL *dal.SystemDAL) *SystemService {
	return &SystemService{
		SystemDAL: systemDAL,
	}
}

func (s *SystemService) Health(ctx context.Context) (*vo.HealthStatus, error) {
	seed, err := s.SystemDAL.GetHealth(ctx)
	if err != nil {
		return nil, err
	}

	return &vo.HealthStatus{
		Name:      seed.Name,
		Status:    seed.Status,
		Version:   version.Version,
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}
