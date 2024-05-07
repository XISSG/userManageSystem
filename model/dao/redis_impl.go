package dao

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xissg/userManageSystem/model/entity"
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
func (redisService *RedisService) AddUser(user entity.User, ctx *gin.Context) error {
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
func (redisService *RedisService) AddUsers(users []entity.User, ctx *gin.Context) error {
	for _, user := range users {
		err := redisService.AddUser(user, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetUserByName 获取用户信息
func (redisService *RedisService) GetUserByName(name string, ctx *gin.Context) (entity.User, error) {
	resultJson, err := redisService.rdb.Get(ctx, name).Result()
	if err != nil {
		if err == redis.Nil {
			return entity.User{}, err
		}
		return entity.User{}, err
	}

	var user entity.User
	err = json.Unmarshal([]byte(resultJson), &user)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// UpdateUserInfo 更新用户信息
func (redisService *RedisService) UpdateUserAll(user entity.UpdateUser, ctx *gin.Context) error {
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
func (redisService *RedisService) UpdateUserOne(column string, user entity.UpdateUser, ctx *gin.Context) error {
	val, err := redisService.GetUserByName(user.UserName, ctx)
	err = redisService.DeleteUserByName(user.UserName, ctx)
	if err != nil {
		return err
	}

	switch column {
	case "user_account":
		val.UserAccount = user.UserAccount
	case "user_password":
		val.UserPassword = user.UserPassword
	case "avatar_url":
		val.AvatarUrl = user.AvatarUrl
	case "user_role":
		val.UserRole = user.UserRole
	default:
		return errors.New("invalid column name")
	}

	err = redisService.AddUser(val, ctx)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserByName 删除用户
func (redisService *RedisService) DeleteUserByName(name string, ctx *gin.Context) error {
	err := redisService.rdb.Del(ctx, name).Err()
	if err != nil {
		return err
	}

	return nil
}
