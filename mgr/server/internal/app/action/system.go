package action

import (
	"net/http"

	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/pkg/httpx"
	"github.com/gin-gonic/gin"
)

var _ = vo.ResponseResult{}

type SystemAction struct {
	SystemService *bll.SystemService
}

func NewSystemAction(systemService *bll.SystemService) *SystemAction {
	return &SystemAction{
		SystemService: systemService,
	}
}

// Health 健康检查
// @Summary  健康检查
// @Description  返回服务健康状态与版本信息
// @Produce  json
// @Success  200  {object}  vo.ResponseResult  "执行成功"
// @Router   /api/v1/health [get]
// @ID       health-get
// @Tags     system
func (a *SystemAction) Health(c *gin.Context) {
	data, err := a.SystemService.Health(c.Request.Context())
	if err != nil {
		httpx.WriteError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.WriteData(c.Writer, http.StatusOK, data)
}
