package job

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/canyinxinxi/main-server/handler/auth"
)

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/24
*@description:
 */
func Router(router *gin.Engine) {
	jobs := router.Group("/jobs")
	{
		jobs.Use(auth.Auth())
		jobs.POST("/create", Create)
		jobs.POST("/my", MyJob)
		jobs.POST("/update", Update)
		jobs.GET("/info/:jobId", Info)
	}
}
