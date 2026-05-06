package vo

import "time"

type Credential struct {
	ID              uint       `json:"id"`
	UserID          uint       `json:"userId"`
	Type            string     `json:"type"`
	Name            string     `json:"name"`
	Status          string     `json:"status"`
	AccessKey       string     `json:"accessKey,omitempty"`
	AccessKeyPreview string    `json:"accessKeyPreview,omitempty"`
	SecretKey       string     `json:"secretKey,omitempty"`
	Token           string     `json:"token,omitempty"`
	TokenPreview    string     `json:"tokenPreview,omitempty"`
	ParentID        *uint      `json:"parentId,omitempty"`
	ExpiresAt       *time.Time `json:"expiresAt,omitempty"`
	RotatedAt       *time.Time `json:"rotatedAt,omitempty"`
	LastUsedAt      *time.Time `json:"lastUsedAt,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

type CreateCredentialReq struct {
	Type string `json:"type" binding:"required,max=16"`
	Name string `json:"name" binding:"required,max=128"`
}

type RotateCredentialReq struct {
	Name string `json:"name" binding:"omitempty,max=128"`
}

type ListCredentialReq struct {
	PageReq
	Type   string `form:"type"`
	Status string `form:"status"`
}
