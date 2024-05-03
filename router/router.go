package router

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	redisstore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/xissg/userManageSystem/controller"
	"github.com/xissg/userManageSystem/middleware"
	"github.com/xissg/userManageSystem/model"
	"github.com/xissg/userManageSystem/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// read configuration from yaml config file
func readConfig(database string) *model.Config {
	viper.AddConfigPath("./conf")
	viper.SetConfigName(database)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var config *model.Config
	err = viper.Unmarshal(&config)
	return config
}

// init database connection
func initDB() *gorm.DB {
	config := readConfig("mysql")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func initRedisStore() redisstore.Store {

	config := readConfig("redis")
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	store, err := redisstore.NewStoreWithDB(10, "tcp", address, config.Password, config.Database, []byte(config.Secret))
	if err != nil {
		panic(err)
	}
	return store
}

func initRedis() *redis.Client {
	config := readConfig("redis")
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: config.Password, // 没有密码，默认值
		DB:       0,               // 默认DB 0
	})

	return rdb
}

// NewServer 开启服务器
func NewServer() *gin.Engine {
	// 初始化gin引擎
	r := gin.New()

	//注册序列化模型用户session对象存储
	gob.Register(model.User{})
	gob.Register(model.UserSession{})

	//cors中间件
	r.Use(middleware.CORS)

	// 初始化数据库连接
	db := initDB()
	rdb := initRedis()
	//初始化session
	store := initRedisStore()
	r.Use(sessions.Sessions("session", store))

	//注入依赖
	sessionService := service.NewSessionService(store)
	userService := service.NewUserService(db)
	redisService := service.NewRedisService(rdb)
	userController := controller.NewUserController(userService, sessionService, redisService)

	//映射路由
	v1 := r.Group("v1")
	{
		v1.POST("/user/register", userController.Register)
		v1.POST("/user/login", userController.Login)
		v1.POST("/user/logout", userController.Logout)
		v1.POST("/user/admin/match", userController.MatchUsers)
		v1.GET("/user/admin/delete/:username", userController.DeleteUser)
	}
	return r
}
