package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/src/model"
	"math/big"
)

// IdGenerator 生成用户id
func IdGenerator() (*big.Int, error) {
	id, err := rand.Int(rand.Reader, big.NewInt(0xffffffff))
	if err != nil {
		return nil, err
	}
	return id, nil
}

// MD5Crypt 生成md5 hash
func MD5Crypt(plainText string) string {
	hash := md5.New()
	hash.Write([]byte(plainText))
	cypher := hash.Sum([]byte(plainText))
	return hex.EncodeToString(cypher)
}

const (
	UNKNOWNERR   = 1000
	OPERATIONERR = 2000
	REGISTERERR  = 3000
	LOGINERR     = 4000
	PARAMSERR    = 5000
	AUTHERR      = 6000
	SUCCESS      = 7000
)

func Success(data []*model.ResultUser, msg string) *gin.H {
	return &gin.H{"code": SUCCESS, "data": data, "msg": msg}
}

func Error(code int, msg string) *gin.H {

	switch code {
	case OPERATIONERR:
		return &gin.H{"code": OPERATIONERR, "data": nil, "msg": msg}
	case REGISTERERR:
		return &gin.H{"code": REGISTERERR, "data": nil, "msg": msg}
	case LOGINERR: //
		return &gin.H{"code": LOGINERR, "data": nil, "msg": msg}
	case PARAMSERR:
		return &gin.H{"code": PARAMSERR, "data": nil, "msg": msg}
	case AUTHERR:
		return &gin.H{"code": AUTHERR, "data": nil, "msg": msg}
	}

	return &gin.H{"code": UNKNOWNERR, "data": nil, "msg": msg}
}
