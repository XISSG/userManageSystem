package service

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/entity/modelquestion"
	"github.com/xissg/userManageSystem/service/mysql"
)

type QuestionSubmitService struct {
	qsm mysql.QuestionSubmitService
}

func NewQuestionSubmitService(questionSubmitDBService mysql.QuestionSubmitService) *QuestionSubmitService {
	return &QuestionSubmitService{
		qsm: questionSubmitDBService,
	}
}

/**
 * @param submitQuestion
 * @param ctx
 * @return error
 * @Description: 创建题目提交信息
 * @Author: xissg
 */
func (qs *QuestionSubmitService) CreateSubmitQuestion(submitQuestion modelquestion.QuestionSubmit, c *gin.Context) error {

	err := qs.qsm.AddSubmitQuestion(submitQuestion)
	if err != nil {
		return err
	}

	return nil
}

/**
 * @param submitId
 * @param ctx
 * @return entity.QuestionSubmit, error
 * @Description: 查询题目提交信息
 * @Author: xissg
 */
func (qs *QuestionSubmitService) QuerySubmitQuestion(submitId string, c *gin.Context) (modelquestion.QuestionSubmit, error) {
	//TODO
	return modelquestion.QuestionSubmit{}, nil
}
