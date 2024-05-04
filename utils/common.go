package utils

import "github.com/gin-gonic/gin"

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
