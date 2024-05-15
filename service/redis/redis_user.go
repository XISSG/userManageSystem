package redis

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xissg/userManageSystem/entity/modeluser"
	"github.com/xissg/userManageSystem/utils"
)

type UserService struct {
	rdb *redis.Client
}

func NewRedisService(rdb *redis.Client) *UserService {
	return &UserService{
		rdb: rdb,
	}
}

/**
 * @description 添加用户
 * @param modeluser
 * @param ctx
 * @return error
 * @author xissg
 */
func (redisService *UserService) AddUser(user modeluser.User, ctx *gin.Context) error {
	err := redisService.set(user.UserAccount, user, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (redisService *UserService) AddUsers(users []modeluser.User, ctx *gin.Context) error {

	for _, user := range users {
		err := redisService.AddUser(user, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
 * @description 获取用户信息
 * @param name
 * @param ctx
 * @return error
 * @author xissg
 */
func (redisService *UserService) GetUser(accountName string, ctx *gin.Context) (modeluser.User, error) {
	resultJson, err := redisService.get(accountName, ctx)

	if err != nil {
		return modeluser.User{}, err
	}

	var user modeluser.User
	err = json.Unmarshal([]byte(resultJson), &user)
	if err != nil {
		return modeluser.User{}, err
	}
	return user, nil
}

/**
 * @description 更新用户信息
 * @param modeluser
 * @param ctx
 * @return error
 * @author xissg
 */
func (redisService *UserService) UpdateUser(user modeluser.User, ctx *gin.Context) error {
	err := redisService.delete(user.UserAccount, ctx)
	if err != nil {
		return err
	}

	return nil
}

/**
 * @description 删除用户
 * @param name
 * @param ctx
 * @return error
 * @author xissg
 */
func (redisService *UserService) DeleteUser(accountName string, ctx *gin.Context) error {
	err := redisService.delete(accountName, ctx)
	if err != nil {
		return err
	}

	return nil
}

// redis的通用逻辑
func (redisService *UserService) set(key string, value interface{}, ctx *gin.Context) error {
	str, err := json.Marshal(value)
	if err != nil {
		return err
	}
	status := redisService.rdb.Set(ctx, key, str, utils.RandomExpireTime())
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (redisService *UserService) get(key string, ctx *gin.Context) (string, error) {
	str, err := redisService.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return str, nil
}

func (redisService *UserService) delete(key string, ctx *gin.Context) error {

	err := redisService.rdb.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
