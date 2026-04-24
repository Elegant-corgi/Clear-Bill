package vo

type Bill struct {
	ID           string  `json:"id"`
	CustomerName string  `json:"customerName"`
	BillingMonth string  `json:"billingMonth"`
	Amount       float64 `json:"amount"`
	Status       string  `json:"status"`
	Channel      string  `json:"channel"`
	UpdatedAt    string  `json:"updatedAt"`
}

type Customer struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	CreditLevel       string  `json:"creditLevel"`
	ActiveContracts   int     `json:"activeContracts"`
	OutstandingAmount float64 `json:"outstandingAmount"`
	BillingHealth     int     `json:"billingHealth"`
	PrimaryContact    string  `json:"primaryContact"`
}

type ReconciliationTask struct {
	ID      string `json:"id"`
	Bank    string `json:"bank"`
	Period  string `json:"period"`
	Matched int    `json:"matched"`
	Total   int    `json:"total"`
	Status  string `json:"status"`
	Owner   string `json:"owner"`
}
