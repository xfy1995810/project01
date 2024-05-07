package welcome

import (
	"dcss/global"
	"dcss/models"
	"dcss/pkg/resp"
	"github.com/gin-gonic/gin"
)

// GetWelcome
// @Tags     获取主页相关数据
// @Summary  获取主页相关数据
// @Produce  application/json
// @Success  200   {object}  resp.Result{data=models.RespWelcome}  "获取主页相关数据"
// @Security ApiKeyAuth
// @Router   /welcome [get]
func GetWelcome(ctx *gin.Context) {
	db := global.DB
	// 定义返回的数据结构体
	var welcomeData models.RespWelcome
	var err error

	err = db.Model(&models.SysUser{}).Count(&welcomeData.UserCount).Error
	if err != nil {
		resp.Fail(ctx, nil, "获取用户数失败")
		return
	}

	resp.Success(ctx, welcomeData, "获取主页数据成功")
}
