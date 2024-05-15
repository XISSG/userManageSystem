package service

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xissg/userManageSystem/entity/modeluser"
	"github.com/xissg/userManageSystem/service/mysql"
	redis2 "github.com/xissg/userManageSystem/service/redis"
	"log"
	"time"
)

//读服务封装了对redis和MySQL读服务

// UserService 读逻辑，优先根据查询条件去redis进行读取，如果没有读取到数据则去MySQL中进行读取
// 读取情况:
// redis读取成功：
// 返回成功
// redis读取失败,mysql读取成功：
// 将结果写入redis,返回成功
// redis读取失败，mysql读取失败:
// 返回失败
type UserService struct {
	MysqlService mysql.UserService
	RedisService redis2.UserService
}

func NewUserService(m mysql.UserService, r redis2.UserService) *UserService {

	return &UserService{
		MysqlService: m,
		RedisService: r,
	}
}

// AddUser 写服务，封装了对MySQL和redis的写服务
// 写逻辑写入MySQL的同时，写入redis缓存中
// 写入情况：
// MySQL写入失败
// 直接返回失败
// MySQL写入成功，redis写入失败
// 返回成功
func (userService *UserService) AddUser(user modeluser.User, ctx *gin.Context) error {
	var redisErr error
	ch := make(chan error)

	go func() {
		mysqlErr := userService.MysqlService.AddUser(user)
		if mysqlErr != nil {
			log.Println("Mysql add failed")
		}
		ch <- mysqlErr
		close(ch)
	}()

	retryCount := 3
	retryTime := time.Second * 2
	for i := 0; i < retryCount; i++ {
		redisErr = userService.RedisService.AddUser(user, ctx)
		if redisErr == nil {
			break
		}

		time.Sleep(retryTime)
	}

	if redisErr != nil {
		log.Println("Redis add  failed")
	}

	mysqlErr := <-ch
	if mysqlErr != nil {
		return mysqlErr
	}

	return nil
}

func (userService *UserService) GetUser(userAccount string, ctx *gin.Context) (modeluser.User, error) {
	res := make(chan modeluser.User)
	errChan := make(chan error)

	// 尝试从 Redis 中获取用户信息
	result, redisErr := userService.RedisService.GetUser(userAccount, ctx)

	// 在协程中尝试从 MySQL 中获取用户信息
	go func() {
		mysqlResult, mysqlErr := userService.MysqlService.GetUser(userAccount)
		res <- mysqlResult
		errChan <- mysqlErr
	}()

	mysqlRes := <-res
	mysqlErr := <-errChan

	// 检查 MySQL 中是否存在错误
	if mysqlErr != nil {
		return modeluser.User{}, mysqlErr
	}

	// 如果 Redis 返回了空结果，则将 MySQL 中获取到的用户信息写入 Redis 中
	if redisErr == redis.Nil {
		_ = userService.RedisService.AddUser(mysqlRes, ctx)
		return mysqlRes, nil
	}

	// 如果 Redis 存在其他错误，则返回错误信息
	if redisErr != nil {
		return modeluser.User{}, redisErr
	}

	return result, nil
}

// UpdateUserInfo 更新服务，封装了对MySQL和redis的更新服务
// 更新逻辑写入MySQL的同时，写入redis缓存中
// 写入情况：
// MySQL更新失败
// 直接返回失败
// MySQL写入成功，redis写入失败
// 返回成功

func (userService *UserService) UpdateUserInfo(user modeluser.User, ctx *gin.Context) error {
	var redisErr error
	ch := make(chan error)

	go func() {
		mysqlErr := userService.MysqlService.UpdateUser(user)
		ch <- mysqlErr
		close(ch)
	}()

	//错误重试
	retryCount := 3
	retryTime := time.Second * 2
	for i := 0; i < retryCount; i++ {
		redisErr = userService.RedisService.UpdateUser(user, ctx)
		if redisErr == nil {
			break
		}

		time.Sleep(retryTime)
	}

	mysqlErr := <-ch
	if mysqlErr != nil {
		log.Println("Mysql update error: ", mysqlErr)
		return mysqlErr
	}

	if redisErr != nil {
		log.Println("Redis cache update failed")
		return redisErr
	}

	return nil
}

// DeleteUser 删除服务，封装了对MySQL和redis的删除服务
// 删除逻辑删除MySQL的同时，删除redis缓存
// 删除情况：
// MySQL删除失败或redis删除失败
// 直接返回失败
// MySQL删除成功，redis删除成功
// 返回成功
func (userService *UserService) DeleteUser(accountName string, ctx *gin.Context) error {
	var mysqlErr error
	var redisErr error

	ch := make(chan error)

	go func() {
		mysqlErr := userService.MysqlService.DeleteUser(accountName)
		ch <- mysqlErr
		close(ch)
	}()

	retryCount := 3
	retryTime := time.Second * 2
	for i := 0; i < retryCount; i++ {
		redisErr = userService.RedisService.DeleteUser(accountName, ctx)
		if redisErr == nil {
			break
		}

		time.Sleep(retryTime)
	}

	mysqlErr = <-ch
	if mysqlErr != nil {
		return mysqlErr
	}

	if redisErr != nil {
		log.Println("Redis cache update failed")
		return redisErr
	}

	return nil
}
