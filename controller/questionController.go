package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/common/api_response"
	"github.com/xissg/userManageSystem/common/constant"
	"github.com/xissg/userManageSystem/entity/modelquestion"
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
//	@Summary		add question
//	@Description	add question
//	@Tags			Question
//	@Accept			json
//	@Produce		json
//	@Param			question	body		modelquestion.AddQuestionRequest	true	"add question"
//	@Success		200			{object}	utils.ApiResponse{data=nil}			"add successful"
//	@Failure		404			{object}	utils.ApiResponse{data=nil}			"add failed"
//	@Router			/api/question/admin/create    [post]
func (qc *QuestionController) AddQuestion(c *gin.Context) {

	session, _ := qc.session.GetSession(c)
	if session.UserRole != constant.Admin {
		log.Printf("you are not allowed to create a new question")
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "you are not allowed to create a new question").Response(api_response.AUTHERR))

		return
	}

	var receiveQuestion modelquestion.AddQuestionRequest
	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&receiveQuestion); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "unmarshal error ").Response(api_response.OPERATIONERR))

		return
	}

	add := modelquestion.AddQuestionToQuestion(receiveQuestion)
	//数据校验
	err := qc.checkQuestion(add)
	if err != nil {
		log.Printf("validate %v", err)
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "validate error").Response(api_response.OPERATIONERR))

		return
	}
	question := modelquestion.AddQuestionToQuestion(receiveQuestion)
	err = qc.questionService.AddQuestion(question)
	result := modelquestion.QuestionToReturnQuestion(question)
	if err != nil {
		log.Printf("add modelquestion %v", err)
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "add question error").Response(api_response.OPERATIONERR))
		return
	}

	log.Printf("add modelquestion success")
	c.JSON(http.StatusOK, api_response.NewResponse(result, "add question success").Response(api_response.SUCCESS))
}

// GetQuestion 获取题目
//
//	@Summary		get question
//	@Description	get question
//	@Tags			Question
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string													true	"Get Question"
//	@Success		200	{object}	utils.ApiResponse{data=modelquestion.ReturnQuestion}	"get question successful"
//	@Failure		404	{object}	utils.ApiResponse{data=nil}								"EditUserRequest failed"	"get question failed"
//	@Router			/api/question/get    [get]
func (qc *QuestionController) GetQuestion(c *gin.Context) {
	session, err := qc.session.GetSession(c)
	if session.UserRole != constant.Common && session.UserRole != constant.Admin {
		log.Printf("you are not allowed to get a new question")
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "you are not allowed to get a new question").Response(api_response.AUTHERR))

		return
	}
	id := c.Param("id")
	question, err := qc.questionService.GetQuestion(id)
	if err != nil {
		log.Printf("query modelquestion %v", err)
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "query question error").Response(api_response.OPERATIONERR))
		return
	}
	res := modelquestion.QuestionToReturnQuestion(question)

	log.Printf("query modelquestion success")
	c.JSON(http.StatusOK, api_response.NewResponse(res, "query question success").Response(api_response.SUCCESS))
}

// GetQuestionList 获取题目
//
//	@Summary		get questions
//	@Description	get questions
//	@Tags			Question
//	@Accept			json
//	@Produce		json
//	@Param			question	body		modelquestion.QueryQuestionRequest							true	"Get Questions"
//	@Success		200	{object}	utils.ApiResponse{data=[]modelquestion.ReturnQuestion}	"get questions successful"
//	@Failure		404	{object}	utils.ApiResponse{data=nil}								"EditUserRequest failed"	"get questions failed"
//	@Router			/api/question/get    [post]
func (qc *QuestionController) GetQuestionList(c *gin.Context) {
	session, _ := qc.session.GetSession(c)
	if session.UserRole != constant.Admin && session.UserRole != constant.Common {
		log.Printf("you are not allowed to get questions")
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "you are not allowed to create get questions ").Response(api_response.AUTHERR))

		return
	}
	var receiveQuestion modelquestion.QueryQuestionRequest
	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&receiveQuestion); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "JSON unmarshal error").Response(api_response.OPERATIONERR))

		return
	}

	//字段验证
	query := modelquestion.QueryQuestionToQuestion(receiveQuestion)
	err := qc.checkQueryOrUpdateQuestion(query)
	if err != nil {
		log.Printf("validate %v", err)
		c.JSON(http.StatusOK, api_response.NewResponse(nil, err.Error()).Response(api_response.OPERATIONERR))

		return
	}

	if err != nil {
		log.Printf("validate %v", err)
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "validate error").Response(api_response.OPERATIONERR))

		return
	}
	commonQuery := modelquestion.QueryQToCommonQueryQ(receiveQuestion)
	questionList, err := qc.questionService.GetQuestionList(commonQuery)
	if err != nil {
		log.Printf("query questions %v", err)
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "query questions error").Response(api_response.OPERATIONERR))
		return
	}
	res := modelquestion.QuestionsToReturnQuestions(questionList)

	log.Printf("query questions success")
	c.JSON(http.StatusOK, api_response.NewResponse(res, "query questions success").Response(api_response.SUCCESS))
}

