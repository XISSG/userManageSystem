package api_response

import (
	"github.com/gin-gonic/gin"
)

// 封装了全局返回信息
const (
	UNKNOWNERR   = 1000
	OPERATIONERR = 2000
	REGISTERERR  = 3000
	PARAMSERR    = 5000
	AUTHERR      = 6000
	SUCCESS      = 7000
)
// ApiResponse 用户生成swagger文档
type ApiResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func NewResponse(data interface{}, msg string) *ApiResponse {

	return &ApiResponse{
		Data: data,
		Msg:  msg,
	}
}


func (api *ApiResponse) Response(code int) *gin.H {
	api.Code = code
	return &gin.H{"code": api.Code, "data": api.Data, "msg": api.Msg}
}

