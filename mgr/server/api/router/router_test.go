package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
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
	userDAL := dal.NewUserDAL(nil)
	sessionDAL := dal.NewSessionDAL(nil)
	apiTokenDAL := dal.NewAPITokenDAL(nil)

	authService := bll.NewAuthService(userDAL, sessionDAL, apiTokenDAL)
	systemService := bll.NewSystemService(systemDAL)
	billingService := bll.NewBillingService(billingDAL)
	tenantService := bll.NewTenantService(tenantDAL, userDAL)
	userService := bll.NewUserService(userDAL, tenantDAL)

	authAction := action.NewAuthAction(authService)
	systemAction := action.NewSystemAction(systemService)
	billingAction := action.NewBillingAction(billingService)
	tenantAction := action.NewTenantAction(tenantService)
	userAction := action.NewUserAction(userService)

	return New(
		config.Config{WebRoot: "website"},
		authService,
		authAction,
		systemAction,
		billingAction,
		tenantAction,
		userAction,
	)
}

func loginSession(t *testing.T, handler http.Handler, username, password string) string {
	t.Helper()
	body := `{"username":"` + username + `","password":"` + password + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("login failed: code=%d body=%s", resp.Code, resp.Body.String())
	}

	var payload struct {
		Success bool         `json:"success"`
		Data    vo.LoginResp `json:"data"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse login response: %v", err)
	}
	for _, cookie := range resp.Result().Cookies() {
		if cookie.Name == "clear_bill_session" {
			return cookie.Value
		}
	}
	t.Fatalf("expected session cookie")
	return ""
}

func authorizedJSONRequest(method, path, session, body string) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if session != "" {
		req.Header.Set("Cookie", "clear_bill_session="+session)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
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

func TestAuthTenantAndUserFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := gin.New()
	newTestRouter().Register(handler)

	sysadminSession := loginSession(t, handler, "sysadmin", "bill123;")

	createTenantReq := authorizedJSONRequest(
		http.MethodPost,
		"/api/v1/tenants",
		sysadminSession,
		`{"code":"tenant-a","name":"Tenant A"}`,
	)
	createTenantResp := httptest.NewRecorder()
	handler.ServeHTTP(createTenantResp, createTenantReq)
	if createTenantResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", createTenantResp.Code, createTenantResp.Body.String())
	}

	var tenantPayload struct {
		Success bool                `json:"success"`
		Data    vo.CreateTenantResp `json:"data"`
	}
	if err := json.Unmarshal(createTenantResp.Body.Bytes(), &tenantPayload); err != nil {
		t.Fatalf("failed to parse tenant response: %v", err)
	}
	if tenantPayload.Data.AdminUsername == "" {
		t.Fatalf("expected tenant admin username")
	}

	tenantAdminSession := loginSession(t, handler, tenantPayload.Data.AdminUsername, "bill123;")

	createUserReq := authorizedJSONRequest(
		http.MethodPost,
		"/api/v1/users",
		tenantAdminSession,
		`{"username":"alice","displayName":"Alice","role":"user","status":"active"}`,
	)
	createUserResp := httptest.NewRecorder()
	handler.ServeHTTP(createUserResp, createUserReq)
	if createUserResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", createUserResp.Code, createUserResp.Body.String())
	}

	var userPayload struct {
		Success bool              `json:"success"`
		Data    vo.CreateUserResp `json:"data"`
	}
	if err := json.Unmarshal(createUserResp.Body.Bytes(), &userPayload); err != nil {
		t.Fatalf("failed to parse user response: %v", err)
	}

	userSession := loginSession(t, handler, "alice", "bill123;")

	changeOwnPasswordReq := authorizedJSONRequest(
		http.MethodPut,
		"/api/v1/auth/password",
		userSession,
		`{"oldPassword":"bill123;","newPassword":"newpass123"}`,
	)
	changeOwnPasswordResp := httptest.NewRecorder()
	handler.ServeHTTP(changeOwnPasswordResp, changeOwnPasswordReq)
	if changeOwnPasswordResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", changeOwnPasswordResp.Code, changeOwnPasswordResp.Body.String())
	}

	resetPasswordReq := authorizedJSONRequest(
		http.MethodPut,
		"/api/v1/users/"+toString(userPayload.Data.User.ID)+"/password",
		tenantAdminSession,
		`{"newPassword":"reset12345"}`,
	)
	resetPasswordResp := httptest.NewRecorder()
	handler.ServeHTTP(resetPasswordResp, resetPasswordReq)
	if resetPasswordResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", resetPasswordResp.Code, resetPasswordResp.Body.String())
	}

	createTokenReq := authorizedJSONRequest(
		http.MethodPost,
		"/api/v1/auth/tokens",
		tenantAdminSession,
		`{"name":"external-client"}`,
	)
	createTokenResp := httptest.NewRecorder()
	handler.ServeHTTP(createTokenResp, createTokenReq)
	if createTokenResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", createTokenResp.Code, createTokenResp.Body.String())
	}
	if !strings.Contains(createTokenResp.Body.String(), "external-client") {
		t.Fatalf("expected api token response")
	}
}

func toString(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}
