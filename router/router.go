package router

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"github.com/xissg/userManageSystem/controller"
	"github.com/xissg/userManageSystem/model"
	"github.com/xissg/userManageSystem/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化数据库连接
func initDB() *gorm.DB {

	viper.AddConfigPath("./conf")
	viper.SetConfigName("mysql")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var config *model.Config
	err = viper.Unmarshal(&config)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	var db *gorm.DB
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

// NewServer 开启服务器
func NewServer(r *gin.Engine) {

	//注册序列化模型用户session对象存储
	gob.Register(model.UserSession{})

	//初始化session
	store := sessions.NewCookieStore(
		securecookie.GenerateRandomKey(32),
		securecookie.GenerateRandomKey(32))

	// 初始化数据库连接
	db := initDB()

	//注入依赖
	userService := service.NewService(db, store)
	userController := controller.NewUserController(userService)

	//映射路由
	v1 := r.Group("v1")
	{
		v1.POST("/user/register", userController.Register)
		v1.POST("/user/login", userController.Login)
		v1.POST("/user/logout", userController.Logout)
		v1.GET("/user/admin/query/:username", userController.QueryUsers)
		v1.GET("/user/admin/delete/:username", userController.DeleteUser)
	}

}
