package redis

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/entity/model_user"
	"time"
)

type SessionService struct {
	store redis.Store
}

func NewSessionService(store redis.Store) *SessionService {

	return &SessionService{
		store: store,
	}
}

func (us *SessionService) NewOrUpdateSession(c *gin.Context, user model_user.UserSession) error {
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
func (us *SessionService) GetSession(c *gin.Context) (model_user.UserSession, error) {

	session := sessions.Default(c)
	sessionInfo := session.Get("user")
	if sessionInfo == nil {
		return model_user.UserSession{}, errors.New("session not found")
	}

	return sessionInfo.(model_user.UserSession), nil
}

// DeleteSession 删除session
func (us *SessionService) DeleteSession(c *gin.Context) error {

	session := sessions.Default(c)
	session.Delete("user")
	err := session.Save()
	if err != nil {
		return err
	}

	return nil
}
