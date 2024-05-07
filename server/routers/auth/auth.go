package auth

import (
	"dcss/handler/api/v1/auth/role"
	"dcss/handler/api/v1/auth/user"
	"dcss/handler/api/v1/system/api"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoute 注册权限相关路由
func RegisterAuthRoute(router *gin.RouterGroup) {
	systemGroup := router.Group("/auth")
	{
		// 用户相关
		systemGroup.GET("/users", user.GetUserList)
		systemGroup.POST("/user", user.AddUser)
		systemGroup.GET("/user/:id", user.GetUserByID)
		systemGroup.PUT("/user/:id", user.EditUserByID)
		systemGroup.DELETE("/user/:id", user.DeleteUserByID)
		systemGroup.PUT("/user/:id/state/", user.EditUserStateByID)
		systemGroup.PUT("/user/:id/password/", user.EditUserPasswordByID)
	}
	{
		// 角色相关
		systemGroup.GET("/all_users", user.GetAllUserList)
		systemGroup.GET("/roles", role.GetRoleList)
		systemGroup.GET("/role/:id", role.GetRoleByID)
		systemGroup.POST("/role", role.AddRole)
		systemGroup.PUT("/role/:id", role.EditRoleByID)
		systemGroup.DELETE("/role/:id", role.DeleteRoleByID)
		systemGroup.PUT("/role/:id/state/", role.EditRoleStateByID)
		systemGroup.POST("/role/:id/api_auth/", role.SetRoleApiAuth)
		systemGroup.GET("/role/:id/api_auth", role.GetRoleApiAuthByID)
		systemGroup.POST("/role/:id/menu_auth/", role.SetRoleMenuAuth)
		systemGroup.GET("/role/:id/menu_auth", role.GetRoleMenuAuthByID)
		systemGroup.GET("/all_apis", api.GetAllApiList)
	}
}
