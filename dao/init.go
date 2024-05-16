package dao

import (
	"fmt"
	redisstore "github.com/gin-contrib/sessions/redis"
	"github.com/spf13/viper"
	"github.com/xissg/userManageSystem/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// read configuration from yaml config file
func readConfig(database string) *conf.Config {
	viper.AddConfigPath("./conf")
	viper.SetConfigName(database)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var config *conf.Config
	err = viper.Unmarshal(&config)
	return config
}

// InitDB init database connection
func InitDB() *gorm.DB {
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

func InitRedisStore() redisstore.Store {

	config := readConfig("redis")
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	store, err := redisstore.NewStoreWithDB(10, "tcp", address, config.Password, config.Database, []byte(config.Secret))
	if err != nil {
		panic(err)
	}
	return store
}

//func InitRedis() *redis.Client {
//	config := readConfig("redis")
//	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
//	rdb := redis.NewClient(&redis.Options{
//		Addr:     address,
//		Password: config.Password, // 没有密码，默认值
//		DB:       0,               // 默认DB 0
//	})
//
//	return rdb
//}
