package ginx

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var duplicateEntryPattern = regexp.MustCompile(`Duplicate entry '(.+)' for key '(.+)'`)

var fieldLabelByKey = map[string]string{
	"idx_tenants_code":      "\u79df\u6237\u7f16\u7801",
	"tenants.code":          "\u79df\u6237\u7f16\u7801",
	"idx_roles_code":        "\u89d2\u8272\u7f16\u7801",
	"roles.code":            "\u89d2\u8272\u7f16\u7801",
	"idx_users_username":    "\u767b\u5f55\u8d26\u53f7",
	"users.username":        "\u767b\u5f55\u8d26\u53f7",
	"uk_role_permission":    "\u89d2\u8272\u6743\u9650",
	"user_api_tokens.token": "API Token",
	"user_sessions.token":   "\u4f1a\u8bdd\u6807\u8bc6",
}

var errorMessageMap = map[string]string{
	"unauthorized":                      "\u672a\u767b\u5f55\u6216\u767b\u5f55\u5df2\u5931\u6548\uff0c\u8bf7\u91cd\u65b0\u767b\u5f55",
	"permission denied":                 "\u5f53\u524d\u8d26\u53f7\u6ca1\u6709\u6743\u9650\u6267\u884c\u8be5\u64cd\u4f5c",
	"permission not registered":         "\u63a5\u53e3\u6743\u9650\u672a\u6ce8\u518c\uff0c\u8bf7\u8054\u7cfb\u7ba1\u7406\u5458\u5904\u7406",
	"missing session":                   "\u767b\u5f55\u4f1a\u8bdd\u4e0d\u5b58\u5728\uff0c\u8bf7\u91cd\u65b0\u767b\u5f55",
	"invalid session":                   "\u767b\u5f55\u4f1a\u8bdd\u5df2\u5931\u6548\uff0c\u8bf7\u91cd\u65b0\u767b\u5f55",
	"missing api token":                 "\u7f3a\u5c11 API Token\uff0c\u8bf7\u68c0\u67e5\u8bf7\u6c42\u5934",
	"invalid api token":                 "API Token \u65e0\u6548\u6216\u5df2\u5931\u6548",
	"invalid username or password":      "\u8d26\u53f7\u6216\u5bc6\u7801\u9519\u8bef",
	"user is disabled":                  "\u5f53\u524d\u8d26\u53f7\u5df2\u88ab\u7981\u7528\uff0c\u8bf7\u8054\u7cfb\u7ba1\u7406\u5458",
	"old password is incorrect":         "\u539f\u5bc6\u7801\u8f93\u5165\u4e0d\u6b63\u786e",
	"invalid role id":                   "\u89d2\u8272 ID \u4e0d\u5408\u6cd5",
	"role not found":                    "\u89d2\u8272\u4e0d\u5b58\u5728\u6216\u5df2\u88ab\u5220\u9664",
	"invalid tenant id":                 "\u79df\u6237 ID \u4e0d\u5408\u6cd5",
	"tenant not found":                  "\u79df\u6237\u4e0d\u5b58\u5728\u6216\u5df2\u88ab\u5220\u9664",
	"invalid user id":                   "\u7528\u6237 ID \u4e0d\u5408\u6cd5",
	"user not found":                    "\u7528\u6237\u4e0d\u5b58\u5728\u6216\u5df2\u88ab\u5220\u9664",
	"role code is required":             "\u89d2\u8272\u7f16\u7801\u4e0d\u80fd\u4e3a\u7a7a",
	"role name is required":             "\u89d2\u8272\u540d\u79f0\u4e0d\u80fd\u4e3a\u7a7a",
	"tenantid is required":              "\u79df\u6237\u4e0d\u80fd\u4e3a\u7a7a",
	"tenant admin missing tenant scope": "\u79df\u6237\u7ba1\u7406\u5458\u7f3a\u5c11\u6240\u5c5e\u79df\u6237\uff0c\u8bf7\u8054\u7cfb\u7ba1\u7406\u5458\u68c0\u67e5\u8d26\u53f7\u6570\u636e",
	"builtin role cannot be modified":   "\u5185\u7f6e\u89d2\u8272\u4e0d\u5141\u8bb8\u4fee\u6539",
	"role is in use":                    "\u8be5\u89d2\u8272\u5df2\u88ab\u4f7f\u7528\uff0c\u6682\u65f6\u65e0\u6cd5\u5220\u9664",
	"invalid permission id":             "\u5b58\u5728\u65e0\u6548\u7684\u6743\u9650\u6807\u8bc6\uff0c\u8bf7\u5237\u65b0\u540e\u91cd\u8bd5",
	"system role cannot bind tenant":    "\u7cfb\u7edf\u7ea7\u89d2\u8272\u4e0d\u80fd\u7ed1\u5b9a\u79df\u6237",
	"role tenant scope mismatch":        "\u89d2\u8272\u4f5c\u7528\u8303\u56f4\u4e0e\u79df\u6237\u914d\u7f6e\u4e0d\u5339\u914d",
	"invalid role scope":                "\u89d2\u8272\u4f5c\u7528\u8303\u56f4\u4e0d\u5408\u6cd5",
	"cannot delete current user":        "\u4e0d\u80fd\u5220\u9664\u5f53\u524d\u767b\u5f55\u8d26\u53f7",
}

