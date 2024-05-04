package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xissg/userManageSystem/model"
	"github.com/xissg/userManageSystem/utils"
)

type RedisService struct {
	rdb *redis.Client
}

func NewRedisService(rdb *redis.Client) *RedisService {
	return &RedisService{
		rdb: rdb,
	}
}

// AddUser 添加用户
func (redisService *RedisService) AddUser(user model.User, ctx *gin.Context) error {
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}
	status := redisService.rdb.Set(ctx, user.UserName, userJson, utils.RandomExpireTime())
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

// AddUsers 添加多个用户
func (redisService *RedisService) AddUsers(users []model.User, ctx *gin.Context) error {
	for _, user := range users {
		err := redisService.AddUser(user, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

//获取用户信息

func (redisService *RedisService) GetUserByName(name string, ctx *gin.Context) (model.User, error) {
	resultJson, err := redisService.rdb.Get(ctx, name).Result()
	if err != nil {
		if err == redis.Nil {
			return model.User{}, err
		}
		return model.User{}, err
	}
	var user model.User
	err = json.Unmarshal([]byte(resultJson), &user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

//更新用户信息

func (redisService *RedisService) UpdateUserInfo(user model.User, ctx *gin.Context) error {
	val, err := redisService.GetUserByName(user.UserName, ctx)
	err = redisService.DeleteUserByName(user.UserName, ctx)
	if err != nil {
		return err
	}

	val.UserAccount = user.UserAccount
	val.UserPassword = user.UserPassword
	val.AvatarUrl = user.AvatarUrl
	val.UserRole = user.UserRole

	err = redisService.AddUser(val, ctx)
	if err != nil {
		return err
	}
	return nil
}

//删除用户

func (redisService *RedisService) DeleteUserByName(name string, ctx *gin.Context) error {
	err := redisService.rdb.Del(ctx, name).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetUsersByTags 匹配用户
func (redisService *RedisService) GetUsersByTags(tags string, ctx *gin.Context) ([]model.User, error) {
	//TODO:
	return []model.User{}, nil
}
