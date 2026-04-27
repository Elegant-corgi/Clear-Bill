package action

import (
	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/internal/app/ginx"
	"github.com/gin-gonic/gin"
)

var _ = vo.ResponseResult{}

type BillingAction struct {
	BillingService *bll.BillingService
}

func NewBillingAction(billingService *bll.BillingService) *BillingAction {
	return &BillingAction{
		BillingService: billingService,
	}
}

// DashboardSummary 查询仪表盘概览
// @Summary  查询仪表盘概览
// @Description  返回账单工作台概览数据
// @Produce  json
// @Success  200  {object}  vo.ResponseResult  "执行成功"
// @Router   /api/v1/dashboard/summary [get]
// @ID       dashboard-summary-get
// @Tags     billing
func (a *BillingAction) DashboardSummary(c *gin.Context) {
	data, err := a.BillingService.GetDashboardSummary(c.Request.Context())
	if err != nil {
		ginx.ResError(c, err, 500)
		return
	}

	ginx.ResSuccess(c, data)
}

// ListBills 查询账单列表
// @Summary  查询账单列表
// @Description  返回账单列表数据
// @Produce  json
// @Success  200  {object}  vo.ResponseResult  "执行成功"
// @Router   /api/v1/bills [get]
// @ID       bills-list
// @Tags     billing
func (a *BillingAction) ListBills(c *gin.Context) {
	data, err := a.BillingService.ListBills(c.Request.Context())
	if err != nil {
		ginx.ResError(c, err, 500)
		return
	}

	ginx.ResSuccess(c, data)
}

// ListCustomers 查询客户列表
// @Summary  查询客户列表
// @Description  返回客户列表数据
// @Produce  json
// @Success  200  {object}  vo.ResponseResult  "执行成功"
// @Router   /api/v1/customers [get]
// @ID       customers-list
// @Tags     billing
func (a *BillingAction) ListCustomers(c *gin.Context) {
	data, err := a.BillingService.ListCustomers(c.Request.Context())
	if err != nil {
		ginx.ResError(c, err, 500)
		return
	}

	ginx.ResSuccess(c, data)
}

// ListReconciliationTasks 查询对账任务列表
// @Summary  查询对账任务列表
// @Description  返回对账任务列表数据
// @Produce  json
// @Success  200  {object}  vo.ResponseResult  "执行成功"
// @Router   /api/v1/reconciliations [get]
// @ID       reconciliations-list
// @Tags     billing
func (a *BillingAction) ListReconciliationTasks(c *gin.Context) {
	data, err := a.BillingService.ListReconciliationTasks(c.Request.Context())
	if err != nil {
		ginx.ResError(c, err, 500)
		return
	}

	ginx.ResSuccess(c, data)
}
