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

type CredentialAction struct {
	CredentialService *bll.CredentialService
}

func NewCredentialAction(credentialService *bll.CredentialService) *CredentialAction {
	return &CredentialAction{CredentialService: credentialService}
}

// CreateCredential 创建凭证
// @Summary  创建凭证
// @Description  为当前登录用户创建 AK/SK 或 token 凭证
// @Accept   json
// @Produce  json
// @Param    body  body      vo.CreateCredentialReq  true  "body参数"
// @Success  200   {object}  vo.ResponseResult       "执行成功"
// @Router   /api/v1/credentials [post]
// @ID       credentials-create
// @Tags     credential
func (a *CredentialAction) CreateCredential(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	var req vo.CreateCredentialReq
	if err := ginx.ParseJSON(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.CredentialService.CreateCredential(c.Request.Context(), actor, &req)
	if err != nil {
		writeCredentialError(c, err)
		return
	}

	ginx.ResSuccess(c, data)
}

// ListCredentials 查询凭证列表
// @Summary  查询凭证列表
// @Description  查询当前用户可见的凭证，系统管理员可查看全部
// @Produce  json
// @Param    type    query     string            false  "凭证类型"
// @Param    status  query     string            false  "凭证状态"
// @Success  200     {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/credentials [get]
// @ID       credentials-list
// @Tags     credential
func (a *CredentialAction) ListCredentials(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	var req vo.ListCredentialReq
	if err := ginx.ParseQuery(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.CredentialService.ListCredentials(c.Request.Context(), actor, &req)
	if err != nil {
		writeCredentialError(c, err)
		return
	}

	ginx.ResSuccess(c, data)
}

// GetCredential 查询凭证详情
// @Summary  查询凭证详情
// @Description  查询指定凭证的元数据，不返回密钥明文
// @Produce  json
// @Param    id   path      int               true  "凭证ID"
// @Success  200  {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/credentials/{id} [get]
// @ID       credentials-get
// @Tags     credential
func (a *CredentialAction) GetCredential(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	id, err := parseCredentialID(c.Param("id"))
	if err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.CredentialService.GetCredential(c.Request.Context(), actor, id)
	if err != nil {
		writeCredentialError(c, err)
		return
	}

	ginx.ResSuccess(c, data)
}

// RotateCredential 轮转凭证
// @Summary  轮转凭证
// @Description  将旧凭证置为已轮转并生成一条新的凭证记录
// @Accept   json
// @Produce  json
// @Param    id    path      int                 true  "凭证ID"
// @Param    body  body      vo.RotateCredentialReq true  "body参数"
// @Success  200   {object}  vo.ResponseResult    "执行成功"
// @Router   /api/v1/credentials/{id}/rotate [post]
// @ID       credentials-rotate
// @Tags     credential
func (a *CredentialAction) RotateCredential(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	id, err := parseCredentialID(c.Param("id"))
	if err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	var req vo.RotateCredentialReq
	if err := ginx.ParseJSON(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.CredentialService.RotateCredential(c.Request.Context(), actor, id, &req)
	if err != nil {
		writeCredentialError(c, err)
		return
	}

	ginx.ResSuccess(c, data)
}

// DeleteCredential 删除凭证
// @Summary  删除凭证
// @Description  删除指定凭证
// @Produce  json
// @Param    id   path      int               true  "凭证ID"
// @Success  200  {object}  vo.ResponseResult "执行成功"
// @Router   /api/v1/credentials/{id} [delete]
// @ID       credentials-delete
// @Tags     credential
func (a *CredentialAction) DeleteCredential(c *gin.Context) {
	actor := middleware.CurrentUser(c)
	if actor == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	id, err := parseCredentialID(c.Param("id"))
	if err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	if err := a.CredentialService.DeleteCredential(c.Request.Context(), actor, id); err != nil {
		writeCredentialError(c, err)
		return
	}

	ginx.ResOK(c)
}

func parseCredentialID(raw string) (uint, error) {
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, errors.New("invalid credential id")
	}
	return uint(id), nil
}

func writeCredentialError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		ginx.ResError(c, errors.New("credential not found"), 404)
	case strings.Contains(strings.ToLower(err.Error()), "duplicate"):
		ginx.ResError(c, err, 409)
	case err.Error() == "permission denied":
		ginx.ResError(c, err, 403)
	default:
		ginx.ResError(c, err, 400)
	}
}
