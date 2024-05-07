package cache

import (
	"dcss/global"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"time"
)

func Init() {
	initCaptchaCache()
}

func initCaptchaCache() {
	openCaptchaTimeOut := viper.GetInt("Captcha.Timeout") // 缓存超时时间
	if openCaptchaTimeOut < 30 || openCaptchaTimeOut > 300 {
		openCaptchaTimeOut = 30
	}
	global.CaptchaCache = cache.New(time.Duration(openCaptchaTimeOut)*time.Second, 10*time.Minute)
	global.LOG.Infof("验证码cache初始化完毕, 验证码超时:(%v)s", openCaptchaTimeOut)
}
