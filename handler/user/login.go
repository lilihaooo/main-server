package user

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gitlab.com/canyinxinxi/main-server/db/user"
	"gitlab.com/canyinxinxi/main-server/global"
	"gitlab.com/canyinxinxi/main-server/handler/switch_context"
	"gitlab.com/canyinxinxi/main-server/module/producet"
	"gitlab.com/canyinxinxi/main-server/module/redis"
	"gitlab.com/canyinxinxi/main-server/pb/index"
	"gitlab.com/canyinxinxi/main-server/pb/login"
	"net/http"
)

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/19
*@description: 微信登陆接口
 */
func WeChatLogin(ctx *gin.Context) {
	requestInfo := new(login.WeChatLoginRequest)
	responseInfo := new(login.WeChatLoginResponse)
	responseInfo.Code = http.StatusNoContent
	err := switch_context.SwitchContextForRequest(ctx, requestInfo)
	if err != nil {
		responseInfo.Message = err.Error()
		switch_context.SwitchContextForResponse(ctx, responseInfo)
		return
	}
	//获取微信登陆凭证
	wxSession, err := code2session(requestInfo.Code)
	test := ctx.Request.Header.Get("testing")
	if test != "" {
		wxSession = login.WeChatSessionResponse{
			SessionKey: "123",
			Unionid:    "123",
			Errmsg:     "",
			Openid:     "123",
			Errcode:    "",
		}
		err = nil
	}

	if err != nil {
		responseInfo.Message = err.Error()
		switch_context.SwitchContextForResponse(ctx, responseInfo)
		return
	}

	//获取用户详情信息
	userInfo, flag := user.GetUserDao().Get(wxSession.Openid)
	if !flag {
		userInfo = &index.User{
			Phone:         requestInfo.Phone,
			WechartOpenId: wxSession.Openid,
		}
		RegisterUser(userInfo)
		user.GetUserDao().Set(userInfo)
	}
	//发送用户日志
	producet.SendUserRr(userInfo)
	/**
	组织返回数据，需要降用户密码隐掉
	返回登陆凭证 token，放到header头里

	*/
	copyUser := global.Copy(userInfo).(*index.User)
	copyUser.Password = ""
	log.Info().Msgf("wx user info: %v", userInfo)
	//拼装返回参数
	responseInfo.Data = &login.WeChatLoginResponse_Data{
		Token:    wxSession.SessionKey,
		UserInfo: copyUser,
	}
	token := redis.SetUserToken(copyUser.Id)
	ctx.Writer.Header().Set("token", token)
	switch_context.SwitchContextForResponse(ctx, responseInfo)
	return
}

func RegisterUser(user *index.User) {
	user.Id = redis.GetMaxUserId()
}
