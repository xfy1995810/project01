package system

import (
	"dcss/handler/api/v1/system/api"
	"dcss/handler/api/v1/system/menu"
	"github.com/gin-gonic/gin"
)

// RegisterSystemRoute 注册系统相关路由
func RegisterSystemRoute(router *gin.RouterGroup) {
	systemGroup := router.Group("/system")
	{
		// Api相关
		systemGroup.GET("/apis", api.GetApiList)
		systemGroup.GET("/api/:id", api.GetApiByID)
		systemGroup.POST("/api", api.AddApi)
		systemGroup.PUT("/api/:id", api.EditApiByID)
		systemGroup.DELETE("/api/:id", api.DeleteApiByID)
	}
	{
		// 菜单相关
		systemGroup.GET("/menus", menu.GetMenuList)
		systemGroup.GET("/menu/:id", menu.GetMenuByID)
		systemGroup.POST("/menu", menu.AddMenu)
		systemGroup.PUT("/menu/:id", menu.EditMenuByID)
		systemGroup.DELETE("/menu/:id", menu.DeleteMenuByID)
	}
}

// RegisterNoCasbinSystemRoute 注册系统相关路由,不需要casbin鉴权
func RegisterNoCasbinSystemRoute(router *gin.RouterGroup) {
	systemGroup := router.Group("/system")
	{
		// 菜单相关
		systemGroup.GET("/dynamic_menus", menu.GetDynamicMenuList)
	}

}
