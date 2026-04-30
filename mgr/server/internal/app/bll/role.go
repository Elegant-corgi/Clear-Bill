package bll

import (
	"context"
	"errors"
	"sort"
	"strings"

	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/dal"
	"clearbill/mgr/server/internal/app/dal/dbmodel"
)

type RoleService struct {
	RoleDAL           *dal.RoleDAL
	RolePermissionDAL *dal.RolePermissionDAL
	UserDAL           *dal.UserDAL
	TenantDAL         *dal.TenantDAL
	PermissionCatalog *PermissionCatalog
}

func NewRoleService(
	roleDAL *dal.RoleDAL,
	rolePermissionDAL *dal.RolePermissionDAL,
	userDAL *dal.UserDAL,
	tenantDAL *dal.TenantDAL,
	permissionCatalog *PermissionCatalog,
) *RoleService {
	return &RoleService{
		RoleDAL:           roleDAL,
		RolePermissionDAL: rolePermissionDAL,
		UserDAL:           userDAL,
		TenantDAL:         tenantDAL,
		PermissionCatalog: permissionCatalog,
	}
}

func (s *RoleService) ListPermissions() []vo.Permission {
	items := s.PermissionCatalog.List()
	result := make([]vo.Permission, 0, len(items))
	for _, item := range items {
		result = append(result, vo.Permission{
			ID:     item.ID,
			Method: item.Method,
			Path:   item.Path,
			Tag:    item.Tag,
		})
	}
	return result
}

func (s *RoleService) ResolvePermissionByRoute(method, path string) (PermissionItem, bool) {
	return s.PermissionCatalog.GetByRoute(method, path)
}

func (s *RoleService) HasPermission(ctx context.Context, actor *dbmodel.User, permissionID string) (bool, error) {
	if actor == nil {
		return false, nil
	}
	if !s.PermissionCatalog.Exists(permissionID) {
		return false, nil
	}

	role, err := s.GetActorRole(ctx, actor)
	if err != nil {
		return false, err
	}
	permissionIDs, err := s.permissionIDsForRole(ctx, role)
	if err != nil {
		return false, err
	}
	for _, item := range permissionIDs {
		if item == permissionID {
			return true, nil
		}
	}
	return false, nil
}

func (s *RoleService) GetActorRole(ctx context.Context, actor *dbmodel.User) (*dbmodel.Role, error) {
	if actor == nil {
		return nil, errors.New("unauthorized")
	}
	return s.RoleDAL.GetByCode(ctx, actor.Role)
}

func (s *RoleService) IsSystemScoped(ctx context.Context, actor *dbmodel.User) (bool, error) {
	role, err := s.GetActorRole(ctx, actor)
	if err != nil {
		return false, err
	}
	return role.Scope == dbmodel.RoleScopeSystem, nil
}

