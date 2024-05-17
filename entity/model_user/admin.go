package model_user

import (
	"github.com/xissg/userManageSystem/utils"
	"time"
)

// 管理员编辑用户信息
type EditUserRequest struct {
	UserName     string `json:"user_name"`
	UserAccount  string `json:"user_account"`
	AvatarUrl    string `json:"avatar_url"`
	UserPassword string `json:"user_password"`
	UserRole     string `json:"user_role"`
}

func EditUserToUser(oldInfo User, editUser EditUserRequest) User {
	if editUser.UserName != "" {
		oldInfo.UserName = editUser.UserName
	}
	if editUser.AvatarUrl != "" {
		oldInfo.AvatarUrl = editUser.AvatarUrl
	}
	if oldInfo.UserRole != "" {
		oldInfo.UserRole = editUser.UserRole
	}
	if editUser.UserPassword != "" {
		oldInfo.UserPassword = utils.MD5Crypt(editUser.UserPassword)
	}

	return oldInfo
}

// 管理员查询模型
type AdminUserQueryRequest struct {
	// id
	ID string `json:"id" `
	// 用户昵称
	UserName string `json:"user_name" `
	// 用户账号
	UserAccount string `json:"user_account" `
	// 创建时间
	CreateTime time.Time `json:"create_time" `
	// 是否删除
	IsDelete int8 `json:"is_delete" `
	// 匿名用户，普通用户，管理员，禁用用户
	UserRole string `json:"user_role" `
}

func AdminUserQueryToUser(admin AdminUserQueryRequest) User {
	return User{
		ID:          admin.ID,
		UserName:    admin.UserName,
		UserAccount: admin.UserAccount,
		CreateTime:  admin.CreateTime,
		UserRole:    admin.UserRole,
		IsDelete:    admin.IsDelete,
	}
}

// ReturnAdminUser 返回给管理员的用户信息
type ReturnAdminUser struct {
	ID          string    `json:"id"`
	UserName    string    `json:"user_name"`
	UserAccount string    `json:"user_account"`
	AvatarUrl   string    `json:"avatar_url"`
	CreateTime  time.Time `json:"create_time" `
	UserRole    string    `json:"user_role"`
}

func UserToAdminReturnUser(user User) ReturnAdminUser {
	return ReturnAdminUser{
		ID:          user.ID,
		UserName:    user.UserName,
		UserAccount: user.UserAccount,
		AvatarUrl:   user.AvatarUrl,
		CreateTime:  user.CreateTime,
		UserRole:    user.UserRole,
	}
}

func UsersToAdminReturnUsers(users []User) []ReturnAdminUser {
	var adminReturnUsers []ReturnAdminUser
	for _, user := range users {
		adminReturnUsers = append(adminReturnUsers, UserToAdminReturnUser(user))
	}
	return adminReturnUsers
}
