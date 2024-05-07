package main

import (
	"dcss/global"
	"dcss/global/cache"
	"dcss/global/casbin"
	"dcss/global/config"
	"dcss/global/db"
	"dcss/global/gpool"
	"dcss/global/logging"
	"dcss/global/param_check"
	"dcss/pkg/utils"
	"dcss/routers"
)

//go:generate go install github.com/swaggo/swag/cmd/swag@latest
//go:generate swag init

//	@title	Dcss Swagger API
//	@version		1.0
//	@BasePath	/api/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				jwt 认证
func main() {
	// 获取运行程序的绝对路径
	global.BaseDir = utils.GetCurrentAbPathByExecutable()

	// 读取配置文件
	err := config.Init()
	if err != nil {
		panic(err)
	}
	// 初始化日志框架
	err = logging.Init()
	if err != nil {
		panic(err)
	}
	// 初始化gorm
	err = db.InitDB()
	if err != nil {
		global.LOG.Errorln(err)
		panic(err)
	}
	// 初始化CasBin Api权限校验
	casbin.Init()

	// 初始化 协程池
	err = gpool.Init()
	if err != nil {
		global.LOG.Errorln(err)
		panic(err)
	}
	defer global.GPool.Release()

	// 初始化Cache
	cache.Init()

	// 初始化websocket
	//go ws.WebsocketManager.Start()

	// 初始化Validator数据校验
	param_check.InitValidate()

	// 初始化 gin router,并监听对应端口
	routers.InitRoute()
}
