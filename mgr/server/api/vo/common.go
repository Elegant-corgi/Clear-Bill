package vo

type ResponseResult struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

const OKStatus = "ok"

type StatusResult struct {
	Status string `json:"status"`
}

type ListResult struct {
	List any `json:"list"`
}

type HealthStatus struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}
