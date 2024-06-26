package model_question

import (
	"encoding/json"
	"github.com/xissg/userManageSystem/common/constant"
	"github.com/xissg/userManageSystem/utils"
	"time"
)

type Question struct {
	ID string `json:"id" gorm:"column id;type varchar(256); primaryKey"`
	// "标题"
	Title string `json:"title" gorm:"column title; type varchar(512)"`
	// "内容"
	Content string `json:"content" gorm:"column content; type text"`
	// "标签列表json数组"
	Tag string `json:"tags" gorm:"column tag; type varchar(1024)"`
	// "题目答案"
	Answer string `json:"answer" gorm:"column answer; type text"`
	// "题目提交数
	SubmitNum int `json:"submit_num" gorm:"column submit_num; type int; not null;default: 0"`
	// "题目通过数"
	AcceptNum int `json:"accept_num" gorm:"column accept_num; type int; not null;default: 0"`
	// "判题用例json数组"
	JudgeCase string `json:"judge_case" gorm:"column judge_case; type text"`
	// "判题配置json对象"
	JudgeConfig string `json:"judge_config" gorm:"column judge_config; type text"`
	// "点赞数"
	ThumNum int `json:"thum_num" gorm:"column thum_num; type int; not null;default: 0"`
	// "创建用户id"
	UserId string `json:"user_id" gorm:"index; column user_id;type varchar(256); not null"`
	// "创建时间"
	CreateTime time.Time `json:"create_time" gorm:"column create_time; type datetime"`
	// "更新时间"
	UpdateTime time.Time `json:"update_time" gorm:"column update_time; type datetime"`
	// "是否删除"
	IsDelete int8 `json:"is_delete" gorm:"column is_delete; type int; default: 0"`
}

func (q *Question) TableName() string {
	return "question"
}

type AddQuestionRequest struct {
	// "标题"
	Title string `json:"title"`
	// "内容"
	Content string `json:"content"`
	// "标签列表json数组"
	Tag string `json:"tag" `
	// "题目答案"
	Answer []string `json:"answer" `
	// "判题用例json数组"
	JudgeCase []JudgeCase `json:"judge_case" `
	// "判题配置json对象"
	JudgeConfig JudgeConfig `json:"judge_config" `
	// "创建用户id"
	UserId string `json:"user_id"`
}

func AddQuestionToQuestion(addQuestion AddQuestionRequest) Question {
	var question Question
	var judgeCase, judgeConfig, answer string
	var errs error
	if addQuestion.JudgeCase != nil {
		res, err := json.Marshal(addQuestion.JudgeCase)
		errs = err
		judgeCase = string(res)
	}
	if &addQuestion.JudgeConfig != nil {
		res, err := json.Marshal(addQuestion.JudgeConfig)
		errs = err
		judgeConfig = string(res)
	}
	if addQuestion.Answer != nil {
		res, err := json.Marshal(addQuestion.Answer)
		errs = err
		answer = string(res)
	}
	if errs != nil {
		return Question{}
	}

	question.ID = utils.NewUuid()
	question.Title = addQuestion.Title
	question.Content = addQuestion.Content
	question.AcceptNum = 0
	question.SubmitNum = 0
	question.Tag = addQuestion.Tag
	question.Answer = answer
	question.JudgeCase = judgeCase
	question.JudgeConfig = judgeConfig
	question.UserId = addQuestion.UserId
	question.ThumNum = 0
	question.CreateTime = time.Now()
	question.UpdateTime = time.Now()
	question.IsDelete = constant.ALIVE

	return question
}

type UpdateQuestionRequest struct {
	ID string `json:"id"`
	// "标题"
	Title string `json:"title" `
	// "内容"
	Content string `json:"content" `
	// "标签列表json数组"
	Tag string `json:"tag" `
	// "题目答案"
	Answer []string `json:"answer"`
	// "判题用例json数组"
	JudgeCase []JudgeCase `json:"judge_case" `
	// "判题配置json对象"
	JudgeConfig JudgeConfig `json:"judge_config"`
}

