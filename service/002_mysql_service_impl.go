package service

import (
	"github.com/xissg/userManageSystem/model"
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
func (us *MysqlService) AddUser(user model.User) error {
	err := us.db.AutoMigrate(&model.User{})
	if err != nil {
		return err
	}
	if err = us.db.Table("user").Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (us *MysqlService) AddUsers(users []model.User) error {
	err := us.db.AutoMigrate(&model.User{})
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
func (us *MysqlService) GetUserByName(name string) (model.User, error) {
	_ = us.db.AutoMigrate(&model.User{})
	var res model.User
	tx := us.db.Table("user").Where("user_name = ? AND is_delete = ?", name, ALIVE).First(&res)

	return res, tx.Error

}

// UpdateUser 更新用户信息
func (us *MysqlService) UpdateUser(user model.User) error {
	err := us.db.AutoMigrate(&model.User{})
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

func (us *MysqlService) UpdateUserName(user model.User) error {

	err := us.db.AutoMigrate(&model.User{})
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
func (us *MysqlService) UpdateUserAccount(user model.User) error {
	err := us.db.AutoMigrate(&model.User{})
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
func (us *MysqlService) UpdateUserPassword(user model.User) error {
	err := us.db.AutoMigrate(&model.User{})
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
func (us *MysqlService) UpdateUserAvatar(user model.User) error {
	err := us.db.AutoMigrate(&model.User{})
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
func (us *MysqlService) UpdateUserRole(user model.User) error {
	err := us.db.AutoMigrate(&model.User{})
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
	err := us.db.AutoMigrate(&model.User{})
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

// AddUserTags 添加用户标签
//func (us *MysqlService) AddUserTags(tags model.Tags) error {
//	err := us.db.AutoMigrate(&model.Tags{})
//	tx := us.db.Begin()
//	if err = tx.Table("tags").Create(&tags).Error;err != nil {
//		tx.Rollback()
//        return err
//	}
//	tx.Commit()
//	return nil
//}
//
//// GetUsersByTags 根据查询条件返回用户列表
//func (us *MysqlService) GetUsersByTags(tags []string) ([]model.Tags, error) {
//
//	var res []model.Tags
//	tx :=us.db.Table("tags").Where("tags = ?", tags).Find(&res)
//	if tx.Error!= nil {
//        return nil, tx.Error
//    }
//	return res, nil
//}
//
//// UpdateUserTags 更新用户标签
//func (us *MysqlService) UpdateUserTags(tags model.Tags) error {
//	err := us.db.AutoMigrate(&model.Tags{})
//	if err != nil {
//		return err
//	}
//	tx := us.db.Begin()
//	res := tx.Table("tags").Where("id = ?", tags.ID).Update("tags", tags)
//	if res.Error != nil {
//		tx.Rollback()
//		return res.Error
//	}
//	tx.Commit()
//	return nil
//}
//
//// DeleteUserTags 删除用户标签
//func (us *MysqlService) DeleteUserTags(tags []string) error {
//
//	err := us.db.AutoMigrate(&model.Tags{})
//	if err != nil {
//		return err
//	}
//	tx := us.db.Begin()
//	res := tx.Where("tags = ?", tags).Delete(&model.Tags{})
//	if res.Error != nil {
//		tx.Rollback()
//		return res.Error
//	}
//	tx.Commit()
//	return nil
//}
