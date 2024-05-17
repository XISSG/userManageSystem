package model_user

import (
	"github.com/xissg/userManageSystem/common/constant"
	"github.com/xissg/userManageSystem/utils"
	"time"
)

// User 用户
type User struct {
	// id
	ID string `json:"id" gorm:"column:id; type:varchar(256);primaryKey"`
	// 用户昵称
	UserName string `json:"user_name" gorm:"column:user_name;type:varchar(256)"`
	// 用户账号
	UserAccount string `json:"user_account" gorm:"column:user_account;type:varchar(256)"`
	// 用户头像
	AvatarUrl string `json:"avatar_url" gorm:"column:avatar_url;type:varchar(1024)"`
	// 用户密码
	UserPassword string `json:"user_password" gorm:"column:user_password;type:varchar(512)"`
	// 创建时间
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;type:time;"`
	// 更新时间
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;type:time;"`
	// 是否删除
	IsDelete int8 `json:"is_delete" gorm:"column:is_delete;type:int; default: 0"`
	// 0为匿名用户, 1为登录用户, 2为用户
	UserRole string `json:"user_role" gorm:"column:user_role;type:varchar(64)"`
}

func (u User) TableName() string {
	return "user"
}

type AddUserRequest struct {
	UserName     string `json:"user_name" `
	UserAccount  string `json:"user_account" validate:"required,min=3,max=32"`
	AvatarUrl    string `json:"avatar_url"`
	UserPassword string `json:"user_password" validate:"required,min=7,max=32"`
}

// AddUserToUser 为接收的用户补充字段
func AddUserToUser(addUser AddUserRequest) (user User) {
	user.ID = utils.NewUuid()
	user.UserName = addUser.UserName
	user.UserAccount = addUser.UserAccount
	user.AvatarUrl = addUser.AvatarUrl
	user.UserPassword = utils.MD5Crypt(addUser.UserPassword)
	user.CreateTime = time.Now().UTC()
	user.UpdateTime = time.Now().UTC()
	user.UserRole = constant.Common
	user.IsDelete = constant.ALIVE

	return user
}

type LoginUserRequest struct {
	UserAccount  string `json:"user_account"`
	UserPassword string `json:"user_password"`
}

func LoginUserToUser(loginUser LoginUserRequest) (user User) {
	user.UserAccount = loginUser.UserAccount
	user.UserPassword = utils.MD5Crypt(loginUser.UserPassword)
	return user
}

// 用户更新信息的请求
type UpdateUserRequest struct {
	UserName     string `json:"user_name"`
	AvatarUrl    string `json:"avatar_url"`
	UserPassword string `json:"user_password"`
}

func UpdateUserToUser(oldInfo User, updateUser UpdateUserRequest) User {
	if updateUser.UserName != "" {
		oldInfo.UserName = updateUser.UserName
	}
	if updateUser.AvatarUrl != "" {
		oldInfo.AvatarUrl = updateUser.AvatarUrl
	}
	if updateUser.UserPassword != "" {
		oldInfo.UserPassword = utils.MD5Crypt(updateUser.UserPassword)
	}

	return oldInfo
}

// UserQueryRequest 普通用户查询模型
type UserQueryRequest struct {
	// id
	ID string `json:"id" `
	// 用户昵称
	UserName string `json:"user_name" `
	// 用户账号
	UserAccount string `json:"user_account" `
}

// UserQueryToCommonQuery 将其转为通用查询模型
func UserQueryToCommonQuery(queryModel UserQueryRequest) AdminUserQueryRequest {
	return AdminUserQueryRequest{
		ID:          queryModel.ID,
		UserName:    queryModel.UserName,
		UserAccount: queryModel.UserAccount,
		IsDelete:    constant.ALIVE,
	}
}
func UserQueryToUser(query UserQueryRequest) User {
	return User{
		ID:          query.ID,
		UserName:    query.UserName,
		UserAccount: query.UserAccount,
	}
}

// ReturnUser 返回结果数据
type ReturnUser struct {
	ID          string    `json:"id"`
	UserName    string    `json:"user_name"`
	UserAccount string    `json:"user_account"`
	AvatarUrl   string    `json:"avatar_url"`
	CreateTime  time.Time `json:"create_time" `
}

// UserToReturnUser 将数据库中查询出来的数据进行筛选后返回
func UserToReturnUser(u User) *ReturnUser {

	return &ReturnUser{
		ID:          u.ID,
		UserName:    u.UserName,
		UserAccount: u.UserAccount,
		AvatarUrl:   u.AvatarUrl,
		CreateTime:  u.CreateTime,
	}
}

// UsersToReturnUsers 返回多个用户信息
func UsersToReturnUsers(users []User) []*ReturnUser {
	var returnUsers []*ReturnUser
	for _, user := range users {
		returnUsers = append(returnUsers, UserToReturnUser(user))
	}

	return returnUsers
}
