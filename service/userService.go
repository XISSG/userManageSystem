package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/xissg/userManageSystem/model"
	"gorm.io/gorm"
	"log"
)

const Delete = 0 //is_delete 字段，默认为不删除用户

type UserService struct {
	db    *gorm.DB
	store *sessions.CookieStore
}

func NewService(db *gorm.DB, store *sessions.CookieStore) *UserService {

	return &UserService{
		db:    db,
		store: store,
	}
}

// AddUser 添加用户
func (us *UserService) AddUser(user model.User) error {

	if err := us.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

// QueryUser 根据用户名返回查询用户
func (us *UserService) QueryUser(user model.User) (model.User, error) {

	var res model.User
	tx := us.db.Where("user_name = ? AND is_delete = ?", user.UserName, Delete).First(&res)

	return res, tx.Error

}

// QueryUsers 根据查询条件返回用户列表
func (us *UserService) QueryUsers(user model.User) ([]model.User, error) {

	var res []model.User
	usr := "%" + user.UserName + "%"
	tx := us.db.Where("user_name LIKE ? AND is_delete = ?", usr, Delete).Find(&res)

	return res, tx.Error

}

// LogicDeleteUser 将isDelete字段设为 1 ,不进行实际删除
func (us *UserService) LogicDeleteUser(user model.User) error {

	tx := us.db.Where("user_name = ? and is_delete = ?", user.UserName, Delete).Update("is_delete = ? ", 1)

	return tx.Error
}

// RealDeleteUser 真实删除用户
func (us *UserService) RealDeleteUser(user model.User) error {

	tx := us.db.Delete("user_name = ?", user.UserName)

	return tx.Error
}

func (us *UserService) NewSession(c *gin.Context, userSession model.UserSession) error {

	session, err := us.store.New(c.Request, "SessionID")
	if err != nil {
		log.Printf("session new error: %v", err)
		return err
	}
	session.Values["userSession"] = userSession
	err = session.Save(c.Request, c.Writer)
	return err
}

// GetSession 获取session
func (us *UserService) GetSession(c *gin.Context) model.UserSession {

	session, err := us.store.Get(c.Request, "SessionID")

	if err != nil {
		log.Printf("session get error: %v", err)
		return model.UserSession{}
	}
	if session.IsNew {
		log.Printf("session is new")
		return model.UserSession{}
	}

	userSession := session.Values["userSession"].(model.UserSession)
	return userSession

}

// DeleteSession 删除session
func (us *UserService) DeleteSession(c *gin.Context) error {

	session, _ := us.store.Get(c.Request, "SessionID")
	session.Options.MaxAge = -1
	err := session.Save(c.Request, c.Writer)
	return err
}
