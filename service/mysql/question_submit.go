package mysql

import (
	"github.com/xissg/userManageSystem/common/constant"
	"github.com/xissg/userManageSystem/entity/model_question"
	"gorm.io/gorm"
)

type QuestionSubmitService struct {
	db *gorm.DB
}

func NewQuestionSubmitMysqlService(db *gorm.DB) *QuestionSubmitService {
	return &QuestionSubmitService{
		db: db,
	}
}

/**
 * @Description: 添加题目提交信息
 * @param submitQuestion model_question.QuestionSubmit
 * @return error
 * @author xissg
 */
func (qsds *QuestionSubmitService) AddSubmitQuestion(submitQuestion model_question.QuestionSubmit) error {
	err := qsds.db.AutoMigrate(&model_question.QuestionSubmit{})
	if err != nil {
		return err
	}
	tx := qsds.db.Begin()
	err = tx.Table("question_submit").Create(&submitQuestion).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

/**
 * @Description: 查询题目提交信息
 * @param submitId string
 * @return model_question.QuestionSubmit
 * @return error
 * @author xissg
 */
func (qsds *QuestionSubmitService) GetSubmitQuestion(submitId string) (model_question.QuestionSubmit, error) {
	err := qsds.db.AutoMigrate(&model_question.QuestionSubmit{})
	if err != nil {
		return model_question.QuestionSubmit{}, err
	}
	var res model_question.QuestionSubmit
	err = qsds.db.Table("question_submit").Where("id = ? AND is_delete = ?", submitId, constant.ALIVE).First(&res).Error
	if err != nil {

		return model_question.QuestionSubmit{}, err
	}

	return res, nil
}

/**
 * @Description: 查询题目提交信息
 * @param qsQuery model_question.CommonQueryQS
 * @return model_question.QuestionSubmit
 * @return error
 * @author xissg
 */
func (qsds *QuestionSubmitService) GetSubmitQuestionList(qsQuery model_question.CommonQueryQS) ([]model_question.QuestionSubmit, error) {
	//TODO:使用分页查询
	err := qsds.db.AutoMigrate(&model_question.QuestionSubmit{})
	if err != nil {
		return nil, err
	}

	var res []model_question.QuestionSubmit

	err = qsds.db.Table("question_submit").Where(&qsQuery).Find(&res).Limit(1000).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (qsds *QuestionSubmitService) UpdateSubmitQuestion(request model_question.CommonQuestionSubmitRequest) error {
	err := qsds.db.AutoMigrate(&model_question.QuestionSubmit{})
	if err != nil {
		return err
	}
	tx := qsds.db.Begin()
	err = tx.Table("question_submit").Where("id = ?", request.ID).Updates(request).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
