package embed

import (
	"dcss/global"
	"embed"
	"io/fs"
)

var (
	//go:embed dist/*
	Dist       embed.FS
	AssetsDist fs.FS
	StaticDist fs.FS
)

func init() {
	var err error
	AssetsDist, err = fs.Sub(Dist, "dist/assets")
	if err != nil {
		global.LOG.Errorln("初始化assets dist 目录失败, err: ", err)
	}
	StaticDist, err = fs.Sub(Dist, "dist/static")
	if err != nil {
		global.LOG.Errorln("初始化static dist 目录失败, err: ", err)
	}
}
