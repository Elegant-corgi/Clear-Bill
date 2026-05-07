package vo

import "time"

type AuditLog struct {
	ID         uint      `json:"id"`
	User       string    `json:"user"`
	Operation  string    `json:"operation"`
	OccurredAt time.Time `json:"occurredAt"`
	Resource   string    `json:"resource"`
	Result     string    `json:"result"`
}

type ListAuditLogReq struct {
	PageReq
	User      string `form:"user"`
	Operation string `form:"operation"`
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
}
