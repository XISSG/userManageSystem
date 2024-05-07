package router

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	redisstore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"github.com/xissg/userManageSystem/model/dao"
	"github.com/xissg/userManageSystem/model/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/xissg/userManageSystem/docs"

	"github.com/xissg/userManageSystem/controller"
	"github.com/xissg/userManageSystem/middleware"
	"github.com/xissg/userManageSystem/service"
)

// read configuration from yaml config file
func readConfig(database string) *entity.Config {
	viper.AddConfigPath("./conf")
	viper.SetConfigName(database)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var config *entity.Config
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
func NewServer() {
	//注册序列化模型用户session对象存储
	gob.Register(entity.User{})
	gob.Register(entity.UserSession{})

	// 初始化gin引擎
	r := gin.New()

	//cors中间件
	r.Use(middleware.CORS)

	// 初始化数据库连接
	db := initDB()
	rdb := initRedis()

	//初始化session
	store := initRedisStore()
	r.Use(sessions.Sessions("session", store))

	//注入依赖
	sessionService := dao.NewSessionService(store)
	mysqlService := dao.NewUserService(db)
	redisService := dao.NewRedisService(rdb)
	userService := service.NewUserService(mysqlService, redisService)
	userController := controller.NewUserController(*userService, sessionService)

	//开启定时任务
	cronJob := service.NewCronJob(mysqlService, redisService)
	c := cron.New()
	c.AddFunc("@daily", cronJob.Start)
	c.Start()

	//映射路由
	v1 := r.Group("v1")
	{
		v1.POST("/user/register", userController.Register)

		v1.POST("/user/login", userController.Login)
		v1.GET("/user/logout", userController.Logout)

		v1.POST("/user/update", userController.UpdateUser)
		v1.GET("/user/admin/query/:username", userController.QueryUser)
		v1.GET("/user/admin/delete/:username", userController.DeleteUser)
	}

	//设置swagger api文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//开启服务器
	r.Run(":8081")
}
