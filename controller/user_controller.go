package controller

import (
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/common/api_response"
	"github.com/xissg/userManageSystem/common/constant"
	"github.com/xissg/userManageSystem/entity/model_user"
	"github.com/xissg/userManageSystem/service/mysql"
	"github.com/xissg/userManageSystem/service/redis"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type UserController struct {
	sessionService *redis.SessionService
	userService    *mysql.UserService
}

func NewUserController(userService mysql.UserService, sessionService redis.SessionService) *UserController {

	return &UserController{
		sessionService: &sessionService,
		userService:    &userService,
	}
}

// Register 用户注册
//
//	@Summary		Register
//	@Description	Register a new user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		model_user.AddUserRequest	true	"User information"
//	@Success		200		{object}	api_response.ApiResponse{data=nil}	"Register success"
//	@Failure		400		{object}	api_response.ApiResponse{data=nil}	"Register fail"
//	@Router			/api/user/register  [post]
func (uc *UserController) Register(c *gin.Context) {
	var receiveUser model_user.AddUserRequest

	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&receiveUser); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "unmarshal error ").Response(api_response.OPERATIONERR))

		return
	}

	//校验字段合法性
	err := uc.checkUser(receiveUser.UserAccount, receiveUser.UserPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, err.Error()).Response(api_response.AUTHERR))
		log.Printf("validate %v", err)

		return
	}

	//校验账户是否存在
	_, err = uc.userService.GetUser(receiveUser.UserAccount)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "user account repeated").Response(api_response.AUTHERR))
		log.Println("user account repeated")

		return
	}

	//生成用户
	var user model_user.User
	user = model_user.AddUserToUser(receiveUser)

	//插入数据库
	err = uc.userService.AddUser(user)
	if err != nil {
		log.Printf("create user failed")

		return
	}

	log.Printf("register success")
	c.JSON(http.StatusOK, api_response.NewResponse(nil, "register success").Response(api_response.SUCCESS))

}

// Login 用户登录
//
//	@Summary		Login
//	@Description	User login
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		model_user.LoginUserRequest						true	"User information"
//	@Success		200		{object}	api_response.ApiResponse{data=model_user.ReturnUser}	"Login success"
//	@Failure		400		{object}	api_response.ApiResponse{data=nil}						"Login fail"
//	@Router			/api/user/login     [post]
func (uc *UserController) Login(c *gin.Context) {
	var loginUser model_user.LoginUserRequest

	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "unmarshal error ").Response(api_response.OPERATIONERR))

		return
	}

	//验证字段合法性
	err := uc.checkUser(loginUser.UserAccount, loginUser.UserPassword)
	if err != nil {
		log.Printf("%v invalid user account or password", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, err.Error()).Response(api_response.AUTHERR))
		return
	}

	user := model_user.LoginUserToUser(loginUser)
	//查询用户账户和密码是否匹配
	ret, err := uc.userService.GetUser(user.UserAccount)
	if ret.UserAccount == "" {
		log.Println("The user has not registered yet")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "The user has not registered yet").Response(api_response.AUTHERR))

		return
	}

	//禁用的账号
	if ret.UserRole == constant.Ban {
		log.Println("The user has been banned")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "The user has been banned").Response(api_response.AUTHERR))

		return
	}

	if ret.UserPassword != user.UserPassword {
		log.Println("username or password is wrong")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "username or password is wrong").Response(api_response.AUTHERR))

		return
	}

	//登录成功, 存储session信息
	userSession := model_user.UserToUserSession(ret)
	err = uc.sessionService.NewOrUpdateSession(c, userSession)
	if err != nil {
		log.Println(fmt.Sprintf("session create %v", err))

		return
	}

	//插入成功
	resultUser := model_user.UserToReturnUser(ret)
	log.Printf("login success")
	c.JSON(http.StatusOK, api_response.NewResponse(resultUser, "login success").Response(api_response.SUCCESS))
}

// Logout 登出账户
//
//	@Summary		Logout
//	@Description	User logout
//	@Tags			User
//	@Produce		json
//	@Success		200	{object}	api_response.ApiResponse{data=nil}	"Logout success"
//	@Failure		400	{object}	api_response.ApiResponse{data=nil}	"Logout fail"
//	@Router			/api/user/logout    [get]
func (uc *UserController) Logout(c *gin.Context) {
	//判断用户是否登录
	validity, _ := uc.sessionService.GetSession(c)
	if validity.UserRole == constant.Anonymous {
		log.Printf("you must login first")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "you must login first").Response(api_response.AUTHERR))

		return
	}

	//删除session
	err := uc.sessionService.DeleteSession(c)
	if err != nil {
		log.Println(fmt.Sprintf("session delete %v", err))

		return
	}

	log.Printf("logout success")
	c.JSON(http.StatusOK, api_response.NewResponse(nil, "logout success").Response(api_response.SUCCESS))
}

