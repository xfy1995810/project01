package static

import (
	"dcss/embed"
	"dcss/global"
	"dcss/pkg/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterStaticRoute 注册静态资源相关路由
func RegisterStaticRoute(r *gin.Engine) {

	r.GET("/", func(c *gin.Context) {
		c.Writer.WriteHeader(200)
		b, err := embed.Dist.ReadFile("dist/index.html")
		if err != nil {
			global.LOG.Errorln("read embed index.html err, err: ", err)
			resp.Fail(c, nil, "读取index.html失败")
			return
		}
		_, err = c.Writer.Write(b)
		if err != nil {
			global.LOG.Errorln("write index.html to resp err, err: ", err)
			resp.Fail(c, nil, "回写index.html失败")
			return
		}

		c.Writer.Header().Add("Accept", "text/html")
		c.Writer.Flush()
	})
	r.StaticFS("/assets", http.FS(embed.AssetsDist))
	r.StaticFS("/static", http.FS(embed.StaticDist))
	r.StaticFileFS("/favicon.ico", "dist/favicon.ico", http.FS(embed.Dist))
	r.StaticFileFS("/serverConfig.json", "dist/serverConfig.json", http.FS(embed.Dist))
	r.StaticFileFS("/logo.svg", "dist/logo.svg", http.FS(embed.Dist))
}
