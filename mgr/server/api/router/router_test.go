package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"clearbill/mgr/server/internal/app/action"
	"clearbill/mgr/server/internal/app/bll"
	"clearbill/mgr/server/internal/app/config"
	"clearbill/mgr/server/internal/app/dal"
	"github.com/gin-gonic/gin"
)

func newTestRouter() *Router {
	systemDAL := dal.NewSystemDAL()
	billingDAL := dal.NewBillingDAL()

	systemService := bll.NewSystemService(systemDAL)
	billingService := bll.NewBillingService(billingDAL)

	systemAction := action.NewSystemAction(systemService)
	billingAction := action.NewBillingAction(billingService)

	return New(config.Config{WebRoot: "website"}, systemAction, billingAction)
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
