package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/middleware"
	"github.com/xissg/userManageSystem/router"
)

func main() {

	// 初始化gin引擎
	r := gin.New()

	//cors中间件
	r.Use(middleware.CORS)

	router.NewServer(r)

	//开启服务器
	_ = r.Run(":8080")

}
