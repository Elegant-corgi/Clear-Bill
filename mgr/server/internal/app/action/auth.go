package action

import (
	"errors"

	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/internal/app/ginx"
	"clearbill/mgr/server/internal/app/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthAction struct {
	AuthService *bll.AuthService
}

func NewAuthAction(authService *bll.AuthService) *AuthAction {
	return &AuthAction{
		AuthService: authService,
	}
}

// Login 用户登录
// @Summary  用户登录
// @Description  使用用户名和密码登录，设置后台 session
// @Accept   json
// @Produce  json
// @Param    body  body      vo.LoginReq       true  "body参数"
// @Success  200   {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/auth/login [post]
// @ID       auth-login
// @Tags     auth
func (a *AuthAction) Login(c *gin.Context) {
	var req vo.LoginReq
	if err := ginx.ParseJSON(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	user, sessionToken, err := a.AuthService.Login(c.Request.Context(), &req)
	if err != nil {
		writeAuthError(c, err)
		return
	}

	c.SetCookie("clear_bill_session", sessionToken, 86400, "/", "", false, true)
	ginx.ResSuccess(c, vo.LoginResp{User: bll.ToUserVO(user)})
}

// Logout 用户登出
// @Summary  用户登出
// @Description  清除当前后台 session
// @Produce  json
// @Success  200  {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/auth/logout [post]
// @ID       auth-logout
// @Tags     auth
func (a *AuthAction) Logout(c *gin.Context) {
	if err := a.AuthService.Logout(c.Request.Context(), middleware.CurrentToken(c)); err != nil {
		writeAuthError(c, err)
		return
	}

	c.SetCookie("clear_bill_session", "", -1, "/", "", false, true)
	ginx.ResOK(c)
}

// CurrentUser 当前用户信息
// @Summary  当前用户信息
// @Description  返回当前登录用户信息
// @Produce  json
// @Success  200  {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/auth/me [get]
// @ID       auth-me
// @Tags     auth
func (a *AuthAction) CurrentUser(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	ginx.ResSuccess(c, bll.ToUserVO(user))
}

// ChangeOwnPassword 修改自己的密码
// @Summary  修改自己的密码
// @Description  当前登录用户修改自己的密码
// @Accept   json
// @Produce  json
// @Param    body  body      vo.ChangeOwnPasswordReq  true  "body参数"
// @Success  200   {object}  vo.ResponseResult        "执行成功"
// @Router   /api/v1/auth/password [put]
// @ID       auth-password-change
// @Tags     auth
func (a *AuthAction) ChangeOwnPassword(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	var req vo.ChangeOwnPasswordReq
	if err := ginx.ParseJSON(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	if err := a.AuthService.ChangeOwnPassword(c.Request.Context(), user.ID, &req); err != nil {
		writeAuthError(c, err)
		return
	}

	ginx.ResOK(c)
}

// CreateAPIToken 创建 API token
// @Summary  创建 API token
// @Description  为当前已登录用户签发一个用于外部 API 调用的 token
// @Accept   json
// @Produce  json
// @Param    body  body      vo.CreateAPITokenReq  true  "body参数"
// @Success  200   {object}  vo.ResponseResult     "执行成功"
// @Router   /api/v1/auth/tokens [post]
// @ID       auth-token-create
// @Tags     auth
func (a *AuthAction) CreateAPIToken(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	var req vo.CreateAPITokenReq
	if err := ginx.ParseJSON(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.AuthService.CreateAPIToken(c.Request.Context(), user.ID, &req)
	if err != nil {
		writeAuthError(c, err)
		return
	}

	ginx.ResSuccess(c, data)
}

func writeAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		ginx.ResError(c, errors.New("invalid username or password"), 401)
	case err != nil && (err.Error() == "invalid username or password" || err.Error() == "user is disabled" || err.Error() == "old password is incorrect"):
		ginx.ResError(c, err, 400)
	default:
		ginx.ResError(c, err, 500)
	}
}
