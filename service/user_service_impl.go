package service

import (
	"github.com/xissg/userManageSystem/model"
	"gorm.io/gorm"
)

const ALIVE = 0 //is_delete 字段，默认为不删除用户
const DELETE = 1

type UserServiceImpl struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserServiceImpl {

	return &UserServiceImpl{
		db: db,
	}
}

// AddUser 添加用户
func (us *UserServiceImpl) AddUser(user model.User) error {

	if err := us.db.Table("user").Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (us *UserServiceImpl) AddUsers(users []model.User) error {
	tx := us.db.Begin()
	for user := range users {
		if err := tx.Table("user").Create(&user).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// GetUserByName  根据用户名返回查询用户
func (us *UserServiceImpl) GetUserByName(name string) (model.User, error) {

	var res model.User
	tx := us.db.Table("user").Where("user_name = ? AND is_delete = ?", name, ALIVE).First(&res)

	return res, tx.Error

}

// GetUsersByTags 根据查询条件返回用户列表
func (us *UserServiceImpl) GetUsersByTags(tags string) ([]model.User, error) {

	var res []model.User
	tx := us.db.Table("user").Where("tags = ?", tags, ALIVE).Find(&res)

	return res, tx.Error

}

// UpdateUserAccount 更新用户账户
func (us *UserServiceImpl) UpdateUserAccount(user model.User) error {

	tx := us.db.Begin()
	res := tx.Table("user").Where("user_name = ? AND is_delete = ?", user.UserName, ALIVE).Update("user_account = ?", user.UserAccount)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	tx.Commit()
	return nil
}

func (us *UserServiceImpl) UpdateUserPassword(user model.User) error {

	tx := us.db.Begin()
	res := tx.Table("user").Where("user_name = ? AND is_delete = ?", user.UserName, ALIVE).Update("user_password = ?", user.UserPassword)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	tx.Commit()
	return nil
}

func (us *UserServiceImpl) UpdateUserAvatar(user model.User) error {

	tx := us.db.Begin()
	res := tx.Table("user").Where("user_name = ? AND is_delete = ?", user.UserName, ALIVE).Update("avatar_url = ?", user.AvatarUrl)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	tx.Commit()
	return nil
}

func (us *UserServiceImpl) UpdateUserRole(user model.User) error {

	tx := us.db.Begin()
	res := tx.Table("user").Where("user_name = ? AND is_delete = ?", user.UserName, ALIVE).Update("user_role = ?", user.UserRole)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	tx.Commit()
	return nil
}

func (us *UserServiceImpl) UpdateUserTags(user model.User) error {

	tx := us.db.Begin()
	res := tx.Table("user").Where("user_name = ? AND is_delete = ?", user.UserName, ALIVE).Update("tags = ?", user.Tags)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	tx.Commit()
	return nil
}

// DeleteUserByName  将isDelete字段设为 1 ,不进行实际删除
func (us *UserServiceImpl) DeleteUserByName(name string) error {

	tx := us.db.Begin()
	res := tx.Table("user").Where("user_name = ? and is_delete = ?", name, ALIVE).Update("is_delete = ? ", DELETE)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	return nil
}

func (us *UserServiceImpl) DeleteUsersByTags(tags []string) error {
	tx := us.db.Begin()
	res := tx.Table("user").Where("tags = ? and is_delete = ?", tags, ALIVE).Update("is_delete = ? ", DELETE)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	return nil
}
