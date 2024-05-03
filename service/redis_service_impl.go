package service

import "github.com/redis/go-redis/v9"

type RedisServiceImpl struct {
	rdb *redis.Client
}

func NewRedisService(rdb *redis.Client) *RedisServiceImpl {
	return &RedisServiceImpl{rdb: rdb}
}
