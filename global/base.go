package global

import "encoding/base64"

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/25
*@description:
 */

var base64Codec = base64.URLEncoding.WithPadding(base64.NoPadding)

func Base64(message []byte) string {
	return base64Codec.EncodeToString(message)
}
func Base64Decode(message string) ([]byte, error) {
	return base64Codec.DecodeString(message)
}
