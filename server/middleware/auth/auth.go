package auth

import (
	"dcss/global"
	"dcss/models"
	"dcss/pkg/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// Auth token认证中间件
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//	获取authorization header
		tokenStr := ctx.GetHeader("Authorization")
		if len(tokenStr) == 0 {
			tokenStr = ctx.Query("token")
		}
		//	validate header format
		if strings.TrimSpace(tokenStr) == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足, Authorization头部验证失败",
			})
			ctx.Abort()
			return
		}
		tokenStr = tokenStr[7:]

		token, claims, err := utils.ParseToken(tokenStr)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足，token验证失败",
			})
			ctx.Abort()
			return
		}
		//	验证通过后，获取claims中的userID
		userID := claims.UserID
		db := global.DB
		var user models.SysUser
		err = db.Preload("Roles").First(&user, userID).Error

		//	如果用户不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足, 该用户不存在",
			})
			ctx.Abort()
			return
		}

		// 如果用户被禁用
		if !user.State {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "该用户已被禁用",
			})
			ctx.Abort()
			return
		}

		//	用户存在, 将user和role的信息写入上下文
		ctx.Set("user", &user)
		ctx.Set("userID", userID)
		ctx.Next()
	}
}