func normalizeErrorMessage(err error, httpStatus int) string {
	if err == nil {
		return fallbackStatusMessage(httpStatus)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "\u6570\u636e\u4e0d\u5b58\u5728\u6216\u5df2\u88ab\u5220\u9664"
	}

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return translateValidationError(validationErrs)
	}

	var mysqlErr *mysqlDriver.MySQLError
	if errors.As(err, &mysqlErr) {
		if message := translateMySQLError(mysqlErr); message != "" {
			return message
		}
	}

	raw := strings.TrimSpace(err.Error())
	if raw == "" {
		return fallbackStatusMessage(httpStatus)
	}

	if message, ok := errorMessageMap[strings.ToLower(raw)]; ok {
		return message
	}

	lowerRaw := strings.ToLower(raw)
	if strings.Contains(lowerRaw, "duplicate entry") {
		if message := translateDuplicateMessage(raw); message != "" {
			return message
		}
		return "\u6570\u636e\u5df2\u5b58\u5728\uff0c\u8bf7\u52ff\u91cd\u590d\u63d0\u4ea4"
	}

	if strings.Contains(lowerRaw, "foreign key constraint fails") {
		return "\u5f53\u524d\u6570\u636e\u5b58\u5728\u5173\u8054\u8bb0\u5f55\uff0c\u65e0\u6cd5\u76f4\u63a5\u64cd\u4f5c\uff0c\u8bf7\u5148\u89e3\u9664\u5173\u8054\u540e\u91cd\u8bd5"
	}

	if strings.Contains(lowerRaw, "invalid character") || strings.Contains(lowerRaw, "cannot unmarshal") {
		return "\u8bf7\u6c42\u53c2\u6570\u683c\u5f0f\u4e0d\u6b63\u786e\uff0c\u8bf7\u68c0\u67e5\u540e\u91cd\u8bd5"
	}

	if strings.Contains(lowerRaw, "required") && strings.Contains(lowerRaw, "validation for") {
		return "\u8bf7\u6c42\u53c2\u6570\u6821\u9a8c\u5931\u8d25\uff0c\u8bf7\u68c0\u67e5\u5fc5\u586b\u9879\u540e\u91cd\u8bd5"
	}

	if httpStatus >= http.StatusInternalServerError {
		return "\u670d\u52a1\u5904\u7406\u5931\u8d25\uff0c\u8bf7\u7a0d\u540e\u91cd\u8bd5"
	}

	return raw
}

func fallbackRawErrorMessage(err error, httpStatus int) string {
	if err != nil && strings.TrimSpace(err.Error()) != "" {
		return err.Error()
	}

	switch httpStatus {
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusForbidden:
		return "permission denied"
	case http.StatusNotFound:
		return "not found"
	default:
		return "internal server error"
	}
}

