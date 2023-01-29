package user

import (
	"github.com/gin-gonic/gin"
)

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/24
*@description:
 */
func Router(router *gin.Engine) {
	router.Any("/wechat/login", WeChatLogin)
}
