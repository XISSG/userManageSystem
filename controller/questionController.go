package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/constant"
	"github.com/xissg/userManageSystem/entity/modelquestion"
	mysql2 "github.com/xissg/userManageSystem/service/mysql"
	"github.com/xissg/userManageSystem/service/redis"
	"github.com/xissg/userManageSystem/utils"
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
		c.JSON(http.StatusOK, utils.NewResponse(nil, "you are not allowed to create a new question").AuthERR())

		return
	}

	var receiveQuestion modelquestion.AddQuestionRequest
	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&receiveQuestion); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		c.JSON(http.StatusOK, utils.NewResponse(nil, "unmarshal error ").OperationERR())

		return
	}

	add := modelquestion.AddQuestionToQuestion(receiveQuestion)
	//数据校验
	err := qc.checkQuestion(add)
	if err != nil {
		log.Printf("validate %v", err)
		c.JSON(http.StatusOK, utils.NewResponse(nil, "validate error").OperationERR())

		return
	}
	question := modelquestion.AddQuestionToQuestion(receiveQuestion)
	err = qc.questionService.AddQuestion(question)
	result := modelquestion.QuestionToReturnQuestion(question)
	if err != nil {
		log.Printf("add modelquestion %v", err)
		c.JSON(http.StatusOK, utils.NewResponse(nil, "add question error").OperationERR())
		return
	}

	log.Printf("add modelquestion success")
	c.JSON(http.StatusOK, utils.NewResponse(result, "add question success").Success())
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
		c.JSON(http.StatusOK, utils.NewResponse(nil, "you are not allowed to get a new question").AuthERR())

		return
	}
	id := c.Param("id")
	question, err := qc.questionService.GetQuestion(id)
	if err != nil {
		log.Printf("query modelquestion %v", err)
		c.JSON(http.StatusOK, utils.NewResponse(nil, "query question error").OperationERR())
		return
	}
	res := modelquestion.QuestionToReturnQuestion(question)

	log.Printf("query modelquestion success")
	c.JSON(http.StatusOK, utils.NewResponse(res, "query question success").Success())
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
		c.JSON(http.StatusOK, utils.NewResponse(nil, "you are not allowed to create get questions ").AuthERR())

		return
	}
	var receiveQuestion modelquestion.QueryQuestionRequest
	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&receiveQuestion); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		c.JSON(http.StatusOK, utils.NewResponse(nil, "JSON unmarshal error").OperationERR())

		return
	}

	//字段验证
	query := modelquestion.QueryQuestionToQuestion(receiveQuestion)
	err := qc.checkQueryOrUpdateQuestion(query)
	if err != nil {
		log.Printf("validate %v", err)
		c.JSON(http.StatusOK, utils.NewResponse(nil, err.Error()).OperationERR())

		return
	}

	if err != nil {
		log.Printf("validate %v", err)
		c.JSON(http.StatusOK, utils.NewResponse(nil, "validate error").OperationERR())

		return
	}
	commonQuery := modelquestion.QueryQToCommonQueryQ(receiveQuestion)
	questionList, err := qc.questionService.GetQuestionList(commonQuery)
	if err != nil {
		log.Printf("query questions %v", err)
		c.JSON(http.StatusOK, utils.NewResponse(nil, "query questions error").OperationERR())
		return
	}
	res := modelquestion.QuestionsToReturnQuestions(questionList)

	log.Printf("query questions success")
	c.JSON(http.StatusOK, utils.NewResponse(res, "query questions success").Success())
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
		c.JSON(http.StatusOK, utils.NewResponse(nil, "you are not allowed to update a new question").AuthERR())

		return
	}
	var receiveQuestion modelquestion.UpdateQuestionRequest

	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&receiveQuestion); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		c.JSON(http.StatusOK, utils.NewResponse(nil, "JSON unmarshal error").OperationERR())

		return
	}

	//字段验证
	var old modelquestion.Question
	update := modelquestion.UpdateQuestionToQuestion(old, receiveQuestion)
	err := qc.checkQueryOrUpdateQuestion(update)
	if err != nil {
		log.Printf("validate %v", err)
		c.JSON(http.StatusOK, utils.NewResponse(nil, err.Error()).AuthERR())

		return
	}

	//获取题目
	queryQuestion, err := qc.questionService.GetQuestion(receiveQuestion.ID)
	if err != nil {
		log.Printf("query modelquestion %v", err)
		c.JSON(http.StatusOK, utils.NewResponse(nil, "query question error").OperationERR())
		return
	}

	question := modelquestion.UpdateQuestionToQuestion(queryQuestion, receiveQuestion)
	err = qc.questionService.UpdateQuestion(question)
	if err != nil {
		log.Printf("update modelquestion %v", err)
		c.JSON(http.StatusOK, utils.NewResponse(nil, "update question error").OperationERR())
		return
	}

	log.Printf("update modelquestion success")
	c.JSON(http.StatusOK, utils.NewResponse(nil, "update question success").Success())
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
		c.JSON(http.StatusOK, utils.NewResponse(nil, "you are not allowed to delete a question").AuthERR())

		return
	}

	id := c.Param("id")
	if id == "" {
		log.Printf("id is empty")
		c.JSON(http.StatusOK, utils.NewResponse(nil, "id is empty").OperationERR())

		return
	}

	err := qc.questionService.DeleteQuestion(id)
	if err != nil {
		log.Printf("delete modelquestion %v", err)
		c.JSON(http.StatusOK, utils.NewResponse(nil, "delete question error").OperationERR())
		return
	}

	log.Printf("delete modelquestion success")
	c.JSON(http.StatusOK, utils.NewResponse(nil, "delete question success").Success())
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