func fallbackStatusMessage(httpStatus int) string {
	if httpStatus >= http.StatusInternalServerError {
		return "\u670d\u52a1\u5904\u7406\u5931\u8d25\uff0c\u8bf7\u7a0d\u540e\u91cd\u8bd5"
	}
	if httpStatus == http.StatusUnauthorized {
		return "\u672a\u767b\u5f55\u6216\u767b\u5f55\u5df2\u5931\u6548\uff0c\u8bf7\u91cd\u65b0\u767b\u5f55"
	}
	if httpStatus == http.StatusForbidden {
		return "\u5f53\u524d\u8d26\u53f7\u6ca1\u6709\u6743\u9650\u6267\u884c\u8be5\u64cd\u4f5c"
	}
	if httpStatus == http.StatusNotFound {
		return "\u8bf7\u6c42\u7684\u6570\u636e\u4e0d\u5b58\u5728"
	}
	return "\u8bf7\u6c42\u5904\u7406\u5931\u8d25\uff0c\u8bf7\u68c0\u67e5\u540e\u91cd\u8bd5"
}

func translateValidationError(errs validator.ValidationErrors) string {
	if len(errs) == 0 {
		return "\u8bf7\u6c42\u53c2\u6570\u6821\u9a8c\u5931\u8d25\uff0c\u8bf7\u68c0\u67e5\u540e\u91cd\u8bd5"
	}

	err := errs[0]
	field := err.Field()
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s\u4e0d\u80fd\u4e3a\u7a7a", field)
	case "oneof":
		return fmt.Sprintf("%s\u53d6\u503c\u4e0d\u5408\u6cd5", field)
	case "min", "max", "len":
		return fmt.Sprintf("%s\u957f\u5ea6\u4e0d\u7b26\u5408\u8981\u6c42", field)
	default:
		return fmt.Sprintf("%s\u53c2\u6570\u4e0d\u5408\u6cd5\uff0c\u8bf7\u68c0\u67e5\u540e\u91cd\u8bd5", field)
	}
}

func translateMySQLError(err *mysqlDriver.MySQLError) string {
	switch err.Number {
	case 1062:
		if message := translateDuplicateMessage(err.Message); message != "" {
			return message
		}
		return "\u6570\u636e\u5df2\u5b58\u5728\uff0c\u8bf7\u52ff\u91cd\u590d\u63d0\u4ea4"
	case 1451:
		return "\u5f53\u524d\u6570\u636e\u5df2\u88ab\u5173\u8054\u4f7f\u7528\uff0c\u6682\u65f6\u65e0\u6cd5\u5220\u9664\u6216\u4fee\u6539"
	case 1452:
		return "\u5173\u8054\u6570\u636e\u4e0d\u5b58\u5728\u6216\u5df2\u5931\u6548\uff0c\u8bf7\u68c0\u67e5\u540e\u91cd\u8bd5"
	default:
		return ""
	}
}

func translateDuplicateMessage(message string) string {
	matches := duplicateEntryPattern.FindStringSubmatch(message)
	if len(matches) != 3 {
		return ""
	}

	value := matches[1]
	key := normalizeConstraintKey(matches[2])
	if label, ok := fieldLabelByKey[key]; ok {
		if key == "uk_role_permission" {
			return "\u76f8\u540c\u7684\u89d2\u8272\u6743\u9650\u914d\u7f6e\u5df2\u5b58\u5728\uff0c\u8bf7\u52ff\u91cd\u590d\u63d0\u4ea4"
		}
		return fmt.Sprintf("%s\u201c%s\u201d\u5df2\u5b58\u5728\uff0c\u8bf7\u66f4\u6362\u540e\u91cd\u8bd5", label, value)
	}

	return fmt.Sprintf("\u6570\u636e\u201c%s\u201d\u5df2\u5b58\u5728\uff0c\u8bf7\u52ff\u91cd\u590d\u63d0\u4ea4", value)
}

func normalizeConstraintKey(key string) string {
	key = strings.TrimSpace(key)
	key = strings.Trim(key, "`")
	key = strings.Trim(key, "'")
	key = strings.Trim(key, `"`)
	return strings.ToLower(key)
}
