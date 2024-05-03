package service

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/model"
	"time"
)

type SessionServiceImpl struct {
	store redis.Store
}

func NewSessionService(store redis.Store) *SessionServiceImpl {
	return &SessionServiceImpl{
		store: store,
	}
}

func (us *SessionServiceImpl) NewOrUpdateSession(c *gin.Context, user model.UserSession) error {
	session := sessions.Default(c)

	maxAge := int(time.Now().Add(time.Hour*24).UTC().Unix() - time.Now().UTC().Unix())
	opts := sessions.Options{MaxAge: maxAge}
	session.Options(opts)
	session.Set("user", user)
	err := session.Save()
	if err != nil {
		return err
	}
	return nil
}

// GetSession 获取session
func (us *SessionServiceImpl) GetSession(c *gin.Context) model.UserSession {

	session := sessions.Default(c)
	sessionInfo := session.Get("user")
	if sessionInfo == nil {
		return model.UserSession{}
	}
	return sessionInfo.(model.UserSession)

}

// DeleteSession 删除session
func (us *SessionServiceImpl) DeleteSession(c *gin.Context) error {

	session := sessions.Default(c)
	session.Delete("user")
	err := session.Save()
	if err != nil {
		return err
	}
	return nil
}
