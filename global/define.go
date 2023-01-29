/*
   @Author: 请填写作者名称和联系方式
   @Desc: 全局变量
   @Create At: 20.11.20
*/

package global

import "fmt"

type GeedError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *GeedError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

var (
	ErrNotFound  = &GeedError{Code: 10404, Message: "entity not found:"}
	ErrDuplicate = &GeedError{Code: 10405, Message: "duplicate entity:"}

	ErrUnsupportedEvent = &GeedError{Code: 20001, Message: "unsupported event"}
	//登陆相关错误
	ErrLoginAgain = &GeedError{Code: 30001, Message: "Login expired, please login again"}

	ErrInternal = &GeedError{Code: 99999, Message: "internal error:"}
)

//全局常量
const ()

//全局变量
var ()

//status  状态
const (
	Status_Active  = 1
	Status_Delete  = 2
	Status_Stop    = 3
	Status_Examine = 4
)
