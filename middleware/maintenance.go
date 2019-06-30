package middleware

import (
	"contract_service/common"
	"log"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// MaintenanceHandling 维护模式中间件
func MaintenanceHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		if common.FuncHandler(c, viper.GetBool("basic.maintenance"), false, common.Maintenance) {
			log.Println(c.ClientIP(), "Maintenance mode is on")
			raven.CaptureMessage("Maintenance mode is on", map[string]string{"type": "maintenance"})
			c.Abort()
		}

		c.Next()
	}
}
