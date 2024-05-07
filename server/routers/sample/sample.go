package sample

import (
	"dcss/handler/api/v1/sample"

	"github.com/gin-gonic/gin"
)

func RegisterSampleInfoRoute(router *gin.RouterGroup) {
	systemGroup := router.Group("/sample")
	{
		// 配置样品信息相关
		systemGroup.GET("/sample", sample.GetSampleObjList)
		systemGroup.POST("/sample", sample.AddSampleObj)
		systemGroup.PUT("/sample/:id", sample.EditSampleObjByID)
		systemGroup.DELETE("/sample/:id", sample.DeleteSampleObjByID)
	}
}
