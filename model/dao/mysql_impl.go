package dao

import (
	"errors"
	"github.com/xissg/userManageSystem/model/entity"
	"gorm.io/gorm"
)

const ALIVE = 0 //is_delete 字段，默认为不删除用户
const DELETE = 1

type MysqlService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *MysqlService {

	return &MysqlService{
		db: db,
	}
}

// AddUser 添加用户
func (us *MysqlService) AddUser(user entity.User) error {
	err := us.db.AutoMigrate(&entity.User{})
	if err != nil {
		return err
	}

	tx := us.db.Begin()
	if err = tx.Table("user").Create(&user).Error; err != nil {

		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (us *MysqlService) AddUsers(users []entity.User) error {
	err := us.db.AutoMigrate(&entity.User{})
	if err != nil {
		return err
	}

	tx := us.db.Begin()
	for user := range users {
		if err = tx.Table("user").Create(&user).Error; err != nil {
			tx.Rollback()

			return err
		}
	}

	tx.Commit()

	return nil
}

// GetUserByName  根据用户名返回查询用户
func (us *MysqlService) GetUserByName(name string) (entity.User, error) {
	_ = us.db.AutoMigrate(&entity.User{})
	var res entity.User
	tx := us.db.Table("user").Where("user_name = ? AND is_delete = ?", name, ALIVE).First(&res)

	return res, tx.Error
}

// UpdateUserAll  更新用户信息
func (us *MysqlService) UpdateUserAll(user entity.UpdateUser) error {
	err := us.db.AutoMigrate(&entity.User{})
	if err != nil {
		return err
	}

	tx := us.db.Begin()
	res := tx.Table("user").Where("user_name = ? AND is_delete = ?", user.UserName, ALIVE).Select("user_account", "user_password", "avtar_url", "user_role").Updates(user)
	if res.Error != nil {
		tx.Rollback()

		return res.Error
	}

	tx.Commit()

	return nil
}

func (us *MysqlService) UpdateUserOne(column string, user entity.UpdateUser) error {

	err := us.db.AutoMigrate(&entity.User{})
	if err != nil {
		return err
	}

	var res *gorm.DB
	tx := us.db.Begin()
	switch column {
	case "user_account":
		res = tx.Table("user").Where("user_name = ? AND is_delete = ?", user.UserName, ALIVE).Update("user_account", user.UserAccount)

	case "avatar_url":
		res = tx.Table("user").Where("user_name = ? AND is_delete = ?", user.UserName, ALIVE).Update("avatar_url", user.AvatarUrl)

	case "user_password":
		res = tx.Table("user").Where("user_name = ? AND is_delete = ?", user.UserName, ALIVE).Update("user_password", user.UserPassword)

	case "user_role":
		res = tx.Table("user").Where("user_name = ? AND is_delete = ?", user.UserName, ALIVE).Update("user_role", user.UserRole)
	default:
		return errors.New("invalid column name")
	}

	if res.Error != nil {
		tx.Rollback()

		return res.Error
	}

	tx.Commit()

	return nil
}
func (us *MysqlService) UpdateUserName(user entity.User) error {

	err := us.db.AutoMigrate(&entity.User{})
	if err != nil {
		return err
	}

	tx := us.db.Begin()
	res := tx.Table("user").Where("user_name = ? AND is_delete = ?", user.UserName, ALIVE).Update("user_name", user.UserName)
	if res.Error != nil {
		tx.Rollback()

		return res.Error
	}

	tx.Commit()

	return nil
}

// UpdateUserAccount 更新用户账户
func (us *MysqlService) UpdateUserAccount(user entity.User) error {
	err := us.db.AutoMigrate(&entity.User{})
	if err != nil {
		return err
	}

	tx := us.db.Begin()
	res := tx.Table("user").Where("user_name = ? AND is_delete = ?", user.UserName, ALIVE).Update("user_account", user.UserAccount)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	tx.Commit()

	return nil
}

// UpdateUserPassword 更新用户密码
func (us *MysqlService) UpdateUserPassword(user entity.User) error {
	err := us.db.AutoMigrate(&entity.User{})
	if err != nil {
		return err
	}

	tx := us.db.Begin()
	res := tx.Table("user").Where("user_name = ? AND is_delete = ?", user.UserName, ALIVE).Update("user_password", user.UserPassword)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	tx.Commit()

	return nil
}

// UpdateUserAvatar 更新用户头像
func (us *MysqlService) UpdateUserAvatar(user entity.User) error {
	err := us.db.AutoMigrate(&entity.User{})
	if err != nil {
		return err
	}

	tx := us.db.Begin()
	res := tx.Table("user").Where("user_name = ? AND is_delete = ?", user.UserName, ALIVE).Update("avatar_url", user.AvatarUrl)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	tx.Commit()

	return nil
}

// UpdateUserRole 更行用户权益
func (us *MysqlService) UpdateUserRole(user entity.User) error {
	err := us.db.AutoMigrate(&entity.User{})
	if err != nil {
		return err
	}

	tx := us.db.Begin()
	res := tx.Table("user").Where("user_name = ? AND is_delete = ?", user.UserName, ALIVE).Update("user_role", user.UserRole)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	tx.Commit()

	return nil
}

// DeleteUserByName  将isDelete字段设为 1 ,不进行实际删除
func (us *MysqlService) DeleteUserByName(name string) error {
	err := us.db.AutoMigrate(&entity.User{})
	if err != nil {
		return err
	}

	tx := us.db.Begin()
	res := tx.Table("user").Where("user_name = ? AND is_delete = ?", name, ALIVE).Update("is_delete", DELETE)
	if res.Error != nil {
		tx.Rollback()

		return res.Error
	}

	tx.Commit()

	return nil
}
