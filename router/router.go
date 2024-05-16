package router

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/xissg/userManageSystem/controller"
	"github.com/xissg/userManageSystem/dao"
	_ "github.com/xissg/userManageSystem/docs"
	"github.com/xissg/userManageSystem/entity/modeluser"
	"github.com/xissg/userManageSystem/middleware"
	mysql2 "github.com/xissg/userManageSystem/service/mysql"
	redis2 "github.com/xissg/userManageSystem/service/redis"
)

// NewServer 开启服务器
func NewServer() {
	//注册序列化模型用户session对象存储
	gob.Register(modeluser.User{})
	gob.Register(modeluser.UserSession{})

	// 初始化gin引擎
	r := gin.New()

	//cors中间件
	r.Use(middleware.CORS)

	// 初始化数据库连接
	db := dao.InitDB()
	//rdb := dao.InitRedis()

	//初始化session
	store := dao.InitRedisStore()
	r.Use(sessions.Sessions("session", store))

	//注入依赖
	sessionService := redis2.NewSessionService(store)
	mysqlService := mysql2.NewUserService(db)
	userController := controller.NewUserController(*mysqlService, *sessionService)

	//题目相关依赖
	questionMysqlService := mysql2.NewQuestionMysqlService(db)
	questionController := controller.NewQuestionController(questionMysqlService, sessionService)

	//题目提交相关依赖
	qsMysqlService := mysql2.NewQuestionSubmitMysqlService(db)
	qsService := mysql2.NewQuestionMysqlService(db)
	qsController := controller.NewQuestionSubmitController(qsMysqlService, qsService, sessionService)

	//映射路由
	v1 := r.Group("api")
	{
		userGroup := v1.Group("user")
		{
			userGroup.POST("/login", userController.Login)
			userGroup.GET("/logout", userController.Logout)

			userGroup.POST("/register", userController.Register)
			userGroup.POST("/query", userController.QueryUserList)
			userGroup.POST("/update", userController.UpdateUser)
			//后台操作

			userGroup.POST("/admin/query", userController.AdminQueryUserList)
			userGroup.POST("/admin/update", userController.EditUser)
			userGroup.GET("/admin/delete/:useraccount", userController.DeleteUser)
		}
		questionGroup := v1.Group("question")
		{
			questionGroup.GET("/get/:id", questionController.GetQuestion)
			questionGroup.POST("/get", questionController.GetQuestionList)

			questionGroup.POST("/admin/create", questionController.AddQuestion)
			questionGroup.GET("/admin/delete/:id", questionController.DeleteQuestion)
			questionGroup.POST("/admin/update", questionController.UpdateQuestion)
		}

		questionSubmitGroup := v1.Group("submit")
		{
			questionSubmitGroup.POST("/submission", qsController.Submit)
			questionSubmitGroup.GET("/query/:id", qsController.GetQuestionSubmit)
			questionSubmitGroup.POST("/query", qsController.GetQuestionSubmitList)
		}
	}

	//设置swagger api文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//开启服务器
	r.Run(":8082")
}
