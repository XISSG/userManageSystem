package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/model"
)

// 封装了全局返回信息
const (
	UNKNOWNERR   = 1000
	OPERATIONERR = 2000
	REGISTERERR  = 3000
	LOGINERR     = 4000
	PARAMSERR    = 5000
	AUTHERR      = 6000
	SUCCESS      = 7000
)

func Success(data interface{}, msg string) *gin.H {
	return &gin.H{"code": SUCCESS, "data": data, "msg": msg}
}

func Error(code int, msg string) *gin.H {
	return &gin.H{"code": code, "data": nil, "msg": msg}
}

// ApiResponse 用户生成swaggerapi文档
type ApiResponse struct {
	Code int              `json:"code"`
	Data model.ResultUser `json:"data"`
	Msg  string           `json:"msg"`
}
