package main

import (
	"io"
	"time"

	base_common "base_service/common"
	"contract_service/common"
	contract_common "contract_service/common"
	"contract_service/model"

	"contract_service/controller"

	"contract_service/middleware"

	_ "contract_service/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func migrate(db *gorm.DB) {
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_bin auto_increment=1")
	db.AutoMigrate(&model.Contract{})
}

// init 在 main 之前执行
func init() {
	// init config
	contract_common.DefaultConfig()
	contract_common.SetConfig()
	contract_common.WatchConfig()

	// init logger
	base_common.InitLogger()

	// init Database
	db := base_common.InitMySQL()
	// 禁止在表名后面加s
	db.SingularTable(true)
	migrate(db)
}

// @title YANFEI-CONTRACT API
// @version 0.0.1
func main() {
	// Before init router
	if viper.GetBool("basic.debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
		// Redirect log to file
		gin.DisableConsoleColor()
		logFile := base_common.GetLogFile()
		defer logFile.Close()
		gin.DefaultWriter = io.MultiWriter(logFile)
	}

	r := gin.Default()
	// CORS
	if viper.GetBool("basic.debug") {
		r.Use(cors.New(cors.Config{
			// The value of the 'Access-Control-Allow-Origin' header in the
			// response must not be the wildcard '*' when the request's
			// credentials mode is 'include'.
			AllowOrigins:     common.CORS_ALLOW_DEBUG_ORIGINS,
			AllowMethods:     common.CORS_ALLOW_METHODS,
			AllowHeaders:     common.CORS_ALLOW_HEADERS,
			ExposeHeaders:    common.CORS_EXPOSE_HEADERS,
			AllowCredentials: true,
			AllowWildcard:    true,
			MaxAge:           12 * time.Hour,
		}))
	} else {
		// RELEASE Mode
		r.Use(cors.New(cors.Config{
			AllowOrigins:     common.CORS_ALLOW_ORIGINS,
			AllowMethods:     common.CORS_ALLOW_METHODS,
			AllowHeaders:     common.CORS_ALLOW_HEADERS,
			ExposeHeaders:    common.CORS_EXPOSE_HEADERS,
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}
	// middleware
	r.Use(middleware.ErrorHandling())
	r.Use(middleware.MaintenanceHandling())
	r.Use(middleware.TokenHandling())

	// swagger router
	if viper.GetBool("basic.debug") {
		r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 路由
	r.POST("/contract/new", controller.NewContract)
	r.GET("/oss/signture", controller.CreateSignture)

	r.Run("0.0.0.0:" + viper.GetString("basic.port"))
}
