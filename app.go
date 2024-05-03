package main

import (
	"github.com/xissg/userManageSystem/router"
)

func main() {

	r := router.NewServer()

	//开启服务器
	r.Run(":8081")

}
