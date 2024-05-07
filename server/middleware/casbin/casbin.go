package casbin

import (
	"dcss/global"
	"dcss/models"
	"dcss/pkg/resp"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// CasbinHandler 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		u, exist := ctx.Get("user")
		if !exist {
			global.LOG.Errorln("权限验证失败，ctx获取user失败")
			resp.Fail(ctx, nil, "权限验证失败")
			ctx.Abort()
			return
		}
		user, ok := u.(*models.SysUser)
		if !ok {
			global.LOG.Errorln("权限验证失败，类型断言user失败")
			resp.Fail(ctx, nil, "权限验证失败")
			ctx.Abort()
			return
		}

		if user.IsSuper {
			ctx.Next()
			return
		}

		//获取请求的PATH
		path := ctx.Request.URL.Path
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}

		obj := strings.TrimPrefix(path, global.RouterPrefix)
		// 获取请求方法
		act := ctx.Request.Method
		// 获取用户的角色
		// 循环判断角色是否拥有权限
		var success bool
		for _, role := range user.Roles {
			if !role.State {
				continue
			}
			sub := strconv.Itoa(int(role.ID))
			success, _ := global.SyncedCachedEnforcer.Enforce(sub, obj, act)
			if success {
				ctx.Next()
				return
			}
		}
		if !success {
			resp.Fail(ctx, nil, "权限验证失败,请联系管理员添加权限后访问。")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
