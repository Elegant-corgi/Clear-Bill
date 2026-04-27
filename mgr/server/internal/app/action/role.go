package action

import (
	"errors"
	"strconv"
	"strings"

	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/internal/app/ginx"
	"clearbill/mgr/server/internal/app/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleAction struct {
	RoleService *bll.RoleService
}

func NewRoleAction(roleService *bll.RoleService) *RoleAction {
	return &RoleAction{
		RoleService: roleService,
	}
}

// ListPermissions 查询权限列表
// @Summary  查询权限列表
// @Description  返回基于 Swagger 接口注解生成的权限列表
// @Produce  json
// @Success  200  {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/permissions [get]
// @ID       permissions-list
// @Tags     role
func (a *RoleAction) ListPermissions(c *gin.Context) {
	ginx.ResSuccess(c, a.RoleService.ListPermissions())
}

// CreateRole 创建角色
// @Summary  创建角色
// @Description  创建角色并分配权限
// @Accept   json
// @Produce  json
// @Param    body  body      vo.CreateRoleReq  true  "body参数"
// @Success  200   {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/roles [post]
// @ID       roles-create
// @Tags     role
func (a *RoleAction) CreateRole(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	var req vo.CreateRoleReq
	if err := ginx.ParseJSON(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.RoleService.CreateRole(c.Request.Context(), actor, &req)
	if err != nil {
		writeRoleError(c, err)
		return
	}
	ginx.ResSuccess(c, data)
}

// ListRoles 查询角色列表
// @Summary  查询角色列表
// @Description  查询可见范围内的角色列表
// @Produce  json
// @Param    keyword   query     string            false  "关键字"
// @Param    scope     query     string            false  "角色作用域"
// @Param    tenantId  query     int               false  "租户ID"
// @Success  200       {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/roles [get]
// @ID       roles-list
// @Tags     role
func (a *RoleAction) ListRoles(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	var req vo.ListRoleReq
	if err := ginx.ParseQuery(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.RoleService.ListRoles(c.Request.Context(), actor, &req)
	if err != nil {
		writeRoleError(c, err)
		return
	}
	ginx.ResSuccess(c, data)
}

// GetRole 查询角色详情
// @Summary  查询角色详情
// @Description  查询指定角色信息
// @Produce  json
// @Param    id   path      int               true  "角色ID"
// @Success  200  {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/roles/{id} [get]
// @ID       roles-get
// @Tags     role
func (a *RoleAction) GetRole(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	id, err := parseRoleID(c.Param("id"))
	if err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.RoleService.GetRole(c.Request.Context(), actor, id)
	if err != nil {
		writeRoleError(c, err)
		return
	}
	ginx.ResSuccess(c, data)
}

// UpdateRole 更新角色
// @Summary  更新角色
// @Description  更新指定角色的基础信息
// @Accept   json
// @Produce  json
// @Param    id    path      int               true  "角色ID"
// @Param    body  body      vo.UpdateRoleReq  true  "body参数"
// @Success  200   {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/roles/{id} [put]
// @ID       roles-update
// @Tags     role
func (a *RoleAction) UpdateRole(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	id, err := parseRoleID(c.Param("id"))
	if err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	var req vo.UpdateRoleReq
	if err := ginx.ParseJSON(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.RoleService.UpdateRole(c.Request.Context(), actor, id, &req)
	if err != nil {
		writeRoleError(c, err)
		return
	}
	ginx.ResSuccess(c, data)
}

// UpdateRolePermissions 更新角色权限
// @Summary  更新角色权限
// @Description  覆盖指定角色的权限列表
// @Accept   json
// @Produce  json
// @Param    id    path      int                            true  "角色ID"
// @Param    body  body      vo.UpdateRolePermissionsReq    true  "body参数"
// @Success  200   {object}  vo.ResponseResult              "执行成功"
// @Router   /api/v1/roles/{id}/permissions [put]
// @ID       roles-permissions-update
// @Tags     role
func (a *RoleAction) UpdateRolePermissions(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	id, err := parseRoleID(c.Param("id"))
	if err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	var req vo.UpdateRolePermissionsReq
	if err := ginx.ParseJSON(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.RoleService.UpdateRolePermissions(c.Request.Context(), actor, id, &req)
	if err != nil {
		writeRoleError(c, err)
		return
	}
	ginx.ResSuccess(c, data)
}

// DeleteRole 删除角色
// @Summary  删除角色
// @Description  删除指定角色
// @Produce  json
// @Param    id   path      int               true  "角色ID"
// @Success  200  {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/roles/{id} [delete]
// @ID       roles-delete
// @Tags     role
func (a *RoleAction) DeleteRole(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	id, err := parseRoleID(c.Param("id"))
	if err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	if err := a.RoleService.DeleteRole(c.Request.Context(), actor, id); err != nil {
		writeRoleError(c, err)
		return
	}
	ginx.ResOK(c)
}

func parseRoleID(raw string) (uint, error) {
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, errors.New("invalid role id")
	}
	return uint(id), nil
}

func writeRoleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		ginx.ResError(c, errors.New("role not found"), 404)
	case strings.Contains(strings.ToLower(err.Error()), "duplicate"):
		ginx.ResError(c, err, 409)
	case err.Error() == "permission denied":
		ginx.ResError(c, err, 403)
	case err.Error() == "builtin role cannot be modified" || err.Error() == "role is in use":
		ginx.ResError(c, err, 400)
	default:
		ginx.ResError(c, err, 400)
	}
}
