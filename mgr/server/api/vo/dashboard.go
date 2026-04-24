package vo

type DashboardSummary struct {
	OpenBills       int             `json:"openBills"`
	PendingInvoices int             `json:"pendingInvoices"`
	OverdueAmount   float64         `json:"overdueAmount"`
	AutoMatchedRate float64         `json:"autoMatchedRate"`
	RecentActivity  []ActivityItem  `json:"recentActivity"`
	AttentionList   []AttentionItem `json:"attentionList"`
}

type ActivityItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Time        string `json:"time"`
	Status      string `json:"status"`
}

type AttentionItem struct {
	Title    string `json:"title"`
	Owner    string `json:"owner"`
	DueDate  string `json:"dueDate"`
	Progress int    `json:"progress"`
}
