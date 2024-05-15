package redis

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xissg/userManageSystem/entity/modelquestion"
	"github.com/xissg/userManageSystem/utils"
)

type QuestionSubmitService struct {
	rdb *redis.Client
}

func NewQuestionSubmitRedisService(rdb *redis.Client) *QuestionSubmitService {
	return &QuestionSubmitService{
		rdb: rdb,
	}
}
func (qsc *QuestionSubmitService) AddSubmitQuestion(submitQuestion modelquestion.QuestionSubmit, ctx *gin.Context) error {
	qsc.rdb.Set(ctx, submitQuestion.ID, submitQuestion, utils.RandomExpireTime())
	return nil
}

func (qsc *QuestionSubmitService) QuerySubmitQuestion(submitId string, ctx *gin.Context) (modelquestion.QuestionSubmit, error) {
	qsc.rdb.Get(ctx, submitId)
	return modelquestion.QuestionSubmit{}, nil
}
