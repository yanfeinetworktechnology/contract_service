package middleware

import (
	"context"
	"contract_service/common"
	"fmt"
	"regexp"
	pb "user_service/service/proto"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// TokenHandling 中间件，检查token
func TokenHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenWhiteList = []string{"/docs"}
		var certificationWhiteList = []string{}

		var requestURL = c.Request.RequestURI

		for _, v := range tokenWhiteList {
			match, _ := regexp.MatchString(v, requestURL)
			if match {
				c.Next()
				return
			}
		}

		token := c.Request.Header.Get("token")

		conn, err := grpc.Dial(viper.GetString("token_server.address"), grpc.WithInsecure())
		if common.FuncHandler(c, err, nil, common.TokenServiceError) {
			return
		}
		defer conn.Close()
		client := pb.NewTokenClient(conn)

		// Contact the server and print out its response.
		response, err := client.Check(context.Background(), &pb.TokenRequest{Token: token})
		if common.FuncHandler(c, err, nil, common.TokenServiceError) {
			return
		}

		userID := response.UserId
		status := response.Status

		fmt.Println(userID, status)

		if common.FuncHandler(c, status != -1, true, common.TokenExpired) {
			c.Abort()
			return
		}
		if common.FuncHandler(c, status != -2, true, common.TokenInvalid) {
			c.Abort()
			return
		}
		if common.FuncHandler(c, status != -3, true, common.DatabaseError) {
			c.Abort()
			return
		}

		// 检查是否需要认证，不需要则直接通过
		for _, v := range certificationWhiteList {
			match, _ := regexp.MatchString(v, requestURL)
			if match {
				c.Set("userID", userID)
				c.Next()
				return
			}
		}
		if common.FuncHandler(c, status != -4, true, common.NoCertification) {
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("userID", userID)
		c.Next()
	}
}
