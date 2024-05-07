package utils

import (
	"github.com/gin-gonic/gin"
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
func (api *ApiResponse) Success() *gin.H {

	api.Code = SUCCESS
	return &gin.H{"code": api.Code, "data": api.Data, "msg": api.Msg}
}

func (api *ApiResponse) UnknownERR() *gin.H {

	api.Code = UNKNOWNERR
	return &gin.H{"code": api.Code, "data": api.Data, "msg": api.Msg}
}
func (api *ApiResponse) OperationERR() *gin.H {

	api.Code = OPERATIONERR
	return &gin.H{"code": api.Code, "data": api.Data, "msg": api.Msg}
}
func (api *ApiResponse) RegisterERR() *gin.H {

	api.Code = REGISTERERR
	return &gin.H{"code": api.Code, "data": api.Data, "msg": api.Msg}
}

func (api *ApiResponse) LoginERR() *gin.H {

	api.Code = LOGINERR
	return &gin.H{"code": api.Code, "data": api.Data, "msg": api.Msg}
}

func (api *ApiResponse) ParamsERR() *gin.H {

	api.Code = PARAMSERR
	return &gin.H{"code": api.Code, "data": api.Data, "msg": api.Msg}
}
func (api *ApiResponse) AuthERR() *gin.H {

	api.Code = AUTHERR
	return &gin.H{"code": api.Code, "data": api.Data, "msg": api.Msg}
}