func UpdateQuestionToQuestion(old Question, updateQuestion UpdateQuestionRequest) Question {

	var errs error

	if updateQuestion.Title != "" {
		old.Title = updateQuestion.Title
	}
	if updateQuestion.Content != "" {
		old.Content = updateQuestion.Content
	}
	if updateQuestion.Answer != nil {
		answer, err := json.Marshal(updateQuestion.Answer)
		errs = err
		old.Answer = string(answer)
	}
	if updateQuestion.Tag != "" {

		old.Tag = updateQuestion.Tag
	}
	if updateQuestion.JudgeCase != nil {
		judgeCase, err := json.Marshal(updateQuestion.JudgeCase)
		errs = err
		old.JudgeCase = string(judgeCase)
	}
	if &updateQuestion.JudgeConfig != nil {
		judgeConfig, err := json.Marshal(updateQuestion.JudgeConfig)
		errs = err
		old.JudgeConfig = string(judgeConfig)
	}

	if errs != nil {
		return Question{}
	}
	return old
}

type QueryQuestionRequest struct {
	ID string `json:"id"`
	// "标题"
	Title string `json:"title"`
	// "内容"
	Content string `json:"content"`
	// "标签列表json数组"
	Tag string `json:"tag" `
	// "创建用户id"
	UserId string `json:"user_id" `

	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func QueryQuestionToQuestion(q QueryQuestionRequest) Question {

	return Question{
		ID:      q.ID,
		Title:   q.Title,
		Content: q.Content,
		Tag:     q.Tag,
		UserId:  q.UserId,
	}
}

type CommonQueryQuestion struct {
	ID string `json:"id"`
	// "标题"
	Title string `json:"title"`
	// "内容"
	Content string `json:"content"`
	// "标签列表json数组"
	Tag string `json:"tag" `
	// "创建用户id"
	UserId string `json:"user_id" `
	// "是否删除"
	IsDelete int8 `json:"is_delete"`
}

func QueryQToCommonQueryQ(queryQuestion QueryQuestionRequest) CommonQueryQuestion {

	return CommonQueryQuestion{
		ID:       queryQuestion.ID,
		Title:    queryQuestion.Title,
		Content:  queryQuestion.Content,
		Tag:      queryQuestion.Tag,
		UserId:   queryQuestion.UserId,
		IsDelete: constant.ALIVE,
	}
}

type ReturnQuestion struct {
	ID string `json:"id" `
	// "标题"
	Title string `json:"title"`
	// "内容"
	Content string `json:"content"`
	// "标签列表json数组"
	Tag string `json:"tag"`
	// "题目答案"
	Answer []string `json:"answer"`
	// "题目提交数
	SubmitNum int `json:"submit_num" `
	// "题目通过数"
	AcceptNum int `json:"accept_num"`
	// "判题配置json对象"
	JudgeConfig JudgeConfig `json:"judge_config" `
	// "点赞数"
	ThumNum int `json:"thum_num"`
	//用户id
	UserId string `json:"user_id"`
}

func QuestionToReturnQuestion(question Question) ReturnQuestion {
	var answer []string
	var judgeConfig JudgeConfig
	var errs error

	if question.JudgeConfig != "" {
		err := json.Unmarshal([]byte(question.JudgeConfig), &judgeConfig)
		errs = err
	}
	if question.Answer != "" {
		err := json.Unmarshal([]byte(question.Answer), &answer)
		errs = err
	}
	if errs != nil {
		return ReturnQuestion{}
	}
	return ReturnQuestion{
		ID:          question.ID,
		Title:       question.Title,
		Content:     question.Content,
		Answer:      answer,
		Tag:         question.Tag,
		SubmitNum:   question.SubmitNum,
		AcceptNum:   question.AcceptNum,
		JudgeConfig: judgeConfig,
		ThumNum:     question.ThumNum,
		UserId:      question.UserId,
	}
}

func QuestionsToReturnQuestions(questions []Question) []ReturnQuestion {
	var returnQuestions []ReturnQuestion
	for _, question := range questions {
		returnQuestions = append(returnQuestions, QuestionToReturnQuestion(question))
	}
	return returnQuestions
}
