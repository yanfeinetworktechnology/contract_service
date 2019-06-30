package controller

import (
	base_common "base_service/common"

	"contract_service/common"
	"contract_service/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewContract 上传新的合同
// @Summary 新的合同信息
// @Description 新的合同信息
// @Tags 合同相关
// @Param token header string true "token"
// @Param contract body model.ContractRequest true "电子合同信息"
// @Accept json
// @Produce json
// @Success 200 {object} model.Message
// @Router /contract/new [post]
func NewContract(c *gin.Context) {
	userID, exist := c.Get("userID")
	// 获取数据失败
	if common.FuncHandler(c, exist, true, common.SystemError) {
		return
	}

	var contractRequest model.ContractRequest
	if common.FuncHandler(c, c.BindJSON(&contractRequest), nil, common.ParameterError) {
		return
	}

	var contract model.Contract
	contract.ElectronicContract = contractRequest.ElectronicContract
	contract.PaperContract = contractRequest.PaperContract
	contract.UploadUserID = userID.(int64)

	db := base_common.GetMySQL()
	tx := db.Begin()

	err := db.Create(&contract).Error
	if common.FuncHandler(c, err, nil, common.DatabaseError) {
		tx.Rollback()
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, model.Message{
		Msg: "保存合同成功",
	})
}
