package mysql

import (
	"github.com/xissg/userManageSystem/common/constant"
	"github.com/xissg/userManageSystem/entity/modeluser"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {

	return &UserService{
		db: db,
	}
}

/**
 * @Description: 新增用户
 * @param user modeluser.User
 * @return error
 * @author xissg
 */
func (us *UserService) AddUser(user modeluser.User) error {
	err := us.db.AutoMigrate(&modeluser.User{})
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

/**
 * @Description: 获取用户
 * @param accountName string
 * @return modeluser.User
 * @return error
 * @author xissg
 */
func (us *UserService) GetUser(accountName string) (modeluser.User, error) {
	_ = us.db.AutoMigrate(&modeluser.User{})
	var res modeluser.User
	tx := us.db.Table("user").Where("user_account = ? AND is_delete = ?", accountName, constant.ALIVE).First(&res)

	return res, tx.Error
}

/**
 * @Description: 获取用户列表
 * @param queryModel modeluser.AdminUserQueryRequest
 * @return []modeluser.User
 * @return error
 * @author xissg
 */
func (us *UserService) GetUserList(queryModel modeluser.AdminUserQueryRequest) ([]modeluser.User, error) {
	var users []modeluser.User
	err := us.db.AutoMigrate(&modeluser.User{})
	if err != nil {
		return nil, err
	}

	err = us.db.Table("user").Where(&queryModel).Find(&users).Limit(1000).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

/**
 * @Description: 更新用户信息
 * @param modeluser.User
 * @return error
 * @author xissg
 */
func (us *UserService) UpdateUser(user modeluser.User) error {
	err := us.db.AutoMigrate(&modeluser.User{})
	if err != nil {
		return err
	}

	tx := us.db.Begin()
	res := tx.Table("user").Where("user_account = ? AND is_delete = ?", user.UserAccount, constant.ALIVE).Updates(user)
	if res.Error != nil {
		tx.Rollback()

		return res.Error
	}

	tx.Commit()

	return nil
}

/**
 * @Description: 删除用户
 * @param accountName string
 * @return error
 * @author xissg
 */
func (us *UserService) DeleteUser(accountName string) error {
	err := us.db.AutoMigrate(&modeluser.User{})
	if err != nil {
		return err
	}

	tx := us.db.Begin()
	res := tx.Table("user").Where("user_account = ? AND is_delete = ?", accountName, constant.ALIVE).Update("is_delete", constant.DELETE)
	if res.Error != nil {
		tx.Rollback()

		return res.Error
	}

	tx.Commit()

	return nil
}
