package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/constant"
	"github.com/xissg/userManageSystem/core/judge"
	"github.com/xissg/userManageSystem/entity/modelquestion"
	"github.com/xissg/userManageSystem/service/mysql"
	"github.com/xissg/userManageSystem/service/redis"
	"github.com/xissg/userManageSystem/utils"
	"log"
	"net/http"
	"sync"
)

type QuestionSubmitController struct {
	qsService       *mysql.QuestionSubmitService
	questionService *mysql.QuestionService
	sessionService  *redis.SessionService
}

func NewQuestionSubmitController(qsService *mysql.QuestionSubmitService, questionService *mysql.QuestionService, sessionService *redis.SessionService) *QuestionSubmitController {
	return &QuestionSubmitController{
		qsService:       qsService,
		questionService: questionService,
		sessionService:  sessionService,
	}
}

// Submit 提交代码
//
//	@Summary		submit
//	@Description	submit
//	@Tags			QuestionSubmit
//	@Accept			json
//	@Produce		json
//	@Param			modelquestion	body		modelquestion.AddQuestionSubmitRequest		true	"submit code"
//	@Success		200			{object}	utils.ApiResponse{data=string}	"submit success"
//	@Failure		404			{object}	utils.ApiResponse{data=nil}		"submit failed"
//	@Router			/api/submit/submission    [post]
func (qsc *QuestionSubmitController) Submit(c *gin.Context) {
	//用户身份校验
	session, err := qsc.sessionService.GetSession(c)
	if session.UserRole != constant.Common && session.UserRole != constant.Admin {
		log.Printf("you are not login")
		c.JSON(http.StatusOK, utils.NewResponse(nil, "you are not login").AuthERR())

		return
	}

	var qsAdd modelquestion.AddQuestionSubmitRequest
	//取出数据
	if err = c.ShouldBindJSON(&qsAdd); err != nil {
		log.Printf("Failed to unmarshal")
		c.JSON(http.StatusOK, utils.NewResponse(nil, "unmarshal error ").OperationERR())

		return
	}

	//校验编程语言是否合法
	language := checkLanguage(qsAdd.Language)
	if language == "" || qsAdd.QuestionId == "" || qsAdd.Code == "" {
		log.Printf("invalid language ")
		c.JSON(http.StatusOK, utils.NewResponse(nil, "invalid language ").ParamsERR())

		return
	}

	//转换成数据中的存储类型
	questionSubmit := modelquestion.AddQSToQS(qsAdd)
	//获取用户id
	question, _ := qsc.questionService.GetQuestion(qsAdd.QuestionId)
	questionSubmit.UserId = question.UserId

	err = qsc.qsService.AddSubmitQuestion(questionSubmit)
	if err != nil {
		log.Printf("Failed to submit")
		c.JSON(http.StatusOK, utils.NewResponse(nil, "submit error ").OperationERR())

		return
	}

	log.Printf("submit success")
	c.JSON(http.StatusOK, utils.NewResponse(questionSubmit.ID, "submit success").Success())

	//协程池优化性能
	var wg sync.WaitGroup

	wg.Add(1)
	go func(string) {
		defer wg.Done()
		judge.NewJudgeService(qsc.questionService, qsc.qsService).Judge(questionSubmit.ID)
	}(questionSubmit.ID)
	wg.Wait()
}

