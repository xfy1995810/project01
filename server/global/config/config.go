package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func setDefaultValue() {
	// SERVER
	viper.SetDefault("SERVER.PORT", 80)
	viper.SetDefault("SERVER.DEBUG", false)
	// 携程池配置
	viper.SetDefault("GPool.Size", 1000)
	// 数据库配置
	viper.SetDefault("DB.PATH", "./sqlite.db")
	// 日志配置
	viper.SetDefault("LOG.PATH", "./logs/dcss.log")
	// 0： panic, 1: fatal, 2: error, 3: warn, 4: info, 5: debug, 6: trace
	viper.SetDefault("LOG.LEVEL", 5)
	viper.SetDefault("LOG.MaxAge", 60)
	// CMD
	viper.SetDefault("CMD.TIMEOUT", 5)
	// 防止爆破每个接口每秒请求次数开关,0不限制
	viper.SetDefault("PreventBlast.frequency", 3)
	// 验证码过期时间
	viper.SetDefault("Captcha.Timeout", 60)
	// 验证码长度
	viper.SetDefault("Captcha.Length", 6)
}

func Init() error {
	viper.SetConfigName("setting")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("init config file failed，err: %s", err.Error())
		}
	}
	setDefaultValue()
	// 设置获取环境变量， RICHAMIL_{变量名大写}
	viper.SetEnvPrefix("richmail")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()
	return nil
}
