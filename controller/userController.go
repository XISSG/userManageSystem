package controller

import (
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/model/entity"
	"github.com/xissg/userManageSystem/service"
	"github.com/xissg/userManageSystem/utils"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// 用户的登录状态
const (
	ANONYMOUS = iota
	COMMON
	ADMIN //user_role 字段判断是否为admin用户
)

type UserController struct {
	sessionService service.SessionService
	userService    *service.UserService
}

func NewUserController(userService service.UserService, sessionService service.SessionService) *UserController {

	return &UserController{
		sessionService: sessionService,
		userService:    &userService,
	}
}

// 校验用户是否登录以及是否为admin用户
func (uc *UserController) checkValidity(c *gin.Context) (int, entity.UserSession) {
	userSession, err := uc.sessionService.GetSession(c)
	if err != nil {
		return ANONYMOUS, entity.UserSession{}
	}
	if userSession.Role == ADMIN {

		return ADMIN, userSession
	}

	return COMMON, userSession
}

// Register 用户注册
//
//	@Summary		User registration
//	@Description	Register a new user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		string										true	"User object"
//	@Success		200		{object}	utils.apiResponse{data=model.ReturnUser}	"Success registered"
//	@Failure		400		{object}	utils.apiResponse{data=nil}					"Registration failed"
//	@Router			/user/register [post]
func (uc *UserController) Register(c *gin.Context) {
	var receiveUser entity.AddUser

	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&receiveUser); err != nil {
		log.Printf("JSON unmarshal  %v", err)

		return
	}

	//校验用户名长度
	if len(receiveUser.UserName) < 4 || len(receiveUser.UserName) > 32 {
		c.JSON(http.StatusBadRequest, utils.NewResponse(nil, "receiveUser name too short or too long").LoginERR())
		log.Printf("wrong receiveUser name length")

		return
	}

	//校验密码长度
	if len(receiveUser.UserPassword) < 8 || len(receiveUser.UserPassword) > 32 {
		c.JSON(http.StatusBadRequest, utils.NewResponse(nil, "wrong password length").LoginERR())
		log.Printf("wrong password length")

		return
	}

	//校验密码合法性
	expr := `^(?![0-9a-zA-Z]+$)(?![a-zA-Z!@#$%^&*]+$)(?![0-9!@#$%^&*]+$)[0-9A-Za-z!@#$%^&*]{8,16}$`
	reg, _ := regexp2.Compile(expr, 0)
	m, _ := reg.MatchString(receiveUser.UserPassword)
	if !m {
		c.JSON(http.StatusBadRequest, utils.NewResponse(nil, "Invalid password, At least one special character, lowercase and uppercase, is required").LoginERR())
		log.Printf("illegal receiveUser password")

		return
	}

	//校验用户名是否存在
	_, err := uc.userService.GetUser(receiveUser.UserName, c)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, utils.NewResponse(nil, "username repeated").LoginERR())
		log.Println("username repeated")

		return
	}

	//用户密码脱敏
	receiveUser.UserPassword = utils.MD5Crypt(receiveUser.UserPassword)

	//生成用户
	var user entity.User
	user = entity.AddUserToUser(receiveUser)

	//插入数据库
	err = uc.userService.AddUser(user, c)
	if err != nil {
		log.Printf("create receiveUser failed")

		return
	}

	//插入成功
	resultUser := entity.SafetyUser(user)
	log.Printf("register success")
	c.JSON(http.StatusOK, utils.NewResponse(resultUser, "register success").Success())

}

// Login 用户登录
//
//	@Summary		User login
//	@Description	Authenticate user login
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		string										true	"User object"
//	@Success		200		{object}	utils.apiResponse{data=model.ReturnUser}	"Login successful"
//	@Failure		400		{object}	utils.apiResponse{data=nil}					"Login failed"
//	@Router			/user/login [post]
func (uc *UserController) Login(c *gin.Context) {
	var user entity.LoginUser

	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("JSON unmarshal  %v", err)

		return
	}

	//判断用户名和密码是否为空
	if user.UserName == "" || user.UserPassword == "" {
		log.Printf("username  or password  is empty")
		c.JSON(http.StatusBadRequest, utils.NewResponse(nil, "username or password is empty").LoginERR())

		return
	}

	//用户密码加密
	user.UserPassword = utils.MD5Crypt(user.UserPassword)

	//查询用户用户名和密码是否匹配
	ret, err := uc.userService.GetUser(user.UserName, c)
	if ret == nil {
		log.Println("The user has not registered yet")
		c.JSON(http.StatusBadRequest, utils.NewResponse(nil, "The user has not registered yet").AuthERR())

		return
	}

	re := ret.(entity.User)
	if re.UserPassword != user.UserPassword {
		log.Println("username or password is wrong")
		c.JSON(http.StatusBadRequest, utils.NewResponse(nil, "username or password is wrong").AuthERR())

		return
	}

	//登录成功, 存储session信息
	if re.UserRole == ANONYMOUS {
		re.UserRole = COMMON
	}
	userSession := entity.UserSession{
		ID:       re.ID,
		UserName: re.UserName,
		Role:     re.UserRole,
	}

	err = uc.sessionService.NewOrUpdateSession(c, userSession)
	if err != nil {
		log.Println(fmt.Sprintf("session create %v", err))

		return
	}

	log.Printf("login success")
	c.JSON(http.StatusOK, utils.NewResponse(nil, "login success").Success())
}

