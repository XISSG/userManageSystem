package model

import (
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

// User 用户
type User struct {
	// id
	ID int64 `json:"id" gorm:"column:id; type:uint;primaryKey"`
	// 用户昵称
	UserName string `json:"user_name" gorm:"column:user_name;type:string;size:256"`
	// 用户账号
	UserAccount string `json:"user_account" gorm:"column:user_account;type:string;size:256"`
	// 用户头像
	AvatarUrl string `json:"avatar_url" gorm:"column:avatar_url;type:string;size:1024"`
	// 用户密码
	UserPassword string `json:"user_password" gorm:"column:user_password;type:string;size 512"`
	// 创建时间
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;type:time;"`
	// 更新时间
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;type:time;"`
	// 是否删除
	IsDelete bool `json:"is_delete" gorm:"column:is_delete;type:bool;size:1;default:false"`
	// 0为匿名用户, 1为登录用户, 2为用户
	UserRole int32 `json:"user_role" gorm:"column:user_role;type:int;size:32;default:0"`
}

func (u User) TableName() string {
	return "user"
}

// UserSession 存储的用户session信息
type UserSession struct {
	ID     int64 `json:"id"`
	Expire int64 `json:"expire"`
	Role   int32 `json:"role"`
}

// ResultUser 返回结果数据
type ResultUser struct {
	// id
	ID int64 `json:"id"`
	// 用户昵称
	UserName string `json:"user_name"`
	// 用户账号
	UserAccount string `json:"user_account"`
	// 用户头像
	AvatarUrl string `json:"avatar_url"`
	// 创建时间
	CreateTime time.Time `json:"create_time" `
	// 0为普通用户, 1为会员用户
	UserRole int32 `json:"user_role"`
}

// UserProc 将数据库中查询出来的数据进行筛选后返回
func UserProc(u User) *ResultUser {
	return &ResultUser{
		ID:          u.ID,
		UserName:    u.UserName,
		UserAccount: u.UserAccount,
		AvatarUrl:   u.AvatarUrl,
		CreateTime:  u.CreateTime,
		UserRole:    u.UserRole,
	}
}

// UsersProc 将数据库中查询出来的数据进行筛选后返回
func UsersProc(u []User) []*ResultUser {
	result := make([]*ResultUser, len(u))
	for i, v := range u {
		result[i] = UserProc(v)
	}
	return result
}
