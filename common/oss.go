package common

import (
	"contract_service/model"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"hash"
	"io"
	"time"

	"github.com/spf13/viper"
)

var uploadDir = "contract/"
var expireTime int64 = 30

func getGmtIso8601(expireEnd int64) string {
	var tokenExpire = time.Unix(expireEnd, 0).Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

// GetPolicyToken 签名
func GetPolicyToken() model.PolicyToken {
	now := time.Now().Unix()
	expireEnd := now + expireTime
	var tokenExpire = getGmtIso8601(expireEnd)

	//create post policy json
	var config model.ConfigStruct
	config.Expiration = tokenExpire
	var condition []interface{}
	condition = append(condition, "content-length-range")
	condition = append(condition, 0)
	condition = append(condition, 1048576000)

	config.Conditions = append(config.Conditions, condition)

	//calucate signature
	result, _ := json.Marshal(config)
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(viper.GetString("oss.accessKeySecret")))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var policyToken model.PolicyToken
	policyToken.AccessKeyID = viper.GetString("oss.accessKeyId")
	policyToken.Host = viper.GetString("oss.host")
	policyToken.Expire = expireEnd
	policyToken.Signature = string(signedStr)
	policyToken.Directory = uploadDir
	policyToken.Policy = string(debyte)

	return policyToken
}