func (s *RoleService) ListRoles(ctx context.Context, actor *dbmodel.User, req *vo.ListRoleReq) (*vo.PageResult[vo.Role], error) {
	actorRole, err := s.GetActorRole(ctx, actor)
	if err != nil {
		return nil, err
	}

	pageReq := req.PageReq.Normalize()
	roleScope := strings.TrimSpace(req.Scope)
	tenantID := req.TenantID
	switch actorRole.Scope {
	case dbmodel.RoleScopeSystem:
	case dbmodel.RoleScopeTenant:
		if actor.TenantID == nil {
			return nil, errors.New("permission denied")
		}
		roleScope = dbmodel.RoleScopeTenant
		tenantID = actor.TenantID
	default:
		return nil, errors.New("permission denied")
	}

	roles, query, total, err := s.RoleDAL.List(ctx, req.Keyword, roleScope, tenantID, pageReq.Page, pageReq.PageSize)
	if err != nil {
		return nil, err
	}

	visibleRoles := make([]dbmodel.Role, 0, len(roles))
	for _, role := range roles {
		if !s.canViewRole(actorRole, actor, &role, req) {
			continue
		}
		visibleRoles = append(visibleRoles, role)
	}
	items, err := s.toRoleVOs(ctx, visibleRoles)
	if err != nil {
		return nil, err
	}
	return &vo.PageResult[vo.Role]{
		List:     items,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

func (s *RoleService) GetRole(ctx context.Context, actor *dbmodel.User, id uint) (*vo.Role, error) {
	actorRole, err := s.GetActorRole(ctx, actor)
	if err != nil {
		return nil, err
	}

	role, err := s.RoleDAL.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !s.canAccessRole(actorRole, actor, role) {
		return nil, errors.New("permission denied")
	}

	items, err := s.toRoleVOs(ctx, []dbmodel.Role{*role})
	if err != nil {
		return nil, err
	}
	return &items[0], nil
}

func (s *RoleService) CreateRole(ctx context.Context, actor *dbmodel.User, req *vo.CreateRoleReq) (*vo.Role, error) {
	actorRole, err := s.GetActorRole(ctx, actor)
	if err != nil {
		return nil, err
	}

	roleCode := normalizeRoleCode(req.Code)
	if roleCode == "" {
		return nil, errors.New("role code is required")
	}
	roleName := strings.TrimSpace(req.Name)
	if roleName == "" {
		return nil, errors.New("role name is required")
	}

	permissionIDs, err := s.normalizePermissionIDs(req.PermissionIDs)
	if err != nil {
		return nil, err
	}

	role := &dbmodel.Role{
		Code:    roleCode,
		Name:    roleName,
		Builtin: false,
	}
	switch actorRole.Scope {
	case dbmodel.RoleScopeSystem:
		role.Scope = normalizeRoleScope(req.Scope, req.TenantID)
		if role.Scope == dbmodel.RoleScopeTenant {
			if req.TenantID == nil {
				return nil, errors.New("tenantId is required")
			}
			if _, err := s.TenantDAL.GetByID(ctx, *req.TenantID); err != nil {
				return nil, err
			}
			role.TenantID = req.TenantID
		}
	case dbmodel.RoleScopeTenant:
		if actor.TenantID == nil {
			return nil, errors.New("tenant admin missing tenant scope")
		}
		role.Scope = dbmodel.RoleScopeTenant
		role.TenantID = actor.TenantID
	default:
		return nil, errors.New("permission denied")
	}

	if err := s.RoleDAL.Create(ctx, role); err != nil {
		return nil, err
	}
	if err := s.RolePermissionDAL.Replace(ctx, role.ID, permissionIDs); err != nil {
		_ = s.RoleDAL.Delete(ctx, role.ID)
		return nil, err
	}

	items, err := s.toRoleVOs(ctx, []dbmodel.Role{*role})
	if err != nil {
		return nil, err
	}
	return &items[0], nil
}

func (s *RoleService) UpdateRole(ctx context.Context, actor *dbmodel.User, id uint, req *vo.UpdateRoleReq) (*vo.Role, error) {
	role, err := s.RoleDAL.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role.Builtin {
		return nil, errors.New("builtin role cannot be modified")
	}

	actorRole, err := s.GetActorRole(ctx, actor)
	if err != nil {
		return nil, err
	}
	if !s.canManageRole(actorRole, actor, role) {
		return nil, errors.New("permission denied")
	}

	roleName := strings.TrimSpace(req.Name)
	if roleName == "" {
		return nil, errors.New("role name is required")
	}
	role.Name = roleName
	if err := s.RoleDAL.Update(ctx, role); err != nil {
		return nil, err
	}

	items, err := s.toRoleVOs(ctx, []dbmodel.Role{*role})
	if err != nil {
		return nil, err
	}
	return &items[0], nil
}

func (s *RoleService) UpdateRolePermissions(ctx context.Context, actor *dbmodel.User, id uint, req *vo.UpdateRolePermissionsReq) (*vo.Role, error) {
	role, err := s.RoleDAL.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role.Builtin {
		return nil, errors.New("builtin role cannot be modified")
	}

	actorRole, err := s.GetActorRole(ctx, actor)
	if err != nil {
		return nil, err
	}
	if !s.canManageRole(actorRole, actor, role) {
		return nil, errors.New("permission denied")
	}

	permissionIDs, err := s.normalizePermissionIDs(req.PermissionIDs)
	if err != nil {
		return nil, err
	}
	if err := s.RolePermissionDAL.Replace(ctx, role.ID, permissionIDs); err != nil {
		return nil, err
	}

	items, err := s.toRoleVOs(ctx, []dbmodel.Role{*role})
	if err != nil {
		return nil, err
	}
	return &items[0], nil
}

func (s *RoleService) DeleteRole(ctx context.Context, actor *dbmodel.User, id uint) error {
	role, err := s.RoleDAL.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if role.Builtin {
		return errors.New("builtin role cannot be modified")
	}

	actorRole, err := s.GetActorRole(ctx, actor)
	if err != nil {
		return err
	}
	if !s.canManageRole(actorRole, actor, role) {
		return errors.New("permission denied")
	}

	count, err := s.UserDAL.CountByRole(ctx, role.Code)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("role is in use")
	}

	if err := s.RolePermissionDAL.DeleteByRoleID(ctx, role.ID); err != nil {
		return err
	}
	return s.RoleDAL.Delete(ctx, role.ID)
}

func (s *RoleService) ValidateRoleAssignment(ctx context.Context, actor *dbmodel.User, roleCode string, tenantID *uint) (*dbmodel.Role, *uint, error) {
	role, err := s.RoleDAL.GetByCode(ctx, normalizeRoleCode(roleCode))
	if err != nil {
		return nil, nil, err
	}
	actorRole, err := s.GetActorRole(ctx, actor)
	if err != nil {
		return nil, nil, err
	}

	switch actorRole.Scope {
	case dbmodel.RoleScopeSystem:
		return s.validateSystemRoleAssignment(ctx, role, tenantID)
	case dbmodel.RoleScopeTenant:
		if actor.TenantID == nil {
			return nil, nil, errors.New("tenant admin missing tenant scope")
		}
		if role.Scope != dbmodel.RoleScopeTenant {
			return nil, nil, errors.New("permission denied")
		}
		if role.TenantID != nil && *role.TenantID != *actor.TenantID {
			return nil, nil, errors.New("permission denied")
		}
		return role, actor.TenantID, nil
	default:
		return nil, nil, errors.New("permission denied")
	}
}

func (s *RoleService) normalizePermissionIDs(permissionIDs []string) ([]string, error) {
	seen := make(map[string]struct{}, len(permissionIDs))
	result := make([]string, 0, len(permissionIDs))
	for _, permissionID := range permissionIDs {
		permissionID = strings.TrimSpace(permissionID)
		if permissionID == "" {
			continue
		}
		if !s.PermissionCatalog.Exists(permissionID) {
			return nil, errors.New("invalid permission id")
		}
		if _, ok := seen[permissionID]; ok {
			continue
		}
		seen[permissionID] = struct{}{}
		result = append(result, permissionID)
	}
	sort.Strings(result)
	return result, nil
}

func (s *RoleService) permissionIDsForRole(ctx context.Context, role *dbmodel.Role) ([]string, error) {
	if role.Builtin {
		return s.builtinPermissionIDs(role.Code), nil
	}

	items, err := s.RolePermissionDAL.ListByRoleIDs(ctx, []uint{role.ID})
	if err != nil {
		return nil, err
	}
	return items[role.ID], nil
}

func (s *RoleService) builtinPermissionIDs(roleCode string) []string {
	allowed := make([]string, 0, len(s.PermissionCatalog.items))
	for _, item := range s.PermissionCatalog.List() {
		switch roleCode {
		case dbmodel.RoleSysadmin:
			allowed = append(allowed, item.ID)
		case dbmodel.RoleTenantAdmin:
			if tenantAdminPermissionSet[item.ID] {
				allowed = append(allowed, item.ID)
			}
		case dbmodel.RoleUser:
			if userPermissionSet[item.ID] {
				allowed = append(allowed, item.ID)
			}
		}
	}
	return allowed
}

func (s *RoleService) toRoleVOs(ctx context.Context, roles []dbmodel.Role) ([]vo.Role, error) {
	customRoleIDs := make([]uint, 0, len(roles))
	for _, role := range roles {
		if !role.Builtin {
			customRoleIDs = append(customRoleIDs, role.ID)
		}
	}

	customPermissions, err := s.RolePermissionDAL.ListByRoleIDs(ctx, customRoleIDs)
	if err != nil {
		return nil, err
	}

	result := make([]vo.Role, 0, len(roles))
	for _, role := range roles {
		permissionIDs := s.builtinPermissionIDs(role.Code)
		if !role.Builtin {
			permissionIDs = customPermissions[role.ID]
		}
		result = append(result, vo.Role{
			ID:            role.ID,
			Code:          role.Code,
			Name:          role.Name,
			Scope:         role.Scope,
			TenantID:      role.TenantID,
			Builtin:       role.Builtin,
			PermissionIDs: permissionIDs,
			CreatedAt:     role.CreatedAt,
			UpdatedAt:     role.UpdatedAt,
		})
	}
	return result, nil
}

func (s *RoleService) canViewRole(actorRole *dbmodel.Role, actor *dbmodel.User, role *dbmodel.Role, req *vo.ListRoleReq) bool {
	if req.Scope != "" && role.Scope != req.Scope {
		return false
	}

	switch actorRole.Scope {
	case dbmodel.RoleScopeSystem:
		if req.TenantID != nil {
			if role.Scope != dbmodel.RoleScopeTenant {
				return false
			}
			return role.TenantID == nil || *role.TenantID == *req.TenantID
		}
		return true
	case dbmodel.RoleScopeTenant:
		if actor.TenantID == nil || role.Scope != dbmodel.RoleScopeTenant {
			return false
		}
		if req.TenantID != nil && *req.TenantID != *actor.TenantID {
			return false
		}
		return role.TenantID == nil || *role.TenantID == *actor.TenantID
	default:
		return false
	}
}

func (s *RoleService) canAccessRole(actorRole *dbmodel.Role, actor *dbmodel.User, role *dbmodel.Role) bool {
	return s.canViewRole(actorRole, actor, role, &vo.ListRoleReq{})
}

func (s *RoleService) canManageRole(actorRole *dbmodel.Role, actor *dbmodel.User, role *dbmodel.Role) bool {
	switch actorRole.Scope {
	case dbmodel.RoleScopeSystem:
		return true
	case dbmodel.RoleScopeTenant:
		if actor.TenantID == nil || role.Scope != dbmodel.RoleScopeTenant {
			return false
		}
		return role.TenantID == nil || *role.TenantID == *actor.TenantID
	default:
		return false
	}
}

func (s *RoleService) validateSystemRoleAssignment(ctx context.Context, role *dbmodel.Role, tenantID *uint) (*dbmodel.Role, *uint, error) {
	switch role.Scope {
	case dbmodel.RoleScopeSystem:
		if tenantID != nil {
			return nil, nil, errors.New("system role cannot bind tenant")
		}
		return role, nil, nil
	case dbmodel.RoleScopeTenant:
		targetTenantID := tenantID
		if targetTenantID == nil {
			if role.TenantID == nil {
				return nil, nil, errors.New("tenantId is required")
			}
			targetTenantID = role.TenantID
		}
		if _, err := s.TenantDAL.GetByID(ctx, *targetTenantID); err != nil {
			return nil, nil, err
		}
		if role.TenantID != nil && *role.TenantID != *targetTenantID {
			return nil, nil, errors.New("role tenant scope mismatch")
		}
		return role, targetTenantID, nil
	default:
		return nil, nil, errors.New("invalid role scope")
	}
}

func normalizeRoleCode(code string) string {
	return strings.ToLower(strings.TrimSpace(code))
}

func normalizeRoleScope(scope string, tenantID *uint) string {
	scope = strings.ToLower(strings.TrimSpace(scope))
	if scope == dbmodel.RoleScopeSystem {
		return dbmodel.RoleScopeSystem
	}
	if scope == dbmodel.RoleScopeTenant || tenantID != nil {
		return dbmodel.RoleScopeTenant
	}
	return dbmodel.RoleScopeSystem
}

var tenantAdminPermissionSet = map[string]bool{
	"auth-logout":              true,
	"auth-me":                  true,
	"auth-password-change":     true,
	"auth-token-create":        true,
	"dashboard-summary-get":    true,
	"bills-list":               true,
	"customers-list":           true,
	"reconciliations-list":     true,
	"users-create":             true,
	"users-list":               true,
	"users-get":                true,
	"users-update":             true,
	"users-delete":             true,
	"users-password-reset":     true,
	"permissions-list":         true,
	"roles-create":             true,
	"roles-list":               true,
	"roles-get":                true,
	"roles-update":             true,
	"roles-delete":             true,
	"roles-permissions-update": true,
}

var userPermissionSet = map[string]bool{
	"auth-logout":           true,
	"auth-me":               true,
	"auth-password-change":  true,
	"auth-token-create":     true,
	"dashboard-summary-get": true,
	"bills-list":            true,
	"customers-list":        true,
	"reconciliations-list":  true,
}
