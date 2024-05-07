package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
	rands "math/rand"
	"strconv"
	"time"
)

// MD5Crypt 生成md5 hash
func MD5Crypt(plainText string) string {
	hash := md5.New()
	hash.Write([]byte(plainText))
	cypher := hash.Sum([]byte(plainText))
	return hex.EncodeToString(cypher)
}

// RandomExpireTime 生成一个随机过期时间，过期时间至少为一天
func RandomExpireTime() time.Duration {
	rands.Seed(time.Now().UnixNano())
	minExpire := time.Minute
	maxExpire := time.Hour
	expire := minExpire + time.Duration(rands.Int63n(int64(maxExpire-minExpire)))
	return expire
}

func NewUuid() int64 {
	uid, _ := strconv.Atoi(uuid.New().String())
	return int64(uid)
}
