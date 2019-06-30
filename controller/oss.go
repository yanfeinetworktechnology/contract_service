package controller

import (
	"contract_service/common"
	"contract_service/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateSignture 生成oss上传签名
// @Summary 生成oss上传签名
// @Description 生成oss上传签名
// @Tags OSS
// @Param token header string true "token"
// @Accept json
// @Produce json
// @Success 200 {object} model.Message
// @Router /oss/signture [get]
func CreateSignture(c *gin.Context) {
	policyToken := common.GetPolicyToken()

	c.JSON(http.StatusOK, model.Message{
		Data: policyToken,
	})
}
