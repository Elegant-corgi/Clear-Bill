package action

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/pkg/httpx"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TenantAction struct {
	TenantService *bll.TenantService
}

func NewTenantAction(tenantService *bll.TenantService) *TenantAction {
	return &TenantAction{
		TenantService: tenantService,
	}
}

func (a *TenantAction) CreateTenant(c *gin.Context) {
	var req vo.CreateTenantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	data, err := a.TenantService.CreateTenant(c.Request.Context(), &req)
	if err != nil {
		writeTenantError(c, err)
		return
	}

	httpx.WriteData(c.Writer, http.StatusCreated, data)
}

func (a *TenantAction) ListTenants(c *gin.Context) {
	var req vo.ListTenantReq
	if err := c.ShouldBindQuery(&req); err != nil {
		httpx.WriteError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	data, err := a.TenantService.ListTenants(c.Request.Context(), &req)
	if err != nil {
		writeTenantError(c, err)
		return
	}

	httpx.WriteData(c.Writer, http.StatusOK, data)
}

func (a *TenantAction) GetTenant(c *gin.Context) {
	id, err := parseTenantID(c.Param("id"))
	if err != nil {
		httpx.WriteError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	data, err := a.TenantService.GetTenant(c.Request.Context(), id)
	if err != nil {
		writeTenantError(c, err)
		return
	}

	httpx.WriteData(c.Writer, http.StatusOK, data)
}

func (a *TenantAction) UpdateTenant(c *gin.Context) {
	id, err := parseTenantID(c.Param("id"))
	if err != nil {
		httpx.WriteError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	var req vo.UpdateTenantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	data, err := a.TenantService.UpdateTenant(c.Request.Context(), id, &req)
	if err != nil {
		writeTenantError(c, err)
		return
	}

	httpx.WriteData(c.Writer, http.StatusOK, data)
}

func (a *TenantAction) DeleteTenant(c *gin.Context) {
	id, err := parseTenantID(c.Param("id"))
	if err != nil {
		httpx.WriteError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.TenantService.DeleteTenant(c.Request.Context(), id); err != nil {
		writeTenantError(c, err)
		return
	}

	httpx.WriteData(c.Writer, http.StatusOK, gin.H{"deleted": true})
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
		httpx.WriteError(c.Writer, http.StatusNotFound, "tenant not found")
	case strings.Contains(strings.ToLower(err.Error()), "duplicate"):
		httpx.WriteError(c.Writer, http.StatusConflict, err.Error())
	default:
		httpx.WriteError(c.Writer, http.StatusInternalServerError, err.Error())
	}
}
