package service

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xissg/userManageSystem/model/entity"
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
	MysqlService DBService
	RedisService CacheService
}

func NewUserService(m DBService, r CacheService) *UserService {

	return &UserService{
		MysqlService: m,
		RedisService: r,
	}
}

func (userService *UserService) GetUser(username string, ctx *gin.Context) (interface{}, error) {
	result, err := userService.RedisService.GetUserByName(username, ctx)
	if err == redis.Nil {
		result, err = userService.MysqlService.GetUserByName(username)
		if err != nil {
			return nil, err
		}

		_ = userService.RedisService.AddUser(result, ctx)
		return result, nil
	}

	if err != nil {
		return nil, err
	}

	return result, nil

}

// AddUser 写服务，封装了对MySQL和redis的写服务
// 写逻辑写入MySQL的同时，写入redis缓存中
// 写入情况：
// MySQL写入失败
// 直接返回失败
// MySQL写入成功，redis写入失败
// 返回成功
func (userService *UserService) AddUser(user entity.User, ctx *gin.Context) error {
	err := userService.MysqlService.AddUser(user)
	if err != nil {

		return err
	}

	retryCount := 3
	retryTime := time.Second * 2
	for i := 0; i < retryCount; i++ {
		err = userService.RedisService.AddUser(user, ctx)
		if err == nil {
			break
		}

		time.Sleep(retryTime)
	}

	if err != nil {
		log.Println("Redis add  failed")
		return err
	}

	return nil
}

// UpdateUserAll 更新服务，封装了对MySQL和redis的更新服务
// 更新逻辑写入MySQL的同时，写入redis缓存中
// 写入情况：
// MySQL更新失败
// 直接返回失败
// MySQL写入成功，redis写入失败
// 返回成功
func (userService *UserService) UpdateUserAll(user entity.UpdateUser, ctx *gin.Context) error {
	err := userService.MysqlService.UpdateUserAll(user)
	if err != nil {

		return err
	}

	retryCount := 3
	retryTime := time.Second * 2
	for i := 0; i < retryCount; i++ {
		err = userService.RedisService.UpdateUserAll(user, ctx)
		if err == nil {
			break
		}

		time.Sleep(retryTime)
	}

	if err != nil {
		log.Println("Redis cache update failed")

		return err
	}

	return nil
}

func (userService *UserService) UpdateUserOne(column string, user entity.UpdateUser, ctx *gin.Context) error {
	err := userService.MysqlService.UpdateUserOne(column, user)
	if err != nil {

		return err
	}

	retryCount := 3
	retryTime := time.Second * 2
	for i := 0; i < retryCount; i++ {
		err = userService.RedisService.UpdateUserOne(column, user, ctx)
		if err == nil {
			break
		}

		time.Sleep(retryTime)
	}

	if err != nil {
		log.Println("Redis cache update failed")

		return err
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
func (userService *UserService) DeleteUser(username string, ctx *gin.Context) error {
	err := userService.MysqlService.DeleteUserByName(username)
	if err != nil {
		return err
	}

	retryCount := 3
	retryTime := time.Second * 2
	for i := 0; i < retryCount; i++ {
		err = userService.RedisService.DeleteUserByName(username, ctx)
		if err == nil {
			break
		}

		time.Sleep(retryTime)
	}

	if err != nil {
		log.Println("Redis cache update failed")

		return err
	}

	return nil
}