// GetUserList 查询用户列表
//
//	@Summary		Query user
//	@Description	Query user list
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		model_user.UserQueryRequest						true	"Query"
//	@Success		200		{object}	api_response.ApiResponse{data=[]model_user.ReturnUser}	"Query success"
//	@Failure		400		{object}	api_response.ApiResponse{data=nil}						"Query fail"
//	@Router			/api/user/query [post]
func (uc *UserController) GetUserList(c *gin.Context) {
	//判断用户权限
	validity, _ := uc.sessionService.GetSession(c)

	if validity.UserRole == constant.Anonymous || validity.UserRole == constant.Ban {
		log.Printf("you must login first")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "you must login first").Response(api_response.AUTHERR))

		return
	}
	var queryRequest model_user.UserQueryRequest
	user := model_user.UserQueryToUser(queryRequest)

	err := uc.checkQueryOrUpdateUser(user)
	if err != nil {
		log.Printf("validate %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, err.Error()).Response(api_response.AUTHERR))
		return
	}
	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&queryRequest); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "unmarshal error ").Response(api_response.OPERATIONERR))

		return
	}

	page := queryRequest.Page
	pageSize := queryRequest.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	commonQuery := model_user.UserQueryToCommonQuery(queryRequest)
	res, err := uc.userService.GetUserList(commonQuery, page, pageSize)
	if err != nil {
		log.Println(fmt.Sprintf("query user %v", err))
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "query user error").Response(api_response.OPERATIONERR))

		return
	}

	ret := model_user.UsersToReturnUsers(res)
	log.Println("query users success")
	c.JSON(http.StatusOK, api_response.NewResponse(ret, "query users success").Response(api_response.SUCCESS))
}

// AdminGetUserList 查询用户列表
//
//	@Summary		Admin query
//	@Description	Query user list for admin
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		model_user.AdminUserQueryRequest							true	"queries"
//	@Success		200		{object}	api_response.ApiResponse{data=[]model_user.ReturnAdminUser}	"Query success"
//	@Failure		400		{object}	api_response.ApiResponse{data=nil}							"Query fail"
//	@Router			/api/user/admin/query [post]
func (uc *UserController) AdminGetUserList(c *gin.Context) {
	//判断用户权限
	validity, _ := uc.sessionService.GetSession(c)
	if validity.UserRole != constant.Admin {
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "you are not admin").Response(api_response.AUTHERR))

		return
	}

	var queryRequest model_user.AdminUserQueryRequest
	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&queryRequest); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "unmarshal error ").Response(api_response.OPERATIONERR))

		return
	}
	page := queryRequest.Page
	pageSize := queryRequest.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	//数据校验
	query := model_user.AdminUserQueryToUser(queryRequest)
	err := uc.checkQueryOrUpdateUser(query)
	if err != nil {
		log.Printf("validate %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, err.Error()).Response(api_response.AUTHERR))

		return
	}

	res, err := uc.userService.GetUserList(queryRequest, page, pageSize)

	if err != nil {
		log.Println(fmt.Sprintf("query user %v", err))
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "query user error").Response(api_response.OPERATIONERR))

		return
	}

	result := model_user.UsersToAdminReturnUsers(res)
	log.Println("query users success")
	c.JSON(http.StatusOK, api_response.NewResponse(result, "query users success").Response(api_response.SUCCESS))
}

// UpdateUser 更新用户信息
//
//	@Summary		Update user
//	@Description	Update user information
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		model_user.UpdateUserRequest	true	"Update information"
//	@Success		200		{object}	api_response.ApiResponse{data=nil}	"Update user request success"
//	@Failure		400		{object}	api_response.ApiResponse{data=nil}	"Update user request fail"
//	@Router			/api/user/update    [post]
func (uc *UserController) UpdateUser(c *gin.Context) {
	//判断用户权限
	validity, _ := uc.sessionService.GetSession(c)

	if validity.UserRole != constant.Common && validity.UserRole != constant.Admin {
		log.Printf("you are not login")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "you are not login").Response(api_response.AUTHERR))

		return
	}

	var updateUser model_user.UpdateUserRequest
	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "unmarshal error ").Response(api_response.OPERATIONERR))

		return
	}

	//校验字段
	var old model_user.User
	update := model_user.UpdateUserToUser(old, updateUser)
	err := uc.checkQueryOrUpdateUser(update)
	if err != nil {
		log.Printf("validate %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, err.Error()).Response(api_response.AUTHERR))

		return
	}

	//更新用户信息
	oldInfo, err := uc.userService.GetUser(validity.UserAccount)
	if err != nil {
		log.Println(fmt.Sprintf("query user %v", err))
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "query user error").Response(api_response.OPERATIONERR))

		return
	}
	user := model_user.UpdateUserToUser(oldInfo, updateUser)
	err = uc.userService.UpdateUser(user)
	if err != nil {
		log.Println(fmt.Sprintf("update user %v", err))
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "update user error").Response(api_response.OPERATIONERR))

		return
	}

	log.Printf("update model_user success")
	c.JSON(http.StatusOK, api_response.NewResponse(nil, "update user success").Response(api_response.SUCCESS))
}

