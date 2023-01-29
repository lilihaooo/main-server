package auth

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/canyinxinxi/main-server/global"
	"gitlab.com/canyinxinxi/main-server/module/redis"
	"gitlab.com/canyinxinxi/main-server/pb/login"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Token")
		userId := ctx.Request.Header.Get("user_id")
		if token == "" || userId == "" {
			responseInfo := login.WeChatLoginResponse{}
			responseInfo.Message = global.ErrLoginAgain.Message
			responseInfo.Code = http.StatusForbidden
			ctx.JSON(http.StatusOK, responseInfo)
			ctx.Abort()
			return
		}
		cacheId := redis.GetUserToken(token)
		if cacheId != userId {
			responseInfo := login.WeChatLoginResponse{}
			responseInfo.Message = global.ErrLoginAgain.Message
			responseInfo.Code = http.StatusForbidden
			ctx.JSON(http.StatusOK, responseInfo)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
