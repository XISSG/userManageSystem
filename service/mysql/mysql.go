package mysql

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Config struct {
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Database        string        `yaml:"database"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
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

// InitDB init database connection
func initDB() *gorm.DB {
	config := readConfig("mysql.yaml")

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

	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(config.MaxIdleConns)       //设置最大连接数
	sqlDb.SetMaxOpenConns(config.MaxOpenConns)       //设置最大的空闲连接数
	sqlDb.SetConnMaxLifetime(config.ConnMaxLifetime) //设置最大连接时间

	return db

}
