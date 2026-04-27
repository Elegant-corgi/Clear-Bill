package vo

import "time"

type User struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"displayName"`
	Role        string    `json:"role"`
	TenantID    *uint     `json:"tenantId,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required,max=64"`
	Password string `json:"password" binding:"required,max=128"`
}

type LoginResp struct {
	User User `json:"user"`
}

type CreateAPITokenReq struct {
	Name string `json:"name" binding:"required,max=128"`
}

type CreateAPITokenResp struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type ChangeOwnPasswordReq struct {
	OldPassword string `json:"oldPassword" binding:"required,max=128"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=128"`
}

type ResetPasswordReq struct {
	NewPassword string `json:"newPassword" binding:"required,min=6,max=128"`
}

type CreateUserReq struct {
	Username    string `json:"username" binding:"required,max=64"`
	DisplayName string `json:"displayName" binding:"required,max=128"`
	Role        string `json:"role" binding:"required,max=64"`
	TenantID    *uint  `json:"tenantId"`
	Status      string `json:"status" binding:"max=32"`
}

type CreateUserResp struct {
	User            User   `json:"user"`
	InitialPassword string `json:"initialPassword"`
}

type UpdateUserReq struct {
	DisplayName string `json:"displayName" binding:"required,max=128"`
	Role        string `json:"role" binding:"required,max=64"`
	TenantID    *uint  `json:"tenantId"`
	Status      string `json:"status" binding:"max=32"`
}

type ListUserReq struct {
	Keyword  string `form:"keyword"`
	TenantID *uint  `form:"tenantId"`
}