// Logout 登出账户
//
//	@Summary		User logout
//	@Description	User logout
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		string						true	"User object"
//	@Success		200		{object}	utils.apiResponse{data=nil}	"Logout successful"
//	@Failure		400		{object}	utils.apiResponse{data=nil}	"Logout failed"
//	@Router			/user/logout [get]
func (uc *UserController) Logout(c *gin.Context) {
	//判断用户是否登录
	validity, _ := uc.checkValidity(c)
	if validity == ANONYMOUS {
		log.Printf("you must login first")
		c.JSON(http.StatusBadRequest, utils.NewResponse(nil, "you must login first").LoginERR())

		return
	}

	//删除session
	err := uc.sessionService.DeleteSession(c)
	if err != nil {
		log.Println(fmt.Sprintf("session delete %v", err))

		return
	}

	log.Printf("logout success")
	c.JSON(http.StatusOK, utils.NewResponse(nil, "logout success").Success())
}

// QueryUser 查询用户
//
//	@Summary		Query user by username
//	@Description	Get user information by username
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			username	query		string										true	"Username"
//	@Success		200			{object}	utils.apiResponse{data=model.ReturnUser}	"Query successful"
//	@Failure		400			{object}	utils.apiResponse{data=nil}					"Query failed"
//	@Router			/user/admin/query [get]
func (uc *UserController) QueryUser(c *gin.Context) {
	//判断用户权限
	validity, _ := uc.checkValidity(c)

	if validity != ADMIN {
		c.JSON(http.StatusBadRequest, utils.NewResponse(nil, "you are not admin").AuthERR())

		return
	}

	//获取username参数的值
	username := c.Param("username")
	if username == "" {
		log.Println("not a valid query username")
		c.JSON(http.StatusBadRequest, utils.NewResponse(nil, "not a valid query username").ParamsERR())

		return
	}

	res, err := uc.userService.GetUser(username, c)
	if res == nil {
		log.Println("query user not found")
		c.JSON(http.StatusBadRequest, utils.NewResponse(nil, "query user not found").ParamsERR())

		return
	}

	if err != nil {
		log.Println(fmt.Sprintf("query user %v", err))
		c.JSON(http.StatusBadRequest, utils.NewResponse(nil, "query user error").OperationERR())

		return
	}

	ret := entity.SafetyUser(res.(entity.User))
	log.Println("query user success")
	c.JSON(http.StatusOK, utils.NewResponse(ret, "query user success").Success())
}

// UpdateUser 更新用户信息
//
//	@Summary		User update
//	@Description	UpdateUser user information
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		string						true	"User object"
//	@Success		200		{object}	utils.apiResponse{data=nil}	"UpdateUser successful"
//	@Failure		400		{object}	utils.apiResponse{data=nil}	"UpdateUser failed"
//	@Router			/user/update [post]
func (uc *UserController) UpdateUser(c *gin.Context) {
	//判断用户权限
	validity, userSession := uc.checkValidity(c)

	if validity == ANONYMOUS {
		log.Printf("you are not login")
		c.JSON(http.StatusBadRequest, utils.NewResponse(nil, "you are not login").AuthERR())

		return
	}

	var user entity.UpdateUser
	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("JSON unmarshal  %v", err)

		return
	}

	//密码要加密
	if user.UserPassword != "" {
		user.UserPassword = utils.MD5Crypt(user.UserPassword)
	}

	//普通用户只允许修改自己的信息
	if validity != ADMIN {
		user.UserName = userSession.UserName
	}

	count, column := entity.CountParams(user)
	if count > 1 {
		//admin用户允许修改其他人的信息
		err := uc.userService.UpdateUserAll(user, c)
		if err != nil {
			log.Printf("user info update  %v", err)

			return
		}
	}

	if count == 1 {
		err := uc.userService.UpdateUserOne(column, user, c)
		if err != nil {
			log.Printf("user info update  %v", err)

			return
		}
	}

	log.Printf("update user success")
	c.JSON(http.StatusOK, utils.NewResponse(nil, "update user success").Success())
}

// DeleteUser 删除用户
//
//	@Summary		DeleteUser user by username
//	@Description	DeleteUser user information by username
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			username	query		string						true	"Username"
//	@Success		200			{object}	utils.apiResponse{data=nil}	"DeleteUser successful"
//	@Failure		400			{object}	utils.apiResponse{data=nil}	"DeleteUser failed"
//	@Router			/user/admin/delete [get]
func (uc *UserController) DeleteUser(c *gin.Context) {

	//判断用户权限
	validity, _ := uc.checkValidity(c)
	if validity != ADMIN {
		c.JSON(http.StatusBadRequest, utils.NewResponse(nil, "you are not validity user").AuthERR())

		return
	}

	//获取username参数的值
	username := c.Param("username")
	if username == "" {
		log.Println("not a valid query username")
		c.JSON(http.StatusBadRequest, utils.NewResponse(nil, "not a valid query username").ParamsERR())

		return
	}

	//逻辑删除用户
	err := uc.userService.DeleteUser(username, c)
	if err != nil {
		log.Printf("delete user  %v", err)
		c.JSON(http.StatusBadRequest, utils.NewResponse(nil, err.Error()).OperationERR())

		return
	}

	log.Printf("delete user success")
	c.JSON(http.StatusOK, utils.NewResponse(nil, "delete user success").Success())
}
