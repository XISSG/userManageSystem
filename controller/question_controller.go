package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/common/api_response"
	"github.com/xissg/userManageSystem/common/constant"
	"github.com/xissg/userManageSystem/entity/model_question"
	mysql2 "github.com/xissg/userManageSystem/service/mysql"
	"github.com/xissg/userManageSystem/service/redis"
	"log"
	"net/http"
)

//题目的增删改查

type QuestionController struct {
	questionService *mysql2.QuestionService
	session         *redis.SessionService
}

func NewQuestionController(questionService *mysql2.QuestionService, session *redis.SessionService) *QuestionController {
	return &QuestionController{
		questionService: questionService,
		session:         session,
	}
}

// AddQuestion 添加题目
//
//	@Summary		Add question
//	@Description	Add question
//	@Tags			Question
//	@Accept			json
//	@Produce		json
//	@Param			question	body		model_question.AddQuestionRequest	true	"Add question"
//	@Success		200			{object}	api_response.ApiResponse{data=nil}			"Add question success"
//	@Failure		400			{object}	api_response.ApiResponse{data=nil}			"Add  question fail"
//	@Router			/api/question/admin/add    [post]
func (qc *QuestionController) AddQuestion(c *gin.Context) {

	session, _ := qc.session.GetSession(c)
	if session.UserRole != constant.Admin {
		log.Printf("you are not allowed to create a new question")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "you are not allowed to create a new question").Response(api_response.AUTHERR))

		return
	}

	var receiveQuestion model_question.AddQuestionRequest
	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&receiveQuestion); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "unmarshal error ").Response(api_response.OPERATIONERR))

		return
	}

	add := model_question.AddQuestionToQuestion(receiveQuestion)
	//数据校验
	err := qc.checkQuestion(add)
	if err != nil {
		log.Printf("validate %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "validate error").Response(api_response.OPERATIONERR))

		return
	}
	question := model_question.AddQuestionToQuestion(receiveQuestion)
	err = qc.questionService.AddQuestion(question)
	result := model_question.QuestionToReturnQuestion(question)
	if err != nil {
		log.Printf("add model_question %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "add question error").Response(api_response.OPERATIONERR))
		return
	}

	log.Printf("add model_question success")
	c.JSON(http.StatusOK, api_response.NewResponse(result, "add question success").Response(api_response.SUCCESS))
}

// GetQuestion 获取题目
//
//	@Summary		Query question
//	@Description	Query question
//	@Tags			Question
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string													true	"Question id"
//	@Success		200	{object}	api_response.ApiResponse{data=model_question.ReturnQuestion}	"Query question success"
//	@Failure		400	{object}	api_response.ApiResponse{data=nil}								"Query question fail"
//	@Router			/api/question/query    [get]
func (qc *QuestionController) GetQuestion(c *gin.Context) {
	session, err := qc.session.GetSession(c)
	if session.UserRole != constant.Common && session.UserRole != constant.Admin {
		log.Printf("you are not allowed to get a new question")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "you are not allowed to get a new question").Response(api_response.AUTHERR))

		return
	}
	id := c.Param("id")
	question, err := qc.questionService.GetQuestion(id)
	if err != nil {
		log.Printf("query model_question %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "query question error").Response(api_response.OPERATIONERR))
		return
	}
	res := model_question.QuestionToReturnQuestion(question)
	if session.UserRole != constant.Admin {
		res.Answer = nil
	}
	log.Printf("query model_question success")
	c.JSON(http.StatusOK, api_response.NewResponse(res, "query question success").Response(api_response.SUCCESS))
}

