package switch_context

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
	"net/http"
)

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/26
*@description:
 */
func SwitchContextForResponse(ctx *gin.Context, responseInfo protoiface.MessageV1) error {
	ContentType := ctx.GetHeader("Content-Type")
	switch ContentType {
	case "application/octet-stream", "application/x-base64":
		if response, err := proto.Marshal(responseInfo); err != nil {
			return err
		} else {
			_, err := ctx.Writer.Write([]byte(base64.StdEncoding.EncodeToString(response)))
			return err
		}
	case "application/json", "application/json;charset=UTF-8":
		if response, err := json.Marshal(responseInfo); err != nil {
			return err
		} else {
			_, err := ctx.Writer.Write(response)
			return err
		}
	default:

		ctx.JSON(http.StatusOK, responseInfo)
		return nil

	}
	return nil
}