// GetQuestionSubmit
//
//	@Summary		get question submit result
//	@Description	get question submit result
//	@Tags			QuestionSubmit
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						true	"query"
//	@Success		200		{object}	utils.ApiResponse{data=modelquestion.ReturnQS}	"query  successful"
//	@Failure		404		{object}	utils.ApiResponse{data=nil}							"query failed"
//	@Router			/api/submit/query [get]
func (qsc *QuestionSubmitController) GetQuestionSubmit(c *gin.Context) {
	//用户身份校验
	session, err := qsc.sessionService.GetSession(c)
	if session.UserRole != constant.Common && session.UserRole != constant.Admin {
		log.Printf("you are not login")
		c.JSON(http.StatusOK, utils.NewResponse(nil, "you are not login").AuthERR())

		return
	}

	id := c.Param("id")
	if id == "" {
		log.Printf("invalid id")
		c.JSON(http.StatusOK, utils.NewResponse(nil, "invalid id ").ParamsERR())

		return
	}
	submit, err := qsc.qsService.GetSubmitQuestion(id)
	question, err := qsc.questionService.GetQuestion(submit.QuestionId)
	if err != nil {
		log.Printf("Failed to get submit result")
		c.JSON(http.StatusOK, utils.NewResponse(nil, "get submit result error ").OperationERR())
		return
	}

	if session.UserRole != constant.Admin && question.UserId != submit.UserId {
		result := modelquestion.QSToReturnQS(submit, "")
		log.Printf("get submit result success")
		c.JSON(http.StatusOK, utils.NewResponse(result, "get submit result success").Success())
		return
	}

	result := modelquestion.QSToReturnQS(submit, question.Answer)
	log.Printf("get submit result success")
	c.JSON(http.StatusOK, utils.NewResponse(result, "get submit result success").Success())
	return
}

// GetQuestionSubmitList 获取代码结果
//
//	@Summary		get question submit result
//	@Description	get question submit result
//	@Tags			QuestionSubmit
//	@Accept			json
//	@Produce		json
//	@Param			query	body		modelquestion.QueryQuestionSubmitRequest						true	"queries"
//	@Success		200		{object}	utils.ApiResponse{data=[]modelquestion.ReturnQS}	"query  successful"
//	@Failure		404		{object}	utils.ApiResponse{data=nil}							"query failed"
//	@Router			/api/submit/query [post]
func (qsc *QuestionSubmitController) GetQuestionSubmitList(c *gin.Context) {
	//用户身份校验
	session, err := qsc.sessionService.GetSession(c)
	if err != nil {
		log.Printf("you are not login")
		c.JSON(http.StatusOK, utils.NewResponse(nil, "you are not login").AuthERR())

		return
	}
	if session.UserRole != constant.Admin {
		log.Printf("you are not admin")
		c.JSON(http.StatusOK, utils.NewResponse(nil, "you are not admin").AuthERR())

		return
	}

	var qsQuery modelquestion.QueryQuestionSubmitRequest
	//取出数据
	if err = c.ShouldBindJSON(&qsQuery); err != nil {
		log.Printf("Failed to unmarshal")
		c.JSON(http.StatusOK, utils.NewResponse(nil, "unmarshal error ").OperationERR())

		return
	}

	//校验数据
	err = qsc.checkQueries(qsQuery)
	if err != nil {
		log.Printf("invalid queries %v", err)
		c.JSON(http.StatusOK, utils.NewResponse(nil, err.Error()).ParamsERR())

		return
	}
	qs := modelquestion.QueryQSToCommonQueryQS(qsQuery)
	res, err := qsc.qsService.GetSubmitQuestionList(qs)
	if err != nil {
		log.Printf("Failed to get submit result")
		c.JSON(http.StatusOK, utils.NewResponse(nil, "get submit result error ").OperationERR())
		return
	}

	qsReturns := modelquestion.QSsToReturnQSs(res, "")
	log.Printf("get submit result success")
	c.JSON(http.StatusOK, utils.NewResponse(qsReturns, "get submit result success").Success())
}

func checkLanguage(lang string) string {
	switch lang {
	case modelquestion.C, modelquestion.Cpp, modelquestion.Java, modelquestion.Go, modelquestion.Python, modelquestion.Js, modelquestion.Ts, modelquestion.Php:
		return lang
	default:
		return ""
	}
}

func (qsc *QuestionSubmitController) checkQueries(qsQuery modelquestion.QueryQuestionSubmitRequest) error {
	if qsQuery.QuestionId != "" && len(qsQuery.QuestionId) > 256 {
		return errors.New("invalid query question id")
	}

	if qsQuery.Answer != "" && len(qsQuery.Answer) > 8192 {
		return errors.New("invalid queries answer")
	}
	if qsQuery.UserId != "" && len(qsQuery.UserId) > 256 {
		return errors.New("invalid query user id")
	}
	if qsQuery.Language != "" && checkLanguage(qsQuery.Language) == "" {
		return errors.New("invalid query language")
	}
	return nil
}
