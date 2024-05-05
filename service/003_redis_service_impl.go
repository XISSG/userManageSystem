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

// GetUserByName 获取用户信息
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

// UpdateUserInfo 更新用户信息
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

// DeleteUserByName 删除用户
func (redisService *RedisService) DeleteUserByName(name string, ctx *gin.Context) error {
	err := redisService.rdb.Del(ctx, name).Err()
	if err != nil {
		return err
	}
	return nil
}

// AddUserTags 添加用户标签
//func (redisService *RedisService) AddUserTags(tags model.Tags, ctx *gin.Context) error {
//	tagsJson, err := json.Marshal(tags)
//	if err != nil {
//		return err
//	}
//	for _, tag := range tags.Tags {
//		status := redisService.rdb.SAdd(ctx, tag, tagsJson, utils.RandomExpireTime())
//		if status.Err() != nil {
//			return status.Err()
//		}
//	}
//	return nil
//}
//
//// GetUsersByTag 返回符合特定标签的用户
//func (redisService *RedisService) GetUsersByTag(tag string, ctx *gin.Context) ([]model.Tags, error) {
//
//	var result model.Tags
//	var results []model.Tags
//	resultsJson, err := redisService.rdb.SMembers(ctx, tag).Result()
//
//	for _, resultJson := range resultsJson {
//
//		err = json.Unmarshal([]byte(resultJson), &result)
//		results = append(results, result)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return results, nil
//}
//
//// GetUsersByTags 获取所有标签的所有用户
//func (redisService *RedisService) GetUsersByTags(tags []string, ctx *gin.Context) ([][]model.Tags, error) {
//
//	var results [][]model.Tags
//	for _, tag := range tags {
//		result, err := redisService.GetUsersByTag(tag, ctx)
//		if err != nil {
//			return nil, err
//		}
//		results = append(results, result)
//	}
//	return results, nil
//}
//
//// UpdateUserTags 更新用户标签
//func (redisService *RedisService) UpdateUserTags(oldTags, newTags model.Tags, ctx *gin.Context) error {
//
//	err := redisService.DeleteUserTags(oldTags, ctx)
//	err = redisService.AddUserTags(newTags, ctx)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//// DeleteUserTags 删除用户标签
//func (redisService *RedisService) DeleteUserTags(tags model.Tags, ctx *gin.Context) error {
//
//	for _, tag := range tags.Tags {
//		resultJson, err := json.Marshal(tags)
//		if err != nil {
//			return err
//		}
//		status := redisService.rdb.SRem(ctx, tag, resultJson)
//		if status.Err() != nil {
//			return status.Err()
//		}
//	}
//	return nil
//}
