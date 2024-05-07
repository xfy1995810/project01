package resp

import (
	"dcss/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Result struct {
	Status    int         `json:"status" default:"200"`
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg" example:"成功"`
	Timestamp string      `json:"timestamp" example:"响应时间"`
}

// Response 统一的返回函数
func Response(ctx *gin.Context, status int, data interface{}, msg string) {
	ctx.JSON(http.StatusOK, Result{
		status,
		data,
		msg,
		utils.GetNowTime(),
	})
}

// Success 成功返回函数
func Success(ctx *gin.Context, data interface{}, msg string) {
	ctx.JSON(http.StatusOK, Result{
		200,
		data,
		msg,
		utils.GetNowTime(),
	})
}

// Fail 失败返回函数
func Fail(ctx *gin.Context, data interface{}, msg string) {
	ctx.JSON(http.StatusOK, Result{
		400,
		data,
		msg,
		utils.GetNowTime(),
	})
}
