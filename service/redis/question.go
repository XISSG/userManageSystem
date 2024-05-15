package redis

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xissg/userManageSystem/entity/modelquestion"
)

type QuestionService struct {
	rdb *redis.Client
}

func NewQuestionRedisService(rdb *redis.Client) *QuestionService {
	return &QuestionService{
		rdb: rdb,
	}
}

func (qcs *QuestionService) AddQuestion(q modelquestion.Question, ctx *gin.Context) error {
	return nil
}
func (qcs *QuestionService) UpdateQuestion(q modelquestion.Question, ctx *gin.Context) error {
	return nil
}
func (qcs *QuestionService) DeleteQuestion(questionId string, ctx *gin.Context) error {
	return nil
}
func (qcs *QuestionService) QueryQuestion(questionId string, ctx *gin.Context) (modelquestion.Question, error) {
	return modelquestion.Question{}, nil
}
