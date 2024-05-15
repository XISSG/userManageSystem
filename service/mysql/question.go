package mysql

import (
	"github.com/xissg/userManageSystem/constant"
	"github.com/xissg/userManageSystem/entity/modelquestion"
	"gorm.io/gorm"
)

type QuestionService struct {
	db *gorm.DB
}

func NewQuestionMysqlService(db *gorm.DB) *QuestionService {
	return &QuestionService{
		db: db,
	}
}

/**
 * @Description: 添加题目
 * @param q modelquestion.Question
 * @return error
 * @author xissg
 */
func (qds *QuestionService) AddQuestion(q modelquestion.Question) error {
	err := qds.db.AutoMigrate(&modelquestion.Question{})
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
 * @param q modelquestion.Question
 * @return error
 * @author xissg
 */
func (qds *QuestionService) UpdateQuestion(q modelquestion.Question) error {
	err := qds.db.AutoMigrate(&modelquestion.Question{})
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
	err := qds.db.AutoMigrate(&modelquestion.Question{})
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
 * @return modelquestion.Question
 * @return error
 * @author xissg
 */
func (qds *QuestionService) GetQuestion(questionId string) (modelquestion.Question, error) {
	err := qds.db.AutoMigrate(&modelquestion.Question{})
	if err != nil {
		return modelquestion.Question{}, err
	}

	var res modelquestion.Question
	tx := qds.db.Begin()
	err = tx.Table("question").Where("id = ? AND is_delete = ?", questionId, constant.ALIVE).First(&res).Error
	if err != nil {
		tx.Rollback()

		return modelquestion.Question{}, err
	}

	tx.Commit()

	return res, nil
}

/**
 * @Description: 查询题目列表
 * @param questionList modelquestion.CommonQueryQuestion
 * @return []modelquestion.Question
 * @return error
 * @author xissg
 */
func (qds *QuestionService) GetQuestionList(questionList modelquestion.CommonQueryQuestion) ([]modelquestion.Question, error) {
	err := qds.db.AutoMigrate(&modelquestion.Question{})
	if err != nil {
		return nil, err
	}

	var res []modelquestion.Question
	tx := qds.db.Begin()
	err = tx.Table("question").Where(&questionList).Find(&res).Error
	if err != nil {
		tx.Rollback()

		return nil, err
	}

	tx.Commit()

	return res, nil
}
