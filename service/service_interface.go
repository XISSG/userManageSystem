package service

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/model"
)

type UserService interface {

	//添加用户

	AddUser(user model.User) error
	AddUsers(users []model.User) error

	//获取用户信息

	GetUserByName(name string) (model.User, error)
	GetUsersByTags(tags string) ([]model.User, error)

	//更新用户信息

	UpdateUserAccount(user model.User) error
	UpdateUserPassword(user model.User) error
	UpdateUserAvatar(user model.User) error
	UpdateUserRole(user model.User) error
	UpdateUserTags(user model.User) error

	//删除用户

	DeleteUserByName(name string) error
	DeleteUsersByTags(tags []string) error
}

type SessionService interface {
	NewOrUpdateSession(c *gin.Context, user model.UserSession) error
	GetSession(c *gin.Context) model.UserSession
	DeleteSession(c *gin.Context) error
}

type RedisService interface {
	//TODO: define some methods
}
