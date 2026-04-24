package bll

import (
	"context"

	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/dal"
)

type BillingService struct {
	BillingDAL *dal.BillingDAL
}

func NewBillingService(billingDAL *dal.BillingDAL) *BillingService {
	return &BillingService{
		BillingDAL: billingDAL,
	}
}

func (s *BillingService) GetDashboardSummary(ctx context.Context) (*vo.DashboardSummary, error) {
	snapshot, err := s.BillingDAL.GetDashboardSummary(ctx)
	if err != nil {
		return nil, err
	}

	return &vo.DashboardSummary{
		OpenBills:       snapshot.OpenBills,
		PendingInvoices: snapshot.PendingInvoices,
		OverdueAmount:   snapshot.OverdueAmount,
		AutoMatchedRate: snapshot.AutoMatchedRate,
		RecentActivity:  snapshot.RecentActivity,
		AttentionList:   snapshot.AttentionList,
	}, nil
}

func (s *BillingService) ListBills(ctx context.Context) ([]vo.Bill, error) {
	return s.BillingDAL.ListBills(ctx)
}

func (s *BillingService) ListCustomers(ctx context.Context) ([]vo.Customer, error) {
	return s.BillingDAL.ListCustomers(ctx)
}

func (s *BillingService) ListReconciliationTasks(ctx context.Context) ([]vo.ReconciliationTask, error) {
	return s.BillingDAL.ListReconciliationTasks(ctx)
}
