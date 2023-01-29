package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"gitlab.com/canyinxinxi/main-server/handler/job"
	"gitlab.com/canyinxinxi/main-server/handler/user"
	"gitlab.com/canyinxinxi/main-server/module/gin_proxy"
)

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/19
*@description:
 */
var routers *gin.Engine

func Router() {
	routers = gin_proxy.Handle()
	routers.Use(Recovery)
	routers.Use(gzip.Gzip(gzip.DefaultCompression))
	user.Router(routers)
	job.Router(routers)
}
