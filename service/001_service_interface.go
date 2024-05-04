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
	GetUsersByTags(tags string) ([]model.User, error)

	//更新用户信息
	UpdateUser(user model.User) error
	UpdateUserName(user model.User) error
	UpdateUserAccount(user model.User) error
	UpdateUserPassword(user model.User) error
	UpdateUserAvatar(user model.User) error
	UpdateUserRole(user model.User) error
	UpdateUserTags(user model.User) error

	//删除用户

	DeleteUserByName(name string) error
}

// CacheService Redis缓存服务的接口
type CacheService interface {
	//添加用户

	AddUser(user model.User, ctx *gin.Context) error
	AddUsers(users []model.User, ctx *gin.Context) error

	//获取用户信息

	GetUserByName(name string, ctx *gin.Context) (model.User, error)
	GetUsersByTags(tags string, ctx *gin.Context) ([]model.User, error)

	//更新用户信息

	UpdateUserInfo(user model.User, ctx *gin.Context) error

	//删除用户

	DeleteUserByName(name string, ctx *gin.Context) error
}

// SessionService Session服务的接口
type SessionService interface {
	NewOrUpdateSession(c *gin.Context, user model.UserSession) error
	GetSession(c *gin.Context) model.UserSession
	DeleteSession(c *gin.Context) error
}