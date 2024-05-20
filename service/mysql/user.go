package mysql

import (
	"github.com/xissg/userManageSystem/common/constant"
	"github.com/xissg/userManageSystem/entity/model_user"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService() *UserService {
	db := initDB()
	return &UserService{
		db: db,
	}
}

/**
 * @Description: 新增用户
 * @param user model_user.User
 * @return error
 * @author xissg
 */
func (us *UserService) AddUser(user model_user.User) error {
	err := us.db.AutoMigrate(&model_user.User{})
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
 * @return model_user.User
 * @return error
 * @author xissg
 */
func (us *UserService) GetUser(accountName string) (model_user.User, error) {
	_ = us.db.AutoMigrate(&model_user.User{})
	var res model_user.User
	tx := us.db.Table("user").Where("user_account = ? AND is_delete = ?", accountName, constant.ALIVE).First(&res)

	return res, tx.Error
}

/**
 * @Description: 获取用户列表
 * @param queryModel model_user.AdminUserQueryRequest
 * @return []model_user.User
 * @return error
 * @author xissg
 */
func (us *UserService) GetUserList(queryModel model_user.AdminUserQueryRequest, page, pageSize int) ([]model_user.User, error) {
	var users []model_user.User
	offset := (page - 1) * pageSize
	err := us.db.AutoMigrate(&model_user.User{})
	if err != nil {
		return nil, err
	}

	err = us.db.Table("user").Where(&queryModel).Limit(pageSize).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

/**
 * @Description: 更新用户信息
 * @param model_user.User
 * @return error
 * @author xissg
 */
func (us *UserService) UpdateUser(user model_user.User) error {
	err := us.db.AutoMigrate(&model_user.User{})
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
	err := us.db.AutoMigrate(&model_user.User{})
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
