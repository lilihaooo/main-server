package gin_proxy

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gitlab.com/canyinxinxi/main-server/config"
	"net/http"
)

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/19
*@description:
 */
var HttpServer = new(Server)

type Server struct {
	g   *gin.Engine
	srv *http.Server
}

//初始化gin框架
//sequence 表示启动编号
//appEngine true 代表结合PaaS, false代表直连
func Init(sequence int) {
	HttpServer.g = gin.New()

	if viper.GetBool("server.debug") {
		gin.SetMode(gin.DebugMode)
		HttpServer.g.Use(gin.Logger(), gin.Recovery())
	} else {
		gin.SetMode(gin.ReleaseMode)
		HttpServer.g.Use(gin.Recovery())
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.GetConfig().Server.Port),
		Handler: HttpServer.g,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("listen: err")
		}
	}()
	HttpServer.srv = srv
	log.Info().Msgf("[%d] Gin Server 初始化 监听端口: %d", sequence, config.GetConfig().Server.Port)
}
func Handle() *gin.Engine {
	return HttpServer.g
}
