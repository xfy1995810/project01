package global

import (
	"github.com/casbin/casbin/v2"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/panjf2000/ants/v2"
	cache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	DB                   *gorm.DB                     // DB 全局数据库池
	LOG                  *logrus.Logger               // 全局日志
	GPool                *ants.Pool                   // 全局线程池
	BaseDir              string                       // 程序运行的目录
	Validate             *validator.Validate          // Validate 全局Validate数据校验实列
	Trans                ut.Translator                // Trans 全局翻译器
	SyncedCachedEnforcer *casbin.SyncedCachedEnforcer // sync cached enforcer
	RouterPrefix         = "/api/v1"                  //路由前缀
	CaptchaCache         *cache.Cache                 // 验证码本地cache
)