// EditUser 更新用户信息
//
//	@Summary		Admin edit user information
//	@Description	Admin edit user information
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		model_user.EditUserRequest	true	"User information"
//	@Success		200		{object}	api_response.ApiResponse{data=nil}	"Edit user request success"
//	@Failure		400		{object}	api_response.ApiResponse{data=nil}	"Edit user request fail"
//	@Router			/api/user/admin/update    [post]
func (uc *UserController) EditUser(c *gin.Context) {
	//判断用户权限
	validity, _ := uc.sessionService.GetSession(c)
	if validity.UserRole != constant.Admin {
		log.Printf("you are not admin")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "you are not admin").Response(api_response.AUTHERR))

		return
	}

	var editUser model_user.EditUserRequest
	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&editUser); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "unmarshal error ").Response(api_response.OPERATIONERR))

		return
	}

	//数据校验
	var old model_user.User
	edit := model_user.EditUserToUser(old, editUser)
	err := uc.checkQueryOrUpdateUser(edit)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, err.Error()).Response(api_response.AUTHERR))

		return
	}

	//获取原始信息
	oldInfo, err := uc.userService.GetUser(editUser.UserAccount)
	if err != nil {
		log.Println(fmt.Sprintf("no such user %v", err))
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "no such user").Response(api_response.OPERATIONERR))
		return
	}

	//更新用户信息
	user := model_user.EditUserToUser(oldInfo, editUser)
	err = uc.userService.UpdateUser(user)
	if err != nil {
		log.Println(fmt.Sprintf("update user %v", err))
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "update user error").Response(api_response.OPERATIONERR))

		return
	}

	log.Printf("update user success")
	c.JSON(http.StatusOK, api_response.NewResponse(nil, "update user success").Response(api_response.SUCCESS))
}

// DeleteUser 删除用户
//
//	@Summary		DeleteUser user
//	@Description	DeleteUser user  by user account
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	path		string						true	"User account"
//	@Success		200		{object}	api_response.ApiResponse{data=nil}	"Delete user success"
//	@Failure		400		{object}	api_response.ApiResponse{data=nil}	"Delete user fail"
//	@Router			/api/user/admin/delete [get]
func (uc *UserController) DeleteUser(c *gin.Context) {

	//判断用户权限
	validity, _ := uc.sessionService.GetSession(c)
	if validity.UserRole != constant.Admin {
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "you are not validity user").Response(api_response.AUTHERR))

		return
	}

	//获取account参数的值
	userAccount := c.Param("account")
	if userAccount == "" {
		log.Println("not a valid query user account")
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, "not a valid query account").Response(api_response.PARAMSERR))

		return
	}

	//逻辑删除用户
	err := uc.userService.DeleteUser(userAccount)
	if err != nil {
		log.Printf("delete user  %v", err)
		c.JSON(http.StatusBadRequest, api_response.NewResponse(nil, err.Error()).Response(api_response.OPERATIONERR))

		return
	}

	log.Printf("delete user success")
	c.JSON(http.StatusOK, api_response.NewResponse(nil, "delete user success").Response(api_response.SUCCESS))
}

func (uc *UserController) checkUser(account string, password string) error {
	if account == "" {
		return errors.New("user account required")
	}
	if password == "" {
		return errors.New("user password required")
	}
	err := uc.checkPassword(password)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserController) checkQueryOrUpdateUser(queryUser model_user.User) error {
	if queryUser.ID != "" && len(queryUser.ID) > 64 {
		return errors.New("invalid id")
	}
	if queryUser.UserAccount != "" && len(queryUser.UserAccount) > 256 {
		return errors.New("invalid user account")
	}
	if queryUser.UserName != "" && len(queryUser.UserName) > 256 {
		return errors.New("invalid user name")
	}
	if queryUser.UserRole != "" {
		if queryUser.UserRole != constant.Common && queryUser.UserRole != constant.Admin {
			return errors.New("invalid user role")
		}
	}
	if queryUser.AvatarUrl != "" && len(queryUser.AvatarUrl) > 256 {
		return errors.New("invalid avatar url")
	}
	if queryUser.IsDelete != 0 {
		if queryUser.IsDelete != constant.ALIVE && queryUser.IsDelete != constant.DELETE {
			return errors.New("invalid is delete")
		}
	}
	if queryUser.UserPassword != "" {
		err := uc.checkPassword(queryUser.UserPassword)
		if err != nil {
			return err
		}
	}
	return nil
}

func (uc *UserController) checkPassword(password string) error {
	if len(password) < 8 || len(password) > 32 {
		return errors.New("invalid password")
	}
	//校验密码合法性
	expr := `^(?![0-9a-zA-Z]+$)(?![a-zA-Z!@#$%^&*]+$)(?![0-9!@#$%^&*]+$)[0-9A-Za-z!@#$%^&*]{8,16}$`
	reg, _ := regexp2.Compile(expr, 0)
	m, _ := reg.MatchString(password)
	if !m {
		return errors.New("invalid password, At least one special character, lowercase and uppercase, is required")
	}
	return nil
}
