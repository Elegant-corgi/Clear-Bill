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
	roleDAL := dal.NewRoleDAL(nil)
	rolePermissionDAL := dal.NewRolePermissionDAL(nil)
	userDAL := dal.NewUserDAL(nil)
	sessionDAL := dal.NewSessionDAL(nil)
	apiTokenDAL := dal.NewAPITokenDAL(nil)
	permissionCatalog := bll.NewPermissionCatalog()

	authService := bll.NewAuthService(userDAL, sessionDAL, apiTokenDAL)
	systemService := bll.NewSystemService(systemDAL)
	billingService := bll.NewBillingService(billingDAL)
	tenantService := bll.NewTenantService(tenantDAL, userDAL)
	roleService := bll.NewRoleService(roleDAL, rolePermissionDAL, userDAL, tenantDAL, permissionCatalog)
	userService := bll.NewUserService(userDAL, tenantDAL, roleService)

	authAction := action.NewAuthAction(authService)
	systemAction := action.NewSystemAction(systemService)
	billingAction := action.NewBillingAction(billingService)
	tenantAction := action.NewTenantAction(tenantService, roleService)
	roleAction := action.NewRoleAction(roleService)
	userAction := action.NewUserAction(userService)

	return New(
		config.Config{WebRoot: "website"},
		authService,
		roleService,
		authAction,
		systemAction,
		billingAction,
		tenantAction,
		roleAction,
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
	sysadminSession := loginSession(t, handler, "sysadmin", "stor123;")
	request := authorizedJSONRequest(http.MethodGet, "/api/v1/bills", sysadminSession, "")
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

	sysadminSession := loginSession(t, handler, "sysadmin", "stor123;")

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

	tenantAdminSession := loginSession(t, handler, tenantPayload.Data.AdminUsername, "stor123;")

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

	userSession := loginSession(t, handler, "alice", "stor123;")

	changeOwnPasswordReq := authorizedJSONRequest(
		http.MethodPut,
		"/api/v1/auth/password",
		userSession,
		`{"oldPassword":"stor123;","newPassword":"newpass123"}`,
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

func TestRBACRolePermissionFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := gin.New()
	newTestRouter().Register(handler)

	sysadminSession := loginSession(t, handler, "sysadmin", "stor123;")

	createTenantReq := authorizedJSONRequest(
		http.MethodPost,
		"/api/v1/tenants",
		sysadminSession,
		`{"code":"tenant-b","name":"Tenant B"}`,
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

	tenantAdminSession := loginSession(t, handler, tenantPayload.Data.AdminUsername, "stor123;")

	listPermissionsReq := authorizedJSONRequest(http.MethodGet, "/api/v1/permissions", tenantAdminSession, "")
	listPermissionsResp := httptest.NewRecorder()
	handler.ServeHTTP(listPermissionsResp, listPermissionsReq)
	if listPermissionsResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", listPermissionsResp.Code, listPermissionsResp.Body.String())
	}
	if !strings.Contains(listPermissionsResp.Body.String(), "users-list") {
		t.Fatalf("expected swagger-derived permissions in response")
	}

	createRoleReq := authorizedJSONRequest(
		http.MethodPost,
		"/api/v1/roles",
		tenantAdminSession,
		`{"code":"auditor","name":"Auditor","permissionIds":["users-list"]}`,
	)
	createRoleResp := httptest.NewRecorder()
	handler.ServeHTTP(createRoleResp, createRoleReq)
	if createRoleResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", createRoleResp.Code, createRoleResp.Body.String())
	}

	var rolePayload struct {
		Success bool    `json:"success"`
		Data    vo.Role `json:"data"`
	}
	if err := json.Unmarshal(createRoleResp.Body.Bytes(), &rolePayload); err != nil {
		t.Fatalf("failed to parse role response: %v", err)
	}
	if rolePayload.Data.Code != "auditor" {
		t.Fatalf("expected role code auditor")
	}

	createUserReq := authorizedJSONRequest(
		http.MethodPost,
		"/api/v1/users",
		tenantAdminSession,
		`{"username":"bob","displayName":"Bob","role":"auditor","status":"active"}`,
	)
	createUserResp := httptest.NewRecorder()
	handler.ServeHTTP(createUserResp, createUserReq)
	if createUserResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", createUserResp.Code, createUserResp.Body.String())
	}

	bobSession := loginSession(t, handler, "bob", "stor123;")

	listUsersReq := authorizedJSONRequest(http.MethodGet, "/api/v1/users", bobSession, "")
	listUsersResp := httptest.NewRecorder()
	handler.ServeHTTP(listUsersResp, listUsersReq)
	if listUsersResp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", listUsersResp.Code, listUsersResp.Body.String())
	}

	createUserDeniedReq := authorizedJSONRequest(
		http.MethodPost,
		"/api/v1/users",
		bobSession,
		`{"username":"eve","displayName":"Eve","role":"user","status":"active"}`,
	)
	createUserDeniedResp := httptest.NewRecorder()
	handler.ServeHTTP(createUserDeniedResp, createUserDeniedReq)
	if createUserDeniedResp.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d body=%s", createUserDeniedResp.Code, createUserDeniedResp.Body.String())
	}
}

func toString(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}
