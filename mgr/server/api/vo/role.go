package vo

import "time"

type Permission struct {
	ID     string `json:"id"`
	Method string `json:"method"`
	Path   string `json:"path"`
	Tag    string `json:"tag"`
}

type Role struct {
	ID            uint      `json:"id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	Scope         string    `json:"scope"`
	TenantID      *uint     `json:"tenantId,omitempty"`
	Builtin       bool      `json:"builtin"`
	PermissionIDs []string  `json:"permissionIds"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type CreateRoleReq struct {
	Code          string   `json:"code" binding:"required,max=64"`
	Name          string   `json:"name" binding:"required,max=128"`
	Scope         string   `json:"scope" binding:"max=16"`
	TenantID      *uint    `json:"tenantId"`
	PermissionIDs []string `json:"permissionIds"`
}

type UpdateRoleReq struct {
	Name string `json:"name" binding:"required,max=128"`
}

type UpdateRolePermissionsReq struct {
	PermissionIDs []string `json:"permissionIds"`
}

type ListRoleReq struct {
	Keyword  string `form:"keyword"`
	Scope    string `form:"scope"`
	TenantID *uint  `form:"tenantId"`
}
