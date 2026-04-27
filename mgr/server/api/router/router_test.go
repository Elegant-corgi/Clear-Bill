package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"clearbill/mgr/server/api/vo"
	"clearbill/mgr/server/internal/app/action"
	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/internal/app/config"
	"clearbill/mgr/server/internal/app/dal"
	"github.com/gin-gonic/gin"
)

func newTestRouter() *Router {
	systemDAL := dal.NewSystemDAL()
	billingDAL := dal.NewBillingDAL()
	tenantDAL := dal.NewTenantDAL(nil)

	systemService := bll.NewSystemService(systemDAL)
	billingService := bll.NewBillingService(billingDAL)
	tenantService := bll.NewTenantService(tenantDAL)

	systemAction := action.NewSystemAction(systemService)
	billingAction := action.NewBillingAction(billingService)
	tenantAction := action.NewTenantAction(tenantService)

	return New(config.Config{WebRoot: "website"}, systemAction, billingAction, tenantAction)
}

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := gin.New()
	newTestRouter().Register(handler)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var payload struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if !payload.Success {
		t.Fatalf("expected success response")
	}
}

func TestBillsEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := gin.New()
	newTestRouter().Register(handler)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/bills", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var payload struct {
		Success bool  `json:"success"`
		Data    []any `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(payload.Data) == 0 {
		t.Fatalf("expected bill data")
	}
}

func TestTenantCRUD(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := gin.New()
	newTestRouter().Register(handler)

	createBody := `{"code":"tenant-a","name":"Tenant A","contactName":"Alex","contactPhone":"13800000000","status":"active","remark":"first tenant"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	handler.ServeHTTP(createResp, createReq)

	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", createResp.Code, createResp.Body.String())
	}

	var created struct {
		Success bool      `json:"success"`
		Data    vo.Tenant `json:"data"`
	}
	if err := json.Unmarshal(createResp.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to parse create response: %v", err)
	}
	if created.Data.ID == 0 {
		t.Fatalf("expected created tenant id")
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/tenants?keyword=tenant", nil)
	listResp := httptest.NewRecorder()
	handler.ServeHTTP(listResp, listReq)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", listResp.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/tenants/1", nil)
	getResp := httptest.NewRecorder()
	handler.ServeHTTP(getResp, getReq)
	if getResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", getResp.Code)
	}

	updateBody := `{"code":"tenant-a","name":"Tenant A Updated","contactName":"Alex","contactPhone":"13900000000","status":"inactive","remark":"updated"}`
	updateReq := httptest.NewRequest(http.MethodPut, "/api/v1/tenants/1", bytes.NewBufferString(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp := httptest.NewRecorder()
	handler.ServeHTTP(updateResp, updateReq)
	if updateResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", updateResp.Code, updateResp.Body.String())
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/tenants/1", nil)
	deleteResp := httptest.NewRecorder()
	handler.ServeHTTP(deleteResp, deleteReq)
	if deleteResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", deleteResp.Code)
	}
}
