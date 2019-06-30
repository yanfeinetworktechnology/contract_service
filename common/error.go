package common

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// FuncHandler 统一错误处理
// i 传入error,bool,int
// judge 不触发错误的值 非error环境下有效
// 如果触发了错误 return True
// Example:
// 1. common.FuncHandler(c, c.BindJSON(&x), nil, http.StatusBadRequest, 20301)
// ==   if(c.BindJSON(&x) != nil){
// 			c.JSON(http.StatusBadRequest, gin.H{
//			"err_code": 20301,
//			"message":  common.Errors[20301],
//			})
// 	 	}
// 2. common.FuncHandler(c, c.BindJSON(&x), nil, http.StatusBadRequest, 20301,fmt.Sprintf("BindJson fail with %v",x))
// ==   if(c.BindJSON(&x) != nil){
// 			log.Println(fmt.Sprintf("BindJson fail with %v",x))
// 			c.JSON(http.StatusBadRequest, gin.H{
//			"err_code": 20301,
//			"message":  common.Errors[20301],
//			})
// 	 	}
// 3. common.FuncHandler(c, isOdd(2), true, fmt.Sprintf("%d is even",2))
// ==   if(isOdd(2) != true){
// 			log.Println(fmt.Sprintf("%d is even",2))
// 	 	}
func FuncHandler(c *gin.Context, i interface{}, judge interface{}, option ...interface{}) bool {
	generalReturn := buildErrorMeta(option)
	errType := gin.ErrorTypePrivate
	// http返回码和错误码齐全则为公开错误
	if generalReturn.HTTPStatus != 0 || generalReturn.AppErrJSON.ErrCode != 0 {
		errType = gin.ErrorTypePublic
	}
	switch i.(type) {
	case nil:
		return false
	case error:
		c.Error(i.(error)).SetMeta(generalReturn).SetType(errType)
		return true
	case bool:
		if i.(bool) == judge.(bool) {
			return false
		}
		if generalReturn.CustomMessage != "" {
			c.Error(fmt.Errorf(generalReturn.CustomMessage)).SetMeta(generalReturn).SetType(errType)
		} else if generalReturn.AppErrJSON.Message != "" {
			c.Error(fmt.Errorf(generalReturn.AppErrJSON.Message)).SetMeta(generalReturn).SetType(errType)
		} else {
			c.Error(fmt.Errorf("no err")).SetMeta(generalReturn).SetType(errType)
		}
		return true
	}
	return true
}
func buildErrorMeta(option []interface{}) GeneralReturn {
	var generalReturn GeneralReturn
	for _, v := range option {
		switch v.(type) {
		case int:
			// RFC 2616 HTTP Status Code 是3位数字代码
			if v.(int) >= 1000 {
				generalReturn.AppErrJSON.ErrCode = v.(int)
				generalReturn.AppErrJSON.Message = Errors[v.(int)]
			} else {
				generalReturn.HTTPStatus = v.(int)
			}
			break
		case string:
			generalReturn.CustomMessage = v.(string)
			break
		}
	}
	return generalReturn
}

// GeneralReturn 通用码
type GeneralReturn struct {
	CustomMessage string
	HTTPStatus    int
	AppErrJSON    appErrJSON
}
type appErrJSON struct {
	ErrCode int    `json:"status"`
	Message string `json:"message"`
}

// 错误码
const (
	OK = 0

	SystemError         = 10001
	DatabaseError       = 10002
	ServiceUnavailable  = 10003
	ParameterError      = 10004
	ResourceUnavailable = 10005
	Maintenance         = 10006
	NoToken             = 10007
	TokenExpired        = 10008
	TokenInvalid        = 10009
	LoginError          = 10010
	NoCertification     = 10011
	TokenServiceError   = 10012
)

// Errors 错误码
var Errors = map[int]string{

	OK: "OK",

	// 系统级错误
	SystemError:         "系统错误",
	DatabaseError:       "数据库错误",
	ServiceUnavailable:  "暂无法提供服务",
	ParameterError:      "参数错误",
	ResourceUnavailable: "资源不可达",
	Maintenance:         "正在维护",
	NoToken:             "无Token",
	TokenExpired:        "Token过期",
	TokenInvalid:        "Token无效",
	LoginError:          "用户名或密码错误",
	NoCertification:     "未实名认证",
	TokenServiceError:   "Token Service错误",
}
