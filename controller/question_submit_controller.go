package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/common/api_response"
	"github.com/xissg/userManageSystem/common/constant"
	"github.com/xissg/userManageSystem/core/judge"
	"github.com/xissg/userManageSystem/entity/model_question"
	"github.com/xissg/userManageSystem/service/mysql"
	"github.com/xissg/userManageSystem/service/redis"
	"log"
	"net/http"
	"strings"
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
//	@Summary		Submit
//	@Description	Submit
//	@Tags			QuestionSubmit
//	@Accept			json
//	@Produce		json
//	@Param			model_question	body		model_question.AddQuestionSubmitRequest		true	"Submit code"
//	@Success		200			{object}	api_response.ApiResponse{data=string}	"Submit success"
//	@Failure		400			{object}	api_response.ApiResponse{data=nil}		"Submit failed"
//	@Router			/api/submit/add    [post]
func (qsc *QuestionSubmitController) Submit(c *gin.Context) {
	//用户身份校验
	session, err := qsc.sessionService.GetSession(c)
	if session.UserRole != constant.Common && session.UserRole != constant.Admin {
		log.Printf("you are not login")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "you are not login").Response(api_response.AUTHERR))

		return
	}

	var qsAdd model_question.AddQuestionSubmitRequest
	//取出数据
	if err = c.ShouldBindJSON(&qsAdd); err != nil {
		log.Printf("Failed to unmarshal")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "unmarshal error ").Response(api_response.OPERATIONERR))

		return
	}

	//校验编程语言是否合法
	language := checkLanguage(qsAdd.Language)
	if language == "" || qsAdd.QuestionId == "" || qsAdd.Code == "" {
		log.Printf("invalid language ")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "invalid language ").Response(api_response.PARAMSERR))

		return
	}

	//转换成数据中的存储类型
	questionSubmit := model_question.AddQSToQS(qsAdd)
	//获取用户id
	question, _ := qsc.questionService.GetQuestion(qsAdd.QuestionId)
	questionSubmit.UserId = question.UserId

	err = qsc.qsService.AddSubmitQuestion(questionSubmit)
	if err != nil {
		log.Printf("Failed to submit")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "submit error ").Response(api_response.OPERATIONERR))

		return
	}

	log.Printf("submit success")
	c.JSON(http.StatusOK, api_response.NewResponse(questionSubmit.ID, "submit success").Response(api_response.SUCCESS))

	//使用消息队列发送信息
	var wg sync.WaitGroup
	wg.Add(1)
	go func(string) {
		judging := judge.NewJudgeService(qsc.questionService, qsc.qsService)
		judging.Judge(questionSubmit.ID)
		wg.Done()
	}(questionSubmit.ID)

	wg.Wait()
}

// GetQuestionSubmit
//
//	@Summary		Get question submit result
//	@Description	Get question submit result
//	@Tags			QuestionSubmit
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						true	"Query id"
//	@Success		200		{object}	api_response.ApiResponse{data=model_question.ReturnQS}	"Query  success"
//	@Failure		400		{object}	api_response.ApiResponse{data=nil}							"Query fail"
//	@Router			/api/submit/query [get]
func (qsc *QuestionSubmitController) GetQuestionSubmit(c *gin.Context) {
	//用户身份校验
	session, err := qsc.sessionService.GetSession(c)
	if session.UserRole != constant.Common && session.UserRole != constant.Admin {
		log.Printf("you are not login")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "you are not login").Response(api_response.AUTHERR))

		return
	}

	id := c.Param("id")
	if id == "" {
		log.Printf("invalid id")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "invalid id ").Response(api_response.PARAMSERR))

		return
	}
	submit, err := qsc.qsService.GetSubmitQuestion(id)
	question, err := qsc.questionService.GetQuestion(submit.QuestionId)
	if err != nil {
		log.Printf("Failed to get submit result")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "get submit result error ").Response(api_response.OPERATIONERR))
		return
	}

	if session.UserRole != constant.Admin && question.UserId != submit.UserId {
		result := model_question.QSToReturnQS(submit, "")
		log.Printf("get submit result success")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(result, "get submit result success").Response(api_response.SUCCESS))
		return
	}

	result := model_question.QSToReturnQS(submit, question.Answer)
	log.Printf("get submit result success")
	c.JSON(http.StatusOK, api_response.NewResponse(result, "get submit result success").Response(api_response.SUCCESS))
	return
}

// GetQuestionSubmitList 获取代码结果
//
//	@Summary		Get question submit list
//	@Description	Get question submit list
//	@Tags			QuestionSubmit
//	@Accept			json
//	@Produce		json
//	@Param			query	body		model_question.QueryQuestionSubmitRequest						true	"Query condition"
//	@Success		200		{object}	api_response.ApiResponse{data=[]model_question.ReturnQS}	"Query  success"
//	@Failure		400		{object}	api_response.ApiResponse{data=nil}							"Query fail"
//	@Router			/api/submit/query [post]
func (qsc *QuestionSubmitController) GetQuestionSubmitList(c *gin.Context) {
	//用户身份校验
	session, err := qsc.sessionService.GetSession(c)
	if err != nil {
		log.Printf("you are not login")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "you are not login").Response(api_response.AUTHERR))

		return
	}
	if session.UserRole != constant.Admin {
		log.Printf("you are not admin")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "you are not admin").Response(api_response.AUTHERR))

		return
	}

	var qsQuery model_question.QueryQuestionSubmitRequest
	//取出数据
	if err = c.ShouldBindJSON(&qsQuery); err != nil {
		log.Printf("Failed to unmarshal")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "unmarshal error ").Response(api_response.OPERATIONERR))

		return
	}

	//校验数据
	err = qsc.checkQueries(qsQuery)
	if err != nil {
		log.Printf("invalid queries %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, err.Error()).Response(api_response.PARAMSERR))

		return
	}

	page := qsQuery.Page
	pageSize := qsQuery.PageSize

	qs := model_question.QueryQSToCommonQueryQS(qsQuery)
	res, err := qsc.qsService.GetSubmitQuestionList(qs, page, pageSize)
	if err != nil {
		log.Printf("Failed to get submit result")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "get submit result error ").Response(api_response.OPERATIONERR))
		return
	}

	qsReturns := model_question.QSsToReturnQSs(res, "")
	log.Printf("get submit result success")
	c.JSON(http.StatusOK, api_response.NewResponse(qsReturns, "get submit result success").Response(api_response.SUCCESS))
}

func checkLanguage(lang string) string {
	lan := strings.ToLower(lang)
	switch lan {
	case constant.C, constant.Cpp, constant.Java, constant.Go, constant.Python:
		return lan
	default:
		return ""
	}
}

func (qsc *QuestionSubmitController) checkQueries(qsQuery model_question.QueryQuestionSubmitRequest) error {
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
	if qsQuery.Page <= 0 {
		qsQuery.Page = 1
	}
	if qsQuery.PageSize <= 0 {
		qsQuery.PageSize = 10
	}
	return nil
}
