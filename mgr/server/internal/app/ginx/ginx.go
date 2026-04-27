package ginx

import (
	"net/http"

	"clearbill/mgr/server/api/vo"
	"github.com/gin-gonic/gin"
)

func ParseJSON(c *gin.Context, obj any) error {
	return c.ShouldBindJSON(obj)
}

func ParseQuery(c *gin.Context, obj any) error {
	return c.ShouldBindQuery(obj)
}

func ResOK(c *gin.Context) {
	ResSuccess(c, vo.StatusResult{Status: vo.OKStatus})
}

func ResList(c *gin.Context, v any) {
	ResSuccess(c, vo.ListResult{List: v})
}

func ResSuccess(c *gin.Context, v any) {
	c.JSON(http.StatusOK, vo.ResponseResult{
		Success: true,
		Data:    v,
	})
	c.Abort()
}

func ResError(c *gin.Context, err error, status ...int) {
	httpStatus := http.StatusInternalServerError
	if len(status) > 0 {
		httpStatus = status[0]
	}

	message := "internal server error"
	if err != nil {
		message = err.Error()
	}

	c.JSON(httpStatus, vo.ResponseResult{
		Success: false,
		Error:   message,
	})
	c.Abort()
}
