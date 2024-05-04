package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	rands "math/rand"
	"time"
)

// IdGenerator 生成用户id
func IdGenerator() (*big.Int, error) {
	id, err := rand.Int(rand.Reader, big.NewInt(0xffffffff))
	if err != nil {
		return nil, err
	}
	return id, nil
}

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
	minExpire := 24 * time.Hour
	maxExpire := 30 * 24 * time.Hour
	expire := minExpire + time.Duration(rands.Int63n(int64(maxExpire-minExpire)))
	return expire
}
