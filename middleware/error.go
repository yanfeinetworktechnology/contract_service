package middleware

import (
	"contract_service/common"
	"fmt"
	"log"
	"net/http"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
)

// ErrorHandling 错误处理中间件
func ErrorHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		//从context获取最新的错误
		err := c.Errors.Last()
		if err == nil {
			return
		}
		// 转义
		var metaData common.GeneralReturn
		switch err.Meta.(type) {
		case common.GeneralReturn:
			metaData = err.Meta.(common.GeneralReturn)
		default:
			return
		}
		switch err.Type {
		case gin.ErrorTypePublic:
			// 公开错误 返回对应Http状态码和错误码
			// 如果有自定义消息 写入日志
			if metaData.CustomMessage != "" {
				log.Println(metaData.CustomMessage)
				raven.CaptureMessage(fmt.Sprintf("[custom] %v", metaData.CustomMessage), map[string]string{"type": "custom"})
			}
			c.JSON(metaData.HTTPStatus, metaData.AppErrJSON)
			if metaData.AppErrJSON.ErrCode < 20000 {
				raven.CaptureMessage(fmt.Sprintf("[%v] %v", metaData.AppErrJSON.ErrCode, metaData.AppErrJSON.Message), map[string]string{"type": "system"})
			} else {
				raven.CaptureMessage(fmt.Sprintf("[%v] %v", metaData.AppErrJSON.ErrCode, metaData.AppErrJSON.Message), map[string]string{"type": "application"})
			}
			return
		case gin.ErrorTypePrivate:
			// 如果有自定义消息 写入日志
			if metaData.CustomMessage != "" {
				log.Println(metaData.CustomMessage)
				raven.CaptureMessage(fmt.Sprintf("[custom] %v", metaData.CustomMessage), map[string]string{"type": "custom"})
			}
			break
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"err_code": common.SystemError,
				"message":  common.Errors[common.SystemError],
			})
			log.Println(c.ClientIP(), "SYSTEM ERROR")
			raven.CaptureMessage(fmt.Sprintf("[%v] %v", 10001, common.Errors[10001]), map[string]string{"type": "system"})
			return
		}
	}
}
