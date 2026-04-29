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
	UserDAL     *dal.UserDAL
	TenantDAL   *dal.TenantDAL
	RoleService *RoleService
}

func NewUserService(userDAL *dal.UserDAL, tenantDAL *dal.TenantDAL, roleService *RoleService) *UserService {
	return &UserService{
		UserDAL:     userDAL,
		TenantDAL:   tenantDAL,
		RoleService: roleService,
	}
}

func (s *UserService) ListUsers(ctx context.Context, actor *dbmodel.User, req *vo.ListUserReq) ([]vo.User, error) {
	role, err := s.RoleService.GetActorRole(ctx, actor)
	if err != nil {
		return nil, err
	}

	tenantID := req.TenantID
	switch role.Scope {
	case dbmodel.RoleScopeSystem:
	case dbmodel.RoleScopeTenant:
		tenantID = actor.TenantID
	default:
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
	if containsChineseCharacters(req.Username) {
		return nil, errors.New("登录账号不能包含中文")
	}

	role, tenantID, err := s.RoleService.ValidateRoleAssignment(ctx, actor, req.Role, req.TenantID)
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
		Role:         role.Code,
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
	canManage, err := s.canManageUser(ctx, actor, user)
	if err != nil {
		return nil, err
	}
	if !canManage && actor.ID != user.ID {
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
	canManage, err := s.canManageUser(ctx, actor, user)
	if err != nil {
		return nil, err
	}
	if !canManage {
		return nil, errors.New("permission denied")
	}

	role, tenantID, err := s.RoleService.ValidateRoleAssignment(ctx, actor, req.Role, req.TenantID)
	if err != nil {
		return nil, err
	}

	user.DisplayName = req.DisplayName
	user.Role = role.Code
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
	canManage, err := s.canManageUser(ctx, actor, user)
	if err != nil {
		return err
	}
	if !canManage {
		return errors.New("permission denied")
	}
	return s.UserDAL.Delete(ctx, id)
}

func (s *UserService) ResetPassword(ctx context.Context, actor *dbmodel.User, id uint, req *vo.ResetPasswordReq) error {
	user, err := s.UserDAL.GetByID(ctx, id)
	if err != nil {
		return err
	}
	canManage, err := s.canManageUser(ctx, actor, user)
	if err != nil {
		return err
	}
	if !canManage {
		return errors.New("permission denied")
	}
	if err := passwordx.ComparePassword(user.PasswordHash, req.NewPassword); err == nil {
		return errors.New("new password must be different from old password")
	}
	hash, err := passwordx.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	return s.UserDAL.Update(ctx, user)
}

func (s *UserService) canManageUser(ctx context.Context, actor, target *dbmodel.User) (bool, error) {
	role, err := s.RoleService.GetActorRole(ctx, actor)
	if err != nil {
		return false, err
	}
	switch role.Scope {
	case dbmodel.RoleScopeSystem:
		return true, nil
	case dbmodel.RoleScopeTenant:
		if actor.TenantID == nil || target.TenantID == nil {
			return false, nil
		}
		return *actor.TenantID == *target.TenantID && target.Role != dbmodel.RoleSysadmin, nil
	default:
		return false, nil
	}
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
