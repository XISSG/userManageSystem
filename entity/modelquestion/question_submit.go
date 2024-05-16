package modelquestion

import (
	"encoding/json"
	"github.com/xissg/userManageSystem/common/constant"
	"github.com/xissg/userManageSystem/utils"
	"time"
)

type QuestionSubmit struct {
	ID string `json:"id" gorm:"column id;type varchar(256); primaryKey"`
	// "编程语言"
	Language string `json:"language" gorm:"column language; type: varchar(128)"`
	//"用户代码"
	Code string `json:"code" gorm:"column code; type: text; not null"`
	//"判题信息json对象(包含上面的枚举值)
	JudgeInfo string `json:"judge_info" gorm:"column judge_info; type: text;"`
	//"判题状态（0-待判题,1-判题中,2-成功,3-失败)",
	Status int `json:"status" gorm:"column status; type: int; default: 0; not null"`
	//"判题id"
	QuestionId string `json:"question_id" gorm:"index; column question_id; type: varchar(256); not null"`
	//"创建用户id"
	UserId string `json:"user_id" gorm:"index; column user_id; type: varchar(256); not null"`
	//"创建时间"
	CreateTime time.Time `json:"create_time" gorm:"column create_time; type: datetime; not null"`
	//"更新时间"
	UpdateTime time.Time `json:"update_time" gorm:"column update_time; type: datetime; not null"`
	//"是否删除",
	IsDelete int8 `json:"is_delete" gorm:"column is_delete; type: int; default: 0; not null"`
}

func (qs *QuestionSubmit) TableName() string {
	return "question_submit"
}

type AddQuestionSubmitRequest struct {
	// "编程语言"
	Language string `json:"language"`
	//"用户代码"
	Code string `json:"code"`
	//"题目id"
	QuestionId string `json:"question_id"`
}

func AddQSToQS(request AddQuestionSubmitRequest) QuestionSubmit {
	var questionSubmit QuestionSubmit

	questionSubmit.ID = utils.NewUuid()
	questionSubmit.CreateTime = time.Now()
	questionSubmit.UpdateTime = time.Now()
	questionSubmit.IsDelete = constant.ALIVE
	questionSubmit.Status = WAITING
	questionSubmit.Language = request.Language
	questionSubmit.Code = request.Code
	questionSubmit.QuestionId = request.QuestionId

	return questionSubmit
}

type QueryQuestionSubmitRequest struct {
	// "编程语言"
	Language string `json:"language"`
	//"判题状态（0-待判题,1-判题中,2-成功,3-失败)",
	Status int `json:"status"`
	//"判题id"
	QuestionId string `json:"question_id"`
	//"创建用户id"
	UserId string `json:"user_id"`
	//"答案"
	Answer string `json:"answer"`
}

type CommonQueryQS struct {
	// "编程语言"
	Language string `json:"language"`
	//"判题状态（1-待判题,2-判题中,3-成功,4-失败)",
	Status int `json:"status" `
	//"判题id"
	QuestionId string `json:"question_id"`
	//"创建用户id"
	UserId string `json:"user_id"`
	//是否删除
	IsDelete int8 `json:"is_delete"`
}

func QueryQSToCommonQueryQS(qsQuery QueryQuestionSubmitRequest) CommonQueryQS {
	return CommonQueryQS{
		Language:   qsQuery.Language,
		Status:     qsQuery.Status,
		QuestionId: qsQuery.QuestionId,
		UserId:     qsQuery.UserId,
		IsDelete:   constant.ALIVE,
	}
}

type UpdateQuestionSubmitRequest struct {
	ID string `json:"id" `
	//"判题信息json对象(包含上面的枚举值)
	JudgeInfo []JudgeInfo `json:"judge_info"`
	//"判题状态（0-待判题,1-判题中,2-成功,3-失败)",
	Status int `json:"status"`
}

type CommonQuestionSubmitRequest struct {
	ID string `json:"id"`
	//"判题信息json对象(包含上面的枚举值)
	JudgeInfo string `json:"judge_info"`
	//"判题状态（0-待判题,1-判题中,2-成功,3-失败)",
	Status int `json:"status"`
}

func UpdateQSToCommonQS(request UpdateQuestionSubmitRequest) CommonQuestionSubmitRequest {
	var questionSubmit CommonQuestionSubmitRequest
	var judgeInfo string
	if request.JudgeInfo != nil {
		res, err := json.Marshal(request.JudgeInfo)
		if err != nil {
			return CommonQuestionSubmitRequest{}
		}
		judgeInfo = string(res)
	}

	questionSubmit.ID = request.ID
	questionSubmit.Status = request.Status
	questionSubmit.JudgeInfo = judgeInfo

	return questionSubmit
}

type ReturnQS struct {
	//题目id
	QuestionId string `json:"question_id"`
	// "编程语言"
	Language string `json:"language" `
	//"判题信息json对象(包含上面的枚举值)
	JudgeInfo JudgeInfo `json:"judge_info" `
	//"判题状态（0-待判题,1-判题中,2-成功,3-失败)",
	Status int `json:"status"`
	//"答案"
	Answer string `json:"answer"`
}

func QSToReturnQS(questionSubmit QuestionSubmit, Answer string) ReturnQS {
	var qsReturn ReturnQS
	var judgeInfo JudgeInfo

	if questionSubmit.JudgeInfo != "" {
		err := json.Unmarshal([]byte(questionSubmit.JudgeInfo), &judgeInfo)
		if err != nil {
			return ReturnQS{}
		}
	}

	qsReturn.QuestionId = questionSubmit.QuestionId
	qsReturn.Language = questionSubmit.Language
	qsReturn.Status = questionSubmit.Status
	qsReturn.JudgeInfo = judgeInfo
	qsReturn.Answer = Answer

	return qsReturn
}

func QSsToReturnQSs(questionSubmit []QuestionSubmit, answer string) []ReturnQS {
	var qsReturn []ReturnQS
	for _, v := range questionSubmit {
		qsReturn = append(qsReturn, QSToReturnQS(v, answer))
	}
	return qsReturn
}
