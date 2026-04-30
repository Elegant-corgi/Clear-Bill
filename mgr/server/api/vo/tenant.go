package vo

import "time"

type Tenant struct {
	ID               uint      `json:"id"`
	Code             string    `json:"code"`
	Name             string    `json:"name"`
	AdminUsername    string    `json:"adminUsername"`
	AdminDisplayName string    `json:"adminDisplayName"`
	ContactName      string    `json:"contactName"`
	ContactPhone     string    `json:"contactPhone"`
	Status           string    `json:"status"`
	Remark           string    `json:"remark"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type CreateTenantReq struct {
	Code             string `json:"code" binding:"required,max=64"`
	Name             string `json:"name" binding:"required,max=128"`
	AdminUsername    string `json:"adminUsername" binding:"max=64"`
	AdminDisplayName string `json:"adminDisplayName" binding:"max=128"`
	ContactName      string `json:"contactName" binding:"max=64"`
	ContactPhone     string `json:"contactPhone" binding:"max=32"`
	Status           string `json:"status" binding:"max=32"`
	Remark           string `json:"remark" binding:"max=255"`
}

type CreateTenantResp struct {
	Tenant          Tenant `json:"tenant"`
	AdminUsername   string `json:"adminUsername"`
	InitialPassword string `json:"initialPassword"`
}

type UpdateTenantReq struct {
	Code         string `json:"code" binding:"required,max=64"`
	Name         string `json:"name" binding:"required,max=128"`
	ContactName  string `json:"contactName" binding:"max=64"`
	ContactPhone string `json:"contactPhone" binding:"max=32"`
	Status       string `json:"status" binding:"max=32"`
	Remark       string `json:"remark" binding:"max=255"`
}

type ListTenantReq struct {
	PageReq
	Keyword string `form:"keyword"`
}
