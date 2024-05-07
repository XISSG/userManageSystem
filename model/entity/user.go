package entity

import (
	"github.com/xissg/userManageSystem/utils"
	"time"
)

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
	ID       int64  `json:"id"`
	UserName string `json:"user_name"`
	Role     int32  `json:"role"`
}

type AddUser struct {
	UserName     string    `json:"user_name"`
	UserAccount  string    `json:"user_account"`
	AvatarUrl    string    `json:"avatar_url"`
	UserPassword string    `json:"user_password"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
	UserRole     int32     `json:"user_role"`
}

// AddUserToUser 为接收的用户补充字段
func AddUserToUser(addUser AddUser) (user User) {
	user.ID = utils.NewUuid()
	user.UserName = addUser.UserName
	user.UserAccount = addUser.UserAccount
	user.AvatarUrl = addUser.AvatarUrl
	user.UserPassword = addUser.UserPassword
	user.CreateTime = time.Now().UTC()
	user.UpdateTime = time.Now().UTC()
	user.UserRole = addUser.UserRole
	user.IsDelete = false

	return user
}

type LoginUser struct {
	UserName     string `json:"user_name"`
	UserPassword string `json:"user_password"`
}

type UpdateUser struct {
	UserName     string `json:"user_name"`
	UserAccount  string `json:"user_account"`
	AvatarUrl    string `json:"avatar_url"`
	UserPassword string `json:"user_password"`
	UserRole     int32  `json:"user_role"`
}

// ReturnUser 返回结果数据
type ReturnUser struct {
	ID          int64     `json:"id"`
	UserName    string    `json:"user_name"`
	UserAccount string    `json:"user_account"`
	AvatarUrl   string    `json:"avatar_url"`
	CreateTime  time.Time `json:"create_time" `
	UserRole    int32     `json:"user_role"`
}

// SafetyUser 将数据库中查询出来的数据进行筛选后返回
func SafetyUser(u User) *ReturnUser {

	return &ReturnUser{
		ID:          u.ID,
		UserName:    u.UserName,
		UserAccount: u.UserAccount,
		AvatarUrl:   u.AvatarUrl,
		CreateTime:  u.CreateTime,
		UserRole:    u.UserRole,
	}
}

// SafetyUsers 将数据库中查询出来的数据进行筛选后返回
func SafetyUsers(u []User) []*ReturnUser {
	result := make([]*ReturnUser, len(u))
	for i, v := range u {
		result[i] = SafetyUser(v)
	}

	return result
}

// CountParams 解析参数
func CountParams(user UpdateUser) (int, string) {
	count := 0
	hash := make(map[string]string)
	if user.UserPassword != "" {
		count++
		hash["column"] = "user_password"
	}
	if user.UserAccount != "" {
		count++
		hash["column"] = "user_account"
	}
	if user.AvatarUrl != "" {
		count++
		hash["column"] = "avatar_url"
	}
	if user.UserRole != 0 {
		count++
		hash["column"] = "user_role"
	}
	if count == 1 {
		return count, hash["column"]
	}
	return count, ""
}
