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

type UserAction struct {
	UserService *bll.UserService
}

func NewUserAction(userService *bll.UserService) *UserAction {
	return &UserAction{
		UserService: userService,
	}
}

// CreateUser 创建用户
// @Summary  创建用户
// @Description  仅超管和租户管理员可创建用户，默认密码为 stor123;
// @Accept   json
// @Produce  json
// @Param    body  body      vo.CreateUserReq   true  "body参数"
// @Success  200   {object}  vo.ResponseResult  "执行成功"
// @Router   /api/v1/users [post]
// @ID       users-create
// @Tags     user
func (a *UserAction) CreateUser(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	var req vo.CreateUserReq
	if err := ginx.ParseJSON(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.UserService.CreateUser(c.Request.Context(), actor, &req)
	if err != nil {
		writeUserError(c, err)
		return
	}

	ginx.ResSuccess(c, data)
}

// ListUsers 查询用户列表
// @Summary  查询用户列表
// @Description  超管可查询所有用户，租户管理员仅可查询本租户用户
// @Produce  json
// @Param    keyword   query     string            false  "关键字"
// @Param    tenantId  query     int               false  "租户ID"
// @Success  200       {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/users [get]
// @ID       users-list
// @Tags     user
func (a *UserAction) ListUsers(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	var req vo.ListUserReq
	if err := ginx.ParseQuery(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.UserService.ListUsers(c.Request.Context(), actor, &req)
	if err != nil {
		writeUserError(c, err)
		return
	}

	ginx.ResSuccess(c, data)
}

// GetUser 查询用户详情
// @Summary  查询用户详情
// @Description  查询指定用户详情
// @Produce  json
// @Param    id   path      int               true  "用户ID"
// @Success  200  {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/users/{id} [get]
// @ID       users-get
// @Tags     user
func (a *UserAction) GetUser(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	id, err := parseUserID(c.Param("id"))
	if err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.UserService.GetUser(c.Request.Context(), actor, id)
	if err != nil {
		writeUserError(c, err)
		return
	}

	ginx.ResSuccess(c, data)
}

// UpdateUser 更新用户
// @Summary  更新用户
// @Description  更新指定用户信息
// @Accept   json
// @Produce  json
// @Param    id    path      int               true  "用户ID"
// @Param    body  body      vo.UpdateUserReq  true  "body参数"
// @Success  200   {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/users/{id} [put]
// @ID       users-update
// @Tags     user
func (a *UserAction) UpdateUser(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	id, err := parseUserID(c.Param("id"))
	if err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	var req vo.UpdateUserReq
	if err := ginx.ParseJSON(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.UserService.UpdateUser(c.Request.Context(), actor, id, &req)
	if err != nil {
		writeUserError(c, err)
		return
	}

	ginx.ResSuccess(c, data)
}

// DeleteUser 删除用户
// @Summary  删除用户
// @Description  删除指定用户
// @Produce  json
// @Param    id   path      int               true  "用户ID"
// @Success  200  {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/users/{id} [delete]
// @ID       users-delete
// @Tags     user
func (a *UserAction) DeleteUser(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	id, err := parseUserID(c.Param("id"))
	if err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	if err := a.UserService.DeleteUser(c.Request.Context(), actor, id); err != nil {
		writeUserError(c, err)
		return
	}

	ginx.ResOK(c)
}

// ResetPassword 重置用户密码
// @Summary  重置用户密码
// @Description  超管可修改所有用户密码，租户管理员可修改本租户所有用户密码
// @Accept   json
// @Produce  json
// @Param    id    path      int                 true  "用户ID"
// @Param    body  body      vo.ResetPasswordReq true  "body参数"
// @Success  200   {object}  vo.ResponseResult   "执行成功"
// @Router   /api/v1/users/{id}/password [put]
// @ID       users-password-reset
// @Tags     user
func (a *UserAction) ResetPassword(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	id, err := parseUserID(c.Param("id"))
	if err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	var req vo.ResetPasswordReq
	if err := ginx.ParseJSON(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	if err := a.UserService.ResetPassword(c.Request.Context(), actor, id, &req); err != nil {
		writeUserError(c, err)
		return
	}

	ginx.ResOK(c)
}

func parseUserID(raw string) (uint, error) {
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, errors.New("invalid user id")
	}
	return uint(id), nil
}

func writeUserError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		ginx.ResError(c, errors.New("user not found"), 404)
	case strings.Contains(strings.ToLower(err.Error()), "duplicate"):
		ginx.ResError(c, err, 409)
	case err.Error() == "permission denied" || err.Error() == "cannot delete current user" || err.Error() == "tenantId is required" || err.Error() == "tenant admin missing tenant scope":
		ginx.ResError(c, err, 403)
	default:
		ginx.ResError(c, err, 400)
	}
}
