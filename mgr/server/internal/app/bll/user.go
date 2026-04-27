package bll

import (
	"context"
	"errors"

	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/dal"
	"clearbill/mgr/server/internal/app/dal/dbmodel"
	"clearbill/mgr/server/pkg/passwordx"
)

type UserService struct {
	UserDAL   *dal.UserDAL
	TenantDAL *dal.TenantDAL
}

func NewUserService(userDAL *dal.UserDAL, tenantDAL *dal.TenantDAL) *UserService {
	return &UserService{
		UserDAL:   userDAL,
		TenantDAL: tenantDAL,
	}
}

func (s *UserService) ListUsers(ctx context.Context, actor *dbmodel.User, req *vo.ListUserReq) ([]vo.User, error) {
	tenantID := req.TenantID
	if actor.Role == dbmodel.RoleTenantAdmin {
		tenantID = actor.TenantID
	} else if actor.Role != dbmodel.RoleSysadmin {
		return nil, errors.New("permission denied")
	}

	users, err := s.UserDAL.List(ctx, req.Keyword, tenantID)
	if err != nil {
		return nil, err
	}

	result := make([]vo.User, 0, len(users))
	for i := range users {
		result = append(result, ToUserVO(&users[i]))
	}
	return result, nil
}

func (s *UserService) CreateUser(ctx context.Context, actor *dbmodel.User, req *vo.CreateUserReq) (*vo.CreateUserResp, error) {
	tenantID, role, err := s.normalizeCreateScope(ctx, actor, req.TenantID, req.Role)
	if err != nil {
		return nil, err
	}

	hash, err := passwordx.HashPassword(dbmodel.DefaultUserPassword)
	if err != nil {
		return nil, err
	}

	user := &dbmodel.User{
		Username:     req.Username,
		DisplayName:  req.DisplayName,
		PasswordHash: hash,
		Role:         role,
		TenantID:     tenantID,
		Status:       normalizeStatus(req.Status),
	}
	if err := s.UserDAL.Create(ctx, user); err != nil {
		return nil, err
	}

	return &vo.CreateUserResp{
		User:            ToUserVO(user),
		InitialPassword: dbmodel.DefaultUserPassword,
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, actor *dbmodel.User, id uint) (*vo.User, error) {
	user, err := s.UserDAL.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !canManageUser(actor, user) && actor.ID != user.ID {
		return nil, errors.New("permission denied")
	}
	result := ToUserVO(user)
	return &result, nil
}

func (s *UserService) UpdateUser(ctx context.Context, actor *dbmodel.User, id uint, req *vo.UpdateUserReq) (*vo.User, error) {
	user, err := s.UserDAL.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !canManageUser(actor, user) {
		return nil, errors.New("permission denied")
	}

	role := req.Role
	tenantID := req.TenantID
	if actor.Role == dbmodel.RoleTenantAdmin {
		role = restrictRole(role)
		tenantID = actor.TenantID
	}
	if role == dbmodel.RoleSysadmin && actor.Role != dbmodel.RoleSysadmin {
		return nil, errors.New("permission denied")
	}
	if role != dbmodel.RoleSysadmin {
		if tenantID == nil {
			return nil, errors.New("tenantId is required")
		}
		if _, err := s.TenantDAL.GetByID(ctx, *tenantID); err != nil {
			return nil, err
		}
	} else {
		tenantID = nil
	}

	user.DisplayName = req.DisplayName
	user.Role = role
	user.TenantID = tenantID
	user.Status = normalizeStatus(req.Status)
	if err := s.UserDAL.Update(ctx, user); err != nil {
		return nil, err
	}
	result := ToUserVO(user)
	return &result, nil
}

func (s *UserService) DeleteUser(ctx context.Context, actor *dbmodel.User, id uint) error {
	user, err := s.UserDAL.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if actor.ID == user.ID {
		return errors.New("cannot delete current user")
	}
	if !canManageUser(actor, user) {
		return errors.New("permission denied")
	}
	return s.UserDAL.Delete(ctx, id)
}

func (s *UserService) ResetPassword(ctx context.Context, actor *dbmodel.User, id uint, req *vo.ResetPasswordReq) error {
	user, err := s.UserDAL.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !canManageUser(actor, user) {
		return errors.New("permission denied")
	}
	hash, err := passwordx.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	return s.UserDAL.Update(ctx, user)
}

func (s *UserService) normalizeCreateScope(ctx context.Context, actor *dbmodel.User, tenantID *uint, role string) (*uint, string, error) {
	switch actor.Role {
	case dbmodel.RoleSysadmin:
		if role == dbmodel.RoleSysadmin {
			return nil, role, nil
		}
		if tenantID == nil {
			return nil, "", errors.New("tenantId is required")
		}
		if _, err := s.TenantDAL.GetByID(ctx, *tenantID); err != nil {
			return nil, "", err
		}
		return tenantID, restrictRole(role), nil
	case dbmodel.RoleTenantAdmin:
		if actor.TenantID == nil {
			return nil, "", errors.New("tenant admin missing tenant scope")
		}
		return actor.TenantID, restrictRole(role), nil
	default:
		return nil, "", errors.New("permission denied")
	}
}

func canManageUser(actor, target *dbmodel.User) bool {
	if actor.Role == dbmodel.RoleSysadmin {
		return true
	}
	if actor.Role != dbmodel.RoleTenantAdmin {
		return false
	}
	if actor.TenantID == nil || target.TenantID == nil {
		return false
	}
	return *actor.TenantID == *target.TenantID && target.Role != dbmodel.RoleSysadmin
}

func restrictRole(role string) string {
	if role == dbmodel.RoleTenantAdmin {
		return dbmodel.RoleTenantAdmin
	}
	return dbmodel.RoleUser
}

func normalizeStatus(status string) string {
	if status == dbmodel.StatusDisabled {
		return dbmodel.StatusDisabled
	}
	return dbmodel.StatusActive
}

func ToUserVO(user *dbmodel.User) vo.User {
	return vo.User{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Role:        user.Role,
		TenantID:    user.TenantID,
		Status:      user.Status,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
