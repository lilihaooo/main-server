package switch_context

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"gitlab.com/canyinxinxi/main-server/global"
	"google.golang.org/protobuf/runtime/protoiface"
	"io/ioutil"
)

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/26
*@description:
 */
func SwitchContextForRequest(ctx *gin.Context, requestInfo protoiface.MessageV1) error {
	ContentType := ctx.GetHeader("Content-Type")
	switch ContentType {
	case "application/octet-stream", "application/x-base64":
		baseBody, _ := ioutil.ReadAll(ctx.Request.Body)
		Body, err := global.Base64Decode(string(baseBody))
		if err != nil {
			log.Error().Msgf("wx user err: %v", err)
			return err
		}
		if err = proto.Unmarshal(Body, requestInfo); err != nil {
			return err
		}
		return nil
	case "application/json", "application/json;charset=UTF-8":
		if err := ctx.Bind(requestInfo); err != nil {

			return err
		}
		return nil
	default:
		return global.ErrNotFound
	}
	return nil
}
