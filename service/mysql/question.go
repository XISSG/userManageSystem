package mysql

import (
	"github.com/xissg/userManageSystem/common/constant"
	"github.com/xissg/userManageSystem/entity/model_question"
	"gorm.io/gorm"
)

type QuestionService struct {
	db *gorm.DB
}

func NewQuestionMysqlService() *QuestionService {
	db := initDB()
	return &QuestionService{
		db: db,
	}
}

/**
 * @Description: 添加题目
 * @param q model_question.Question
 * @return error
 * @author xissg
 */
func (qds *QuestionService) AddQuestion(q model_question.Question) error {
	err := qds.db.AutoMigrate(&model_question.Question{})
	if err != nil {
		return err
	}
	tx := qds.db.Begin()
	err = tx.Table("question").Create(&q).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

/**
 * @Description: 更新题目
 * @param q model_question.Question
 * @return error
 * @author xissg
 */
func (qds *QuestionService) UpdateQuestion(q model_question.Question) error {
	err := qds.db.AutoMigrate(&model_question.Question{})
	if err != nil {
		return err
	}

	tx := qds.db.Begin()
	res := tx.Table("question").Where("id = ? AND is_delete = ?", q.ID, constant.ALIVE).Updates(q)
	if res.Error != nil {
		tx.Rollback()

		return res.Error
	}

	tx.Commit()

	return nil
}

/**
 * @Description: 删除题目
 * @param questionId string
 * @return error
 * @author xissg
 */
func (qds *QuestionService) DeleteQuestion(questionId string) error {
	err := qds.db.AutoMigrate(&model_question.Question{})
	if err != nil {
		return err
	}

	tx := qds.db.Begin()
	res := tx.Table("question").Where("id = ? AND is_delete = ?", questionId, constant.ALIVE).Update("is_delete", constant.DELETE)
	if res.Error != nil {
		tx.Rollback()

		return res.Error
	}

	tx.Commit()

	return nil
}

/**
 * @Description: 查询题目
 * @param questionId string
 * @return model_question.Question
 * @return error
 * @author xissg
 */
func (qds *QuestionService) GetQuestion(questionId string) (model_question.Question, error) {
	err := qds.db.AutoMigrate(&model_question.Question{})
	if err != nil {
		return model_question.Question{}, err
	}

	var res model_question.Question
	err = qds.db.Table("question").Where("id = ? AND is_delete = ?", questionId, constant.ALIVE).First(&res).Error
	if err != nil {

		return model_question.Question{}, err
	}

	return res, nil
}

/**
 * @Description: 查询题目列表
 * @param questionList model_question.CommonQueryQuestion
 * @return []model_question.Question
 * @return error
 * @author xissg
 */
func (qds *QuestionService) GetQuestionList(questionList model_question.CommonQueryQuestion, page, pageSize int) ([]model_question.Question, error) {
	offset := (page - 1) * pageSize
	err := qds.db.AutoMigrate(&model_question.Question{})
	if err != nil {
		return nil, err
	}

	var res []model_question.Question
	tx := qds.db.Begin()
	err = tx.Table("question").Where(&questionList).Limit(pageSize).Offset(offset).Find(&res).Error
	if err != nil {
		tx.Rollback()

		return nil, err
	}

	tx.Commit()

	return res, nil
}
