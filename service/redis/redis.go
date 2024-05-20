package redis

import (
	"fmt"
	redisstore "github.com/gin-contrib/sessions/redis"
	"github.com/spf13/viper"
)

type Config struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Secret   string `yaml:"secret"`
}

func readConfig(database string) *Config {
	viper.AddConfigPath("./conf")
	viper.SetConfigName(database)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var config *Config
	err = viper.Unmarshal(&config)
	return config
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

//func initRedis() *redis.Client {
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
