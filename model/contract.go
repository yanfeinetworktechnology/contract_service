package model

import "time"

// Contract 数据表字段
type Contract struct {
	ID                 int64  `gorm:"primary_key"`
	ElectronicContract string `gorm:"size:10000" gorm:"not null"`
	PaperContract      string
	UploadUserID       int64     `gorm:"not null"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"not null"`
}

// ContractRequest 合同请求字段
type ContractRequest struct {
	ElectronicContract string `json:"electronic_contract"`
	PaperContract      string `json:"paper_contract"`
}
