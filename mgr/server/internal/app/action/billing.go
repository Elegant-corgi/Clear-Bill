package action

import (
	"net/http"

	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/pkg/httpx"
	"github.com/gin-gonic/gin"
)

type BillingAction struct {
	BillingService *bll.BillingService
}

func NewBillingAction(billingService *bll.BillingService) *BillingAction {
	return &BillingAction{
		BillingService: billingService,
	}
}

func (a *BillingAction) DashboardSummary(c *gin.Context) {
	data, err := a.BillingService.GetDashboardSummary(c.Request.Context())
	if err != nil {
		httpx.WriteError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.WriteData(c.Writer, http.StatusOK, data)
}

func (a *BillingAction) ListBills(c *gin.Context) {
	data, err := a.BillingService.ListBills(c.Request.Context())
	if err != nil {
		httpx.WriteError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.WriteData(c.Writer, http.StatusOK, data)
}

func (a *BillingAction) ListCustomers(c *gin.Context) {
	data, err := a.BillingService.ListCustomers(c.Request.Context())
	if err != nil {
		httpx.WriteError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.WriteData(c.Writer, http.StatusOK, data)
}

func (a *BillingAction) ListReconciliationTasks(c *gin.Context) {
	data, err := a.BillingService.ListReconciliationTasks(c.Request.Context())
	if err != nil {
		httpx.WriteError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.WriteData(c.Writer, http.StatusOK, data)
}
