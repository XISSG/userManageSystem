package main

import (
	"github.com/xissg/userManageSystem/router"
)

// @title		用户管理系统
// @version	0.1
// @author		xissg
func main() {

	r := router.NewServer()

	//开启服务器
	r.Run(":8081")

}
