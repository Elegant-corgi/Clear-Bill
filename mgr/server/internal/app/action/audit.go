package action

import (
	"errors"

	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/internal/app/ginx"
	"clearbill/mgr/server/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

type AuditAction struct {
	AuditService *bll.AuditService
}

func NewAuditAction(auditService *bll.AuditService) *AuditAction {
	return &AuditAction{AuditService: auditService}
}

// ListAuditLogs 查询审计日志
// @Summary  查询审计日志
// @Description  支持按用户、操作类型和时间段筛选审计日志
// @Produce  json
// @Param    user       query     string  false  "用户"
// @Param    operation  query     string  false  "操作类型"
// @Param    startTime  query     string  false  "开始时间"
// @Param    endTime    query     string  false  "结束时间"
// @Param    page       query     int     false  "页码"
// @Param    pageSize   query     int     false  "每页数量"
// @Success  200        {object}  vo.ResponseResult
// @Router   /api/v1/audit-logs [get]
// @ID       audit-logs-list
// @Tags     audit
func (a *AuditAction) ListAuditLogs(c *gin.Context) {
	user := middleware.CurrentUser(c)
	if user == nil {
		ginx.ResError(c, errors.New("unauthorized"), 401)
		return
	}

	var req vo.ListAuditLogReq
	if err := ginx.ParseQuery(c, &req); err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	data, err := a.AuditService.ListAuditLogs(c.Request.Context(), &req)
	if err != nil {
		ginx.ResError(c, err, 400)
		return
	}

	ginx.ResSuccess(c, data)
}
