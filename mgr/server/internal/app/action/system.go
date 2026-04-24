package action

import (
	"net/http"

	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type SystemAction struct {
	SystemService *bll.SystemService
}

func NewSystemAction(systemService *bll.SystemService) *SystemAction {
	return &SystemAction{
		SystemService: systemService,
	}
}

func (a *SystemAction) Health(c *gin.Context) {
	data, err := a.SystemService.Health(c.Request.Context())
	if err != nil {
		httpx.WriteError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.WriteData(c.Writer, http.StatusOK, data)
}
