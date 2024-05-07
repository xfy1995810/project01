package gpool

import (
	"dcss/global"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/viper"
)

func Init() error {
	var err error
	gPoolSize := viper.GetInt("GPool.Size")
	global.GPool, err = ants.NewPool(gPoolSize)
	if err != nil {
		return err
	}

	return nil
}
