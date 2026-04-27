package bll

import (
	"context"
	"fmt"

	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/dal"
	"clearbill/mgr/server/internal/app/dal/dbmodel"
	"clearbill/mgr/server/pkg/passwordx"
)

type TenantService struct {
	TenantDAL *dal.TenantDAL
	UserDAL   *dal.UserDAL
}

func NewTenantService(tenantDAL *dal.TenantDAL, userDAL *dal.UserDAL) *TenantService {
	return &TenantService{
		TenantDAL: tenantDAL,
		UserDAL:   userDAL,
	}
}

func (s *TenantService) CreateTenant(ctx context.Context, req *vo.CreateTenantReq) (*vo.CreateTenantResp, error) {
	tenant := &dbmodel.Tenant{
		Code:         req.Code,
		Name:         req.Name,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Status:       req.Status,
		Remark:       req.Remark,
	}
	if tenant.Status == "" {
		tenant.Status = "active"
	}

	if err := s.TenantDAL.Create(ctx, tenant); err != nil {
		return nil, err
	}

	adminUsername := req.AdminUsername
	if adminUsername == "" {
		adminUsername = fmt.Sprintf("%s_admin", req.Code)
	}
	adminDisplayName := req.AdminDisplayName
	if adminDisplayName == "" {
		adminDisplayName = fmt.Sprintf("%s Admin", req.Name)
	}
	hash, err := passwordx.HashPassword(dbmodel.DefaultUserPassword)
	if err != nil {
		return nil, err
	}
	admin := &dbmodel.User{
		Username:     adminUsername,
		DisplayName:  adminDisplayName,
		PasswordHash: hash,
		Role:         dbmodel.RoleTenantAdmin,
		TenantID:     &tenant.ID,
		Status:       dbmodel.StatusActive,
	}
	if err := s.UserDAL.Create(ctx, admin); err != nil {
		_ = s.TenantDAL.Delete(ctx, tenant.ID)
		return nil, err
	}

	return &vo.CreateTenantResp{
		Tenant:          *toTenantVO(tenant),
		AdminUsername:   adminUsername,
		InitialPassword: dbmodel.DefaultUserPassword,
	}, nil
}

func (s *TenantService) ListTenants(ctx context.Context, req *vo.ListTenantReq) ([]vo.Tenant, error) {
	tenants, err := s.TenantDAL.List(ctx, req.Keyword)
	if err != nil {
		return nil, err
	}

	result := make([]vo.Tenant, 0, len(tenants))
	for i := range tenants {
		result = append(result, *toTenantVO(&tenants[i]))
	}

	return result, nil
}

func (s *TenantService) GetTenant(ctx context.Context, id uint) (*vo.Tenant, error) {
	tenant, err := s.TenantDAL.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toTenantVO(tenant), nil
}

func (s *TenantService) UpdateTenant(ctx context.Context, id uint, req *vo.UpdateTenantReq) (*vo.Tenant, error) {
	tenant, err := s.TenantDAL.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	tenant.Code = req.Code
	tenant.Name = req.Name
	tenant.ContactName = req.ContactName
	tenant.ContactPhone = req.ContactPhone
	tenant.Status = req.Status
	tenant.Remark = req.Remark
	if tenant.Status == "" {
		tenant.Status = "active"
	}

	if err := s.TenantDAL.Update(ctx, tenant); err != nil {
		return nil, err
	}

	return toTenantVO(tenant), nil
}

func (s *TenantService) DeleteTenant(ctx context.Context, id uint) error {
	_, err := s.TenantDAL.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.TenantDAL.Delete(ctx, id)
}

func toTenantVO(tenant *dbmodel.Tenant) *vo.Tenant {
	return &vo.Tenant{
		ID:           tenant.ID,
		Code:         tenant.Code,
		Name:         tenant.Name,
		ContactName:  tenant.ContactName,
		ContactPhone: tenant.ContactPhone,
		Status:       tenant.Status,
		Remark:       tenant.Remark,
		CreatedAt:    tenant.CreatedAt,
		UpdatedAt:    tenant.UpdatedAt,
	}
}
