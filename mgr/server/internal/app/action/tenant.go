package action

import (
	"errors"
	"strconv"
	"strings"

	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/internal/app/dal/dbmodel"
	"clearbill/mgr/server/internal/app/ginx"
	"clearbill/mgr/server/internal/app/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TenantAction struct {
	TenantService *bll.TenantService
}

func NewTenantAction(tenantService *bll.TenantService, _ ...*bll.RoleService) *TenantAction {
	return &TenantAction{
		TenantService: tenantService,
	}
}

// CreateTenant 创建租户
// @Summary  创建租户
// @Description  创建一条新的租户记录
// @Accept   json
// @Produce  json
// @Param    body  body      vo.CreateTenantReq  true  "body参数"
// @Success  200   {object}  vo.ResponseResult   "执行成功"
// @Failure  400   {object}  vo.ResponseResult   "参数错误"
// @Router   /api/v1/tenants [post]
// @ID       tenants-create
// @Tags     tenant
func (a *TenantAction) CreateTenant(c *gin.Context) {
	if !requireSysadmin(c) {
		return
	}

	var req vo.CreateTenantReq
	if err := ginx.ParseJSON(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.TenantService.CreateTenant(c.Request.Context(), &req)
	if err != nil {
		writeTenantError(c, err)
		return
	}

	ginx.ResSuccess(c, data)
}

// ListTenants 查询租户列表
// @Summary  查询租户列表
// @Description  根据关键字查询租户列表
// @Produce  json
// @Param    keyword  query     string            false  "关键字"
// @Success  200      {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/tenants [get]
// @ID       tenants-list
// @Tags     tenant
func (a *TenantAction) ListTenants(c *gin.Context) {
	if !requireSysadmin(c) {
		return
	}

	var req vo.ListTenantReq
	if err := ginx.ParseQuery(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.TenantService.ListTenants(c.Request.Context(), &req)
	if err != nil {
		writeTenantError(c, err)
		return
	}

	ginx.ResSuccess(c, data)
}

// GetTenant 查询租户详情
// @Summary  查询租户详情
// @Description  根据租户ID查询租户详情
// @Produce  json
// @Param    id   path      int               true  "租户ID"
// @Success  200  {object}  vo.ResponseResult "执行成功"
// @Failure  404  {object}  vo.ResponseResult "租户不存在"
// @Router   /api/v1/tenants/{id} [get]
// @ID       tenants-get
// @Tags     tenant
func (a *TenantAction) GetTenant(c *gin.Context) {
	if !requireSysadmin(c) {
		return
	}

	id, err := parseTenantID(c.Param("id"))
	if err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.TenantService.GetTenant(c.Request.Context(), id)
	if err != nil {
		writeTenantError(c, err)
		return
	}

	ginx.ResSuccess(c, data)
}

// UpdateTenant 更新租户
// @Summary  更新租户
// @Description  根据租户ID更新租户信息
// @Accept   json
// @Produce  json
// @Param    id    path      int                true  "租户ID"
// @Param    body  body      vo.UpdateTenantReq true  "body参数"
// @Success  200   {object}  vo.ResponseResult  "执行成功"
// @Failure  400   {object}  vo.ResponseResult  "参数错误"
// @Failure  404   {object}  vo.ResponseResult  "租户不存在"
// @Router   /api/v1/tenants/{id} [put]
// @ID       tenants-update
// @Tags     tenant
func (a *TenantAction) UpdateTenant(c *gin.Context) {
	if !requireSysadmin(c) {
		return
	}

	id, err := parseTenantID(c.Param("id"))
	if err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	var req vo.UpdateTenantReq
	if err := ginx.ParseJSON(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.TenantService.UpdateTenant(c.Request.Context(), id, &req)
	if err != nil {
		writeTenantError(c, err)
		return
	}

	ginx.ResSuccess(c, data)
}

// DeleteTenant 删除租户
// @Summary  删除租户
// @Description  根据租户ID删除租户
// @Produce  json
// @Param    id   path      int               true  "租户ID"
// @Success  200  {object}  vo.ResponseResult "执行成功"
// @Failure  404  {object}  vo.ResponseResult "租户不存在"
// @Router   /api/v1/tenants/{id} [delete]
// @ID       tenants-delete
// @Tags     tenant
func (a *TenantAction) DeleteTenant(c *gin.Context) {
	if !requireSysadmin(c) {
		return
	}

	id, err := parseTenantID(c.Param("id"))
	if err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	if err := a.TenantService.DeleteTenant(c.Request.Context(), id); err != nil {
		writeTenantError(c, err)
		return
	}

	ginx.ResOK(c)
}

func requireSysadmin(c *gin.Context) bool {
	user := middleware.CurrentUser(c)
	if user == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return false
	}
	if user.Role != dbmodel.RoleSysadmin {
		ginx.ResError(c, errors.New("permission denied"), 403)
		return false
	}
	return true
}

func parseTenantID(raw string) (uint, error) {
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, errors.New("invalid tenant id")
	}

	return uint(id), nil
}

func writeTenantError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		ginx.ResError(c, errors.New("tenant not found"), 404)
	case strings.Contains(strings.ToLower(err.Error()), "duplicate"):
		ginx.ResError(c, err, 409)
	case err.Error() == "管理员账号不能包含中文":
		ginx.ResError(c, err, 400)
	default:
		ginx.ResError(c, err, 500)
	}
}
