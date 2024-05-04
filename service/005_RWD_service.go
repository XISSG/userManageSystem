package service

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xissg/userManageSystem/model"
)

//读服务封装了对redis和MySQL读服务

// RWDService 读逻辑，优先根据查询条件去redis进行读取，如果没有读取到数据则去MySQL中进行读取
// 读取情况:
// redis读取成功：
// 返回成功
// redis读取失败,mysql读取成功：
// 将结果写入redis,返回成功
// redis读取失败，mysql读取失败:
// 返回失败
type RWDService struct {
	MysqlService DBService
	RedisService CacheService
}

func NewRWDService(m DBService, r CacheService) *RWDService {
	return &RWDService{
		MysqlService: m,
		RedisService: r,
	}
}

func (rw *RWDService) Read(user model.User, ctx *gin.Context) (interface{}, error) {
	result, err := rw.RedisService.GetUserByName(user.UserName, ctx)
	if err == redis.Nil {
		result, err = rw.MysqlService.GetUserByName(user.UserName)
		if err != nil {
			return nil, err
		}
		_ = rw.RedisService.AddUser(result, ctx)
		return result, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Add 写服务，封装了对MySQL和redis的写服务
// 写逻辑写入MySQL的同时，写入redis缓存中
// 写入情况：
// MySQL写入失败
// 直接返回失败
// MySQL写入成功，redis写入失败
// 返回成功
func (rw *RWDService) Add(user model.User, ctx *gin.Context) error {
	err := rw.MysqlService.AddUser(user)
	if err != nil {
		return err
	}
	_ = rw.RedisService.AddUser(user, ctx)

	return nil
}

// Update 更新服务，封装了对MySQL和redis的更新服务
// 更新逻辑写入MySQL的同时，写入redis缓存中
// 写入情况：
// MySQL更新失败
// 直接返回失败
// MySQL写入成功，redis写入失败
// 返回成功
func (rw *RWDService) Update(user model.User, ctx *gin.Context) error {
	err := rw.MysqlService.UpdateUser(user)
	if err != nil {
		return err
	}
	_ = rw.RedisService.UpdateUserInfo(user, ctx)

	return nil
}

// Delete 删除服务，封装了对MySQL和redis的删除服务
// 删除逻辑删除MySQL的同时，删除redis缓存
// 删除情况：
// MySQL删除失败或redis删除失败
// 直接返回失败
// MySQL删除成功，redis删除成功
// 返回成功
func (rw *RWDService) Delete(user model.User, ctx *gin.Context) error {
	err := rw.MysqlService.DeleteUserByName(user.UserName)
	if err != nil {
		return err
	}
	err = rw.RedisService.DeleteUserByName(user.UserName, ctx)
	if err != nil {
		return err
	}
	return nil
}
