package mysql

import (
	"github.com/xissg/userManageSystem/common/constant"
	"github.com/xissg/userManageSystem/entity/modelquestion"
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
 * @param submitQuestion modelquestion.QuestionSubmit
 * @return error
 * @author xissg
 */
func (qsds *QuestionSubmitService) AddSubmitQuestion(submitQuestion modelquestion.QuestionSubmit) error {
	err := qsds.db.AutoMigrate(&modelquestion.QuestionSubmit{})
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
 * @return modelquestion.QuestionSubmit
 * @return error
 * @author xissg
 */
func (qsds *QuestionSubmitService) GetSubmitQuestion(submitId string) (modelquestion.QuestionSubmit, error) {
	err := qsds.db.AutoMigrate(&modelquestion.QuestionSubmit{})
	if err != nil {
		return modelquestion.QuestionSubmit{}, err
	}
	var res modelquestion.QuestionSubmit
	tx := qsds.db.Begin()
	err = tx.Table("question_submit").Where("id AND is_delete", submitId, constant.ALIVE).First(&res).Error
	if err != nil {
		tx.Rollback()

		return modelquestion.QuestionSubmit{}, err
	}

	tx.Commit()

	return res, nil
}

/**
 * @Description: 查询题目提交信息
 * @param qsQuery modelquestion.CommonQueryQS
 * @return modelquestion.QuestionSubmit
 * @return error
 * @author xissg
 */
func (qsds *QuestionSubmitService) GetSubmitQuestionList(qsQuery modelquestion.CommonQueryQS) ([]modelquestion.QuestionSubmit, error) {
	//TODO:使用分页查询
	err := qsds.db.AutoMigrate(&modelquestion.QuestionSubmit{})
	if err != nil {
		return nil, err
	}

	var res []modelquestion.QuestionSubmit

	err = qsds.db.Table("question_submit").Where(&qsQuery).Find(&res).Limit(1000).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (qsds *QuestionSubmitService) UpdateSubmitQuestion(request modelquestion.CommonQuestionSubmitRequest) error {
	err := qsds.db.AutoMigrate(&modelquestion.QuestionSubmit{})
	if err != nil {
		return err
	}
	tx := qsds.db.Begin()
	err = tx.Table("question_submit").Where("id", request.ID).Updates(request).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
