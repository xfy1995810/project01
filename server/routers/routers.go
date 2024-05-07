package routers

import (
	"dcss/docs"
	"dcss/global"
	"dcss/handler/api/v1/auth/user"
	"dcss/handler/api/v1/system/captcha"
	"dcss/handler/api/v1/welcome"
	authMiddleware "dcss/middleware/auth"
	"dcss/middleware/casbin"
	"dcss/middleware/cors"
	"dcss/middleware/limiter"
	"dcss/pkg/resp"
	"dcss/routers/auth"
	"dcss/routers/sample"
	"dcss/routers/static"
	"dcss/routers/system"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRoute 初始化路由信息
func InitRoute() {
	gin.SetMode("debug")
	gin.DisableConsoleColor()
	logger, err := initLogger()
	if err != nil {
		log.Panic("生成gin日志失败,err: ", err)
	}
	gin.DefaultWriter = logger

	r := gin.Default()
	r.Use(cors.Cors())
	if err = r.SetTrustedProxies(nil); err != nil {
		global.LOG.Errorln("set trust proxies failed, err: ", err)
	}
	collectRoute(r)
	webPort := viper.GetString("SERVER.PORT")
	panic(r.Run(fmt.Sprintf(":%s", webPort)))
}

func collectRoute(r *gin.Engine) {
	r.GET("/ping", func(ctx *gin.Context) {
		resp.Success(ctx, 200, "pong")
	})

	// 无效地址处理
	r.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/#/error/404")
	})

	// 静态资源
	static.RegisterStaticRoute(r)

	// 后端Api
	apiV1 := r.Group(global.RouterPrefix)
	// 处理放爆破中间件
	apiV1 = limiter.WarpLimitMiddleware(apiV1)

	// 无需鉴权的组
	publicGroup := apiV1.Group("")
	{
		publicGroup.POST("/login", user.Login)

		//	验证码
		publicGroup.GET("/captcha", captcha.Captcha)
	}
	// 需鉴权的路由组
	privateGroup := apiV1.Group("")
	privateGroup.Use(authMiddleware.Auth())

	// casbin api 权限路由组
	CasbinPrivateGroup := privateGroup.Group("")
	CasbinPrivateGroup.Use(casbin.CasbinHandler())

	//privateWsGroup := privateGroup.Group("/ws")

	// casbin ws api 权限路由组 暂未使用
	// CasbinPrivateWsGroup := CasbinPrivateGroup.Group("/ws")
	{
		// 主页统计数据
		CasbinPrivateGroup.GET("/welcome", welcome.GetWelcome)
		// 权限相关
		auth.RegisterAuthRoute(CasbinPrivateGroup)
		// 系统相关
		system.RegisterSystemRoute(CasbinPrivateGroup)
		system.RegisterNoCasbinSystemRoute(privateGroup)
		//	配置文件-变更操作
		sample.RegisterSampleInfoRoute(CasbinPrivateGroup)
	}

	//	swagger
	if viper.GetBool("SERVER.DEBUG") {
		docs.SwaggerInfo.BasePath = global.RouterPrefix
		apiV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}

func initLogger() (io.Writer, error) {
	logPath := viper.GetString("LOG.PATH")
	logAbsDir := filepath.Join(global.BaseDir, filepath.Dir(logPath))
	logName := "http.log"
	logMaxAge := viper.GetInt("Log.MaxAge")

	if _, err := os.Stat(logAbsDir); os.IsNotExist(err) {
		err := os.MkdirAll(logAbsDir, 0o755)
		if err != nil {
			return nil, fmt.Errorf("create log dir failed, err:%s\n", err.Error())
		}
	}

	logDateName := fmt.Sprintf("%s/%%Y%%m%%d-%s", logAbsDir, logName)
	logLinkName := fmt.Sprintf("%s/%s", logAbsDir, logName)

	logf, err := rotateLogs.New(
		logDateName,
		rotateLogs.WithLinkName(logLinkName),
		rotateLogs.WithRotationTime(24*time.Hour),
		rotateLogs.WithMaxAge(time.Duration(logMaxAge)*24*time.Hour),
	)

	if err != nil {
		return nil, fmt.Errorf("new ratatelog failed, err: %s\n", err.Error())
	}

	return logf, nil
}
