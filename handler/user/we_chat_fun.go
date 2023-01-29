package user

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"gitlab.com/canyinxinxi/main-server/config"
	"gitlab.com/canyinxinxi/main-server/global"
	"gitlab.com/canyinxinxi/main-server/pb/login"
	"io"
	"net/http"
	"net/url"
)

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/24
*@description:
 */
func code2session(code string) (login.WeChatSessionResponse, error) {

	log.Info().Msgf("code: %s", code)

	params := url.Values{}
	params.Set("appid", config.GetConfig().WeChart.AppID)
	params.Set("secret", config.GetConfig().WeChart.Secret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")

	client := http.Client{}
	request, err := http.NewRequest(
		http.MethodGet, config.GetConfig().WeChart.EndpointCode2Session+"?"+params.Encode(), nil)
	if err != nil {
		log.Error().Msgf("create wechat code2session request error: %v", err)
		return login.WeChatSessionResponse{}, global.ErrInternal
	}
	response, err := client.Do(request)
	if err != nil {
		log.Error().Msgf("send wechat code2session request error: %v", err)
		return login.WeChatSessionResponse{}, global.ErrInternal
	}
	if response.StatusCode != http.StatusOK {
		log.Error().Msgf("wechat code2session response error: %d %s", response.StatusCode, response.Status)
		return login.WeChatSessionResponse{}, global.ErrInternal
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error().Msgf("read wechat code2session response error: %v", err)
		return login.WeChatSessionResponse{}, global.ErrInternal
	}
	var result login.WeChatSessionResponse
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Error().Msgf("unmarshal wechat code2session response error: %v", err)
		return login.WeChatSessionResponse{}, global.ErrInternal
	}

	return result, nil

}
