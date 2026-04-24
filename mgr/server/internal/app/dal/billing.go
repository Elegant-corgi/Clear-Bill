package dal

import (
	"context"

	"clearbill/mgr/server/api/vo"
)

type DashboardSnapshot struct {
	OpenBills       int
	PendingInvoices int
	OverdueAmount   float64
	AutoMatchedRate float64
	RecentActivity  []vo.ActivityItem
	AttentionList   []vo.AttentionItem
}

type BillingDAL struct{}

func NewBillingDAL() *BillingDAL {
	return &BillingDAL{}
}

func (d *BillingDAL) GetDashboardSummary(_ context.Context) (*DashboardSnapshot, error) {
	return &DashboardSnapshot{
		OpenBills:       128,
		PendingInvoices: 37,
		OverdueAmount:   482130.50,
		AutoMatchedRate: 92.4,
		RecentActivity: []vo.ActivityItem{
			{
				Title:       "Northwind April bill imported",
				Description: "4,216 line items were normalized and queued for review.",
				Time:        "09:20",
				Status:      "success",
			},
			{
				Title:       "2 reconciliation groups need confirmation",
				Description: "Bank feed and invoice ledger mismatch exceeded tolerance.",
				Time:        "08:40",
				Status:      "warning",
			},
			{
				Title:       "Auto-posting window started",
				Description: "Scheduled posting is processing approved invoices.",
				Time:        "08:00",
				Status:      "processing",
			},
		},
		AttentionList: []vo.AttentionItem{
			{
				Title:    "Government account overdue review",
				Owner:    "Avery",
				DueDate:  "Today 15:00",
				Progress: 68,
			},
			{
				Title:    "Q2 pre-close billing lock",
				Owner:    "Nina",
				DueDate:  "Tomorrow 10:30",
				Progress: 44,
			},
			{
				Title:    "Tax inclusive invoice spot check",
				Owner:    "Leo",
				DueDate:  "Apr 28",
				Progress: 82,
			},
		},
	}, nil
}

func (d *BillingDAL) ListBills(_ context.Context) ([]vo.Bill, error) {
	return []vo.Bill{
		{
			ID:           "CB-2026-0418",
			CustomerName: "Northwind Logistics",
			BillingMonth: "2026-04",
			Amount:       126500.00,
			Status:       "ready",
			Channel:      "EDI",
			UpdatedAt:    "2026-04-24 09:18",
		},
		{
			ID:           "CB-2026-0415",
			CustomerName: "Helio Retail Group",
			BillingMonth: "2026-04",
			Amount:       87420.35,
			Status:       "reviewing",
			Channel:      "Portal",
			UpdatedAt:    "2026-04-24 08:46",
		},
		{
			ID:           "CB-2026-0409",
			CustomerName: "Blue Peak Energy",
			BillingMonth: "2026-04",
			Amount:       192880.90,
			Status:       "overdue",
			Channel:      "Email",
			UpdatedAt:    "2026-04-23 18:05",
		},
		{
			ID:           "CB-2026-0402",
			CustomerName: "Nova Health Labs",
			BillingMonth: "2026-03",
			Amount:       45670.00,
			Status:       "paid",
			Channel:      "Portal",
			UpdatedAt:    "2026-04-22 16:12",
		},
	}, nil
}

func (d *BillingDAL) ListCustomers(_ context.Context) ([]vo.Customer, error) {
	return []vo.Customer{
		{
			ID:                "CUS-101",
			Name:              "Northwind Logistics",
			CreditLevel:       "A",
			ActiveContracts:   12,
			OutstandingAmount: 126500,
			BillingHealth:     91,
			PrimaryContact:    "Olivia Chen",
		},
		{
			ID:                "CUS-204",
			Name:              "Helio Retail Group",
			CreditLevel:       "B+",
			ActiveContracts:   8,
			OutstandingAmount: 87420.35,
			BillingHealth:     76,
			PrimaryContact:    "Mason Wu",
		},
		{
			ID:                "CUS-318",
			Name:              "Blue Peak Energy",
			CreditLevel:       "A-",
			ActiveContracts:   5,
			OutstandingAmount: 192880.90,
			BillingHealth:     58,
			PrimaryContact:    "Sophia Lin",
		},
	}, nil
}

func (d *BillingDAL) ListReconciliationTasks(_ context.Context) ([]vo.ReconciliationTask, error) {
	return []vo.ReconciliationTask{
		{
			ID:      "REC-2401",
			Bank:    "Industrial Bank",
			Period:  "2026-04-24",
			Matched: 182,
			Total:   197,
			Status:  "running",
			Owner:   "Finance Ops",
		},
		{
			ID:      "REC-2398",
			Bank:    "ICBC",
			Period:  "2026-04-23",
			Matched: 210,
			Total:   210,
			Status:  "done",
			Owner:   "Finance Ops",
		},
		{
			ID:      "REC-2396",
			Bank:    "Bank of China",
			Period:  "2026-04-23",
			Matched: 96,
			Total:   122,
			Status:  "attention",
			Owner:   "Settlement Team",
		},
	}, nil
}