// GetQuestionList 获取题目
//
//	@Summary		Get question list
//	@Description	Get question list
//	@Tags			Question
//	@Accept			json
//	@Produce		json
//	@Param			question	body		model_question.QueryQuestionRequest							true	"Query conditions"
//	@Success		200	{object}	api_response.ApiResponse{data=[]model_question.ReturnQuestion}	"Get question list success"
//	@Failure		400	{object}	api_response.ApiResponse{data=nil}								"Get question list failed"
//	@Router			/api/question/query    [post]
func (qc *QuestionController) GetQuestionList(c *gin.Context) {
	session, _ := qc.session.GetSession(c)
	if session.UserRole != constant.Admin && session.UserRole != constant.Common {
		log.Printf("you are not allowed to get questions")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "you are not allowed to create get questions ").Response(api_response.AUTHERR))

		return
	}
	var receiveQuestion model_question.QueryQuestionRequest
	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&receiveQuestion); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "JSON unmarshal error").Response(api_response.OPERATIONERR))

		return
	}

	page := receiveQuestion.Page
	pageSize := receiveQuestion.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	//字段验证
	query := model_question.QueryQuestionToQuestion(receiveQuestion)
	err := qc.checkQueryOrUpdateQuestion(query)
	if err != nil {
		log.Printf("validate %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, err.Error()).Response(api_response.OPERATIONERR))

		return
	}

	if err != nil {
		log.Printf("validate %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "validate error").Response(api_response.OPERATIONERR))

		return
	}
	commonQuery := model_question.QueryQToCommonQueryQ(receiveQuestion)
	questionList, err := qc.questionService.GetQuestionList(commonQuery, page, pageSize)
	if err != nil {
		log.Printf("query questions %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "query questions error").Response(api_response.OPERATIONERR))
		return
	}
	res := model_question.QuestionsToReturnQuestions(questionList)

	if session.UserRole != constant.Admin {
		for _, v := range res {
			v.Answer = nil
		}
	}

	log.Printf("query questions success")
	c.JSON(http.StatusOK, api_response.NewResponse(res, "query questions success").Response(api_response.SUCCESS))
}

// UpdateQuestion 更新题目
//
//	@Summary		Update question
//	@Description	Update question
//	@Tags			Question
//	@Accept			json
//	@Produce		json
//	@Param			question	body		model_question.UpdateQuestionRequest	true	"Update condition"
//	@Success		200			{object}	api_response.ApiResponse{data=nil}			"Update success"
//	@Failure		400			{object}	api_response.ApiResponse{data=nil}			"Update fail"
//	@Router			/api/question/admin/update    [post]
func (qc *QuestionController) UpdateQuestion(c *gin.Context) {
	session, _ := qc.session.GetSession(c)
	if session.UserRole != constant.Admin {
		log.Printf("you are not allowed to update a new question")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "you are not allowed to update a new question").Response(api_response.AUTHERR))

		return
	}
	var receiveQuestion model_question.UpdateQuestionRequest

	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&receiveQuestion); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "JSON unmarshal error").Response(api_response.OPERATIONERR))

		return
	}

	//字段验证
	var old model_question.Question
	update := model_question.UpdateQuestionToQuestion(old, receiveQuestion)
	err := qc.checkQueryOrUpdateQuestion(update)
	if err != nil {
		log.Printf("validate %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, err.Error()).Response(api_response.AUTHERR))

		return
	}

	//获取题目
	queryQuestion, err := qc.questionService.GetQuestion(receiveQuestion.ID)
	if err != nil {
		log.Printf("query model_question %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "query question error").Response(api_response.OPERATIONERR))
		return
	}

	question := model_question.UpdateQuestionToQuestion(queryQuestion, receiveQuestion)
	err = qc.questionService.UpdateQuestion(question)
	if err != nil {
		log.Printf("update model_question %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "update question error").Response(api_response.OPERATIONERR))
		return
	}

	log.Printf("update model_question success")
	c.JSON(http.StatusOK, api_response.NewResponse(nil, "update question success").Response(api_response.SUCCESS))
}

// DeleteQuestion 删除题目
//
//	@Summary		Delete question
//	@Description	Delete question
//	@Tags			Question
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						true	"Question id"
//	@Success		200	{object}	api_response.ApiResponse{data=nil}	"Delete  success"
//	@Failure		400	{object}	api_response.ApiResponse{data=nil}	"Delete fail"
//	@Router			/api/question/admin/delete    [get]
func (qc *QuestionController) DeleteQuestion(c *gin.Context) {
	session, _ := qc.session.GetSession(c)
	if session.UserRole != constant.Admin {
		log.Printf("you are not allowed to delete a question")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "you are not allowed to delete a question").Response(api_response.AUTHERR))

		return
	}

	id := c.Param("id")
	if id == "" {
		log.Printf("id is empty")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "id is empty").Response(api_response.OPERATIONERR))

		return
	}

	err := qc.questionService.DeleteQuestion(id)
	if err != nil {
		log.Printf("delete model_question %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "delete question error").Response(api_response.OPERATIONERR))
		return
	}

	log.Printf("delete model_question success")
	c.JSON(http.StatusOK, api_response.NewResponse(nil, "delete question success").Response(api_response.SUCCESS))
}

func (qc *QuestionController) checkQuestion(question model_question.Question) error {
	if question.Content == "" || len(question.Content) > 8192 {
		return errors.New("content is empty or too long")
	}
	if question.Answer == "" || len(question.Answer) > 8192 {
		return errors.New("answer is empty or too long")
	}
	if question.JudgeCase == "" {
		return errors.New("judge_case is empty")
	}
	if question.JudgeConfig == "" {
		return errors.New("judge_config is empty")
	}
	return nil
}

func (qc *QuestionController) checkQueryOrUpdateQuestion(question model_question.Question) error {
	if question.ID != "" && len(question.ID) > 256 {
		return errors.New("id is too long")
	}
	if question.Title != "" && len(question.Title) > 256 {
		return errors.New("title is  too long")
	}
	if question.Content != "" && len(question.Content) > 8192 {
		return errors.New("content is empty or too long")
	}

	if question.Tag != "" && len(question.Tag) > 256 {
		return errors.New("tags is too long")
	}

	if question.Answer != "" && len(question.Answer) > 8192 {
		return errors.New("answer is  too long")
	}

	if question.JudgeCase != "" && len(question.JudgeCase) > 8192 {
		return errors.New("judge case is  too long")
	}
	if question.UserId != "" && len(question.UserId) > 256 {
		return errors.New("user id is too long")
	}
	if question.JudgeConfig != "" && len(question.JudgeConfig) > 64 {
		return errors.New("judge config is  too long")
	}
	return nil
}
