package limiter

import (
	"dcss/global"
	"dcss/pkg/resp"
	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"time"
)

var Lmt *limiter.Limiter

func WarpLimitMiddleware(routerGroup *gin.RouterGroup) *gin.RouterGroup {
	frequency := viper.GetInt("PreventBlast.frequency")

	if frequency == 0 {
		return routerGroup
	}
	global.LOG.Infoln("放爆破开启：frequency: ", frequency)

	Lmt = tollbooth.NewLimiter(float64(frequency), &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	Lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}).
		SetMethods([]string{"GET", "POST", "PUT", "DELETE"})

	routerGroup.Use(LimitHandler(Lmt))

	return routerGroup
}

// LimitHandler 防止爆破
func LimitHandler(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		httpError := tollbooth.LimitByRequest(lmt, ctx.Writer, ctx.Request)
		if httpError != nil {
			resp.Response(ctx, 429, nil, "请求太过频繁，请稍后再试")
			ctx.Abort()
			return
		} else {
			ctx.Next()
		}
	}
}