// UpdateQuestion 添加题目
//
//	@Summary		update question
//	@Description	update question
//	@Tags			Question
//	@Accept			json
//	@Produce		json
//	@Param			question	body		modelquestion.UpdateQuestionRequest	true	"update question"
//	@Success		200			{object}	utils.ApiResponse{data=nil}			"update successful"
//	@Failure		404			{object}	utils.ApiResponse{data=nil}			"update failed"
//	@Router			/api/question/admin/update    [post]
func (qc *QuestionController) UpdateQuestion(c *gin.Context) {
	session, _ := qc.session.GetSession(c)
	if session.UserRole != constant.Admin {
		log.Printf("you are not allowed to update a new question")
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "you are not allowed to update a new question").Response(api_response.AUTHERR))

		return
	}
	var receiveQuestion modelquestion.UpdateQuestionRequest

	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&receiveQuestion); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "JSON unmarshal error").Response(api_response.OPERATIONERR))

		return
	}

	//字段验证
	var old modelquestion.Question
	update := modelquestion.UpdateQuestionToQuestion(old, receiveQuestion)
	err := qc.checkQueryOrUpdateQuestion(update)
	if err != nil {
		log.Printf("validate %v", err)
		c.JSON(http.StatusOK, api_response.NewResponse(nil, err.Error()).Response(api_response.AUTHERR))

		return
	}

	//获取题目
	queryQuestion, err := qc.questionService.GetQuestion(receiveQuestion.ID)
	if err != nil {
		log.Printf("query modelquestion %v", err)
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "query question error").Response(api_response.OPERATIONERR))
		return
	}

	question := modelquestion.UpdateQuestionToQuestion(queryQuestion, receiveQuestion)
	err = qc.questionService.UpdateQuestion(question)
	if err != nil {
		log.Printf("update modelquestion %v", err)
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "update question error").Response(api_response.OPERATIONERR))
		return
	}

	log.Printf("update modelquestion success")
	c.JSON(http.StatusOK, api_response.NewResponse(nil, "update question success").Response(api_response.SUCCESS))
}

// DeleteQuestion 删除题目
//
//	@Summary		delete question
//	@Description	delete question
//	@Tags			Question
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						true	"delete Question"
//	@Success		200	{object}	utils.ApiResponse{data=nil}	"delete  successful"
//	@Failure		404	{object}	utils.ApiResponse{data=nil}	"delete failed"
//	@Router			/api/question/admin/delete    [get]
func (qc *QuestionController) DeleteQuestion(c *gin.Context) {
	session, _ := qc.session.GetSession(c)
	if session.UserRole != constant.Admin {
		log.Printf("you are not allowed to delete a question")
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "you are not allowed to delete a question").Response(api_response.AUTHERR))

		return
	}

	id := c.Param("id")
	if id == "" {
		log.Printf("id is empty")
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "id is empty").Response(api_response.OPERATIONERR))

		return
	}

	err := qc.questionService.DeleteQuestion(id)
	if err != nil {
		log.Printf("delete modelquestion %v", err)
		c.JSON(http.StatusOK, api_response.NewResponse(nil, "delete question error").Response(api_response.OPERATIONERR))
		return
	}

	log.Printf("delete modelquestion success")
	c.JSON(http.StatusOK, api_response.NewResponse(nil, "delete question success").Response(api_response.SUCCESS))
}

func (qc *QuestionController) checkQuestion(question modelquestion.Question) error {
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

func (qc *QuestionController) checkQueryOrUpdateQuestion(question modelquestion.Question) error {
	if question.ID != "" && len(question.ID) > 256 {
		return errors.New("id is too long")
	}
	if question.Title != "" && len(question.Title) > 256 {
		return errors.New("title is  too long")
	}
	if question.Content != "" && len(question.Content) > 8192 {
		return errors.New("content is empty or too long")
	}

	if question.Tags != "" && len(question.Tags) > 256 {
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
