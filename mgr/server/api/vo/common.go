package vo

type ResponseResult struct {
	Success      bool   `json:"success"`
	Data         any    `json:"data,omitempty"`
	Error        string `json:"error,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

type PageReq struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

func (r PageReq) Normalize() PageReq {
	if r.Page <= 0 {
		r.Page = DefaultPage
	}
	if r.PageSize <= 0 {
		r.PageSize = DefaultPageSize
	}
	if r.PageSize > MaxPageSize {
		r.PageSize = MaxPageSize
	}
	return r
}

type PageResult[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
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
