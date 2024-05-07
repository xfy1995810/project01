package captcha

import (
	"dcss/global"
	"dcss/models"
	"dcss/pkg/resp"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/spf13/viper"
)

var Store models.LocalCacheStore

// Captcha
// @Tags      验证码相关
// @Summary   生成验证码
// @Produce   application/json
// @Success   200  {object}  resp.Result{data=models.RespCaptcha}  "生成验证码,返回包括随机数id,base64,验证码长度,是否开启验证码"
// @Router    /captcha [get]
func Captcha(ctx *gin.Context) {
	captchaLength := viper.GetInt("Captcha.Length")
	// 字符,公式,验证码配置
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(80, 240, captchaLength, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, &Store)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		global.LOG.Errorln("验证码获取失败!, err: ", err)
		resp.Fail(ctx, nil, "验证码获取失败")
		return
	}

	resp.Success(ctx, models.RespCaptcha{
		CaptchaID:   id,
		CaptchaPath: b64s,
	}, "验证码获取成功")
}
