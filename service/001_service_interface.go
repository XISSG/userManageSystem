package service

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/model"
)

// DBService Mysql服务的接口
type DBService interface {
	//添加用户
	AddUser(user model.User) error
	AddUsers(users []model.User) error

	//获取用户信息
	GetUserByName(name string) (model.User, error)

	//更新用户信息
	UpdateUserAll(user model.User) error
	UpdateUserOne(column string, user model.User) error
	UpdateUserName(user model.User) error

	//删除用户
	DeleteUserByName(name string) error

	////用户标签
	//AddUserTags(tags model.Tags) error
	//GetUsersByTags(tags []string) ([]model.Tags, error)
	//UpdateUserTags(tags model.Tags) error
	//DeleteUserTags(tags []string) error
}

// CacheService Redis缓存服务的接口
type CacheService interface {
	//添加用户
	AddUser(user model.User, ctx *gin.Context) error
	AddUsers(users []model.User, ctx *gin.Context) error

	//获取用户信息
	GetUserByName(name string, ctx *gin.Context) (model.User, error)

	//更新用户信息
	UpdateUserAll(user model.User, ctx *gin.Context) error

	//更新一条用户信息
	UpdateUserOne(column string, user model.User, ctx *gin.Context) error
	//删除用户
	DeleteUserByName(name string, ctx *gin.Context) error

	////用户标签
	//AddUserTags(tags model.Tags, ctx *gin.Context) error
	//GetUsersByTag(tag string, ctx *gin.Context) ([]model.Tags, error)
	//GetUsersByTags(tags []string, ctx *gin.Context) ([][]model.Tags, error)
	//UpdateUserTags(oldTags, newTags model.Tags, ctx *gin.Context) error
	//DeleteUserTags(tags model.Tags, ctx *gin.Context) error
}

// SessionService Session服务的接口
type SessionService interface {
	NewOrUpdateSession(c *gin.Context, user model.UserSession) error
	GetSession(c *gin.Context) model.UserSession
	DeleteSession(c *gin.Context) error
}
