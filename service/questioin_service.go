package service

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/entity/modelquestion"
	"github.com/xissg/userManageSystem/service/mysql"
)

type QuestionService struct {
	mysql.QuestionService
}

func NewQuestionService(questionDBService mysql.QuestionService) *QuestionService {
	return &QuestionService{
		QuestionService: questionDBService,
	}
}

/**
 * 创建题目
 * @param  modelquestion entity.Question
 * @param  c *gin.Context
 * @return error
 * @author xissg
 */
func (qs *QuestionService) CreateQuestion(question modelquestion.Question, c *gin.Context) error {

	err := qs.QuestionService.AddQuestion(question)
	if err != nil {
		return err
	}

	return nil
}

/**
 * 更新题目
 * @param  modelquestion entity.Question
 * @param  c *gin.Context
 * @return error
 * @author xissg
 */
func (qs *QuestionService) UpdateQuestion(question modelquestion.Question) error {
	err := qs.QuestionService.UpdateQuestion(question)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 删除题目
 * @param  questionId string
 * @param  c *gin.Context
 * @return error
 */
func (qs *QuestionService) DeleteQuestion(questionId string) error {
	err := qs.QuestionService.DeleteQuestion(questionId)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 查询题目
 * @param  questionId string
 * @param  c *gin.Context
 * @return entity.Question, error
 * @author xissg
 * @return entity.Question, error
 */
func (qs *QuestionService) QueryQuestion(questionId string) (modelquestion.Question, error) {
	question, err := qs.QuestionService.GetQuestion(questionId)
	if err != nil {
		return modelquestion.Question{}, err
	}
	return question, nil
}
