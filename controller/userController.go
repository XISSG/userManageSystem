package controller

import (
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/model"
	"github.com/xissg/userManageSystem/service"
	"github.com/xissg/userManageSystem/utils"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

// 用户的登录状态
const (
	ANONYMOUS = iota
	LOGIN
	ADMIN //user_role 字段判断是否为admin用户
)

type UserController struct {
	sessionService service.SessionService
	rwdService     *service.RWDService
}

func NewUserController(userService service.DBService, redisService service.CacheService, sessionService service.SessionService) *UserController {
	return &UserController{
		sessionService: sessionService,
		rwdService:     service.NewRWDService(userService, redisService),
	}
}

// Register 用户注册
// @Summary User registration
// @Description Register a new user
// @Tags User
// @Accept json
// @Produce json
// @Param user body string true "User object"
// @Success 200 {object} utils.ApiResponse "Success registered"
// @Failure 400 {object} utils.ApiResponse "Registration failed"
// @Router /user/register [post]
func (uc *UserController) Register(c *gin.Context) {
	var user model.User

	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		return
	}

	//校验用户名长度
	if len(user.UserName) < 4 || len(user.UserName) > 32 {
		c.JSON(http.StatusBadRequest, utils.Error(utils.LOGINERR, "user name too short or too long"))
		log.Printf("wrong user name length")
		return
	}

	//校验用户名是否存在
	_, err := uc.rwdService.Read(user, c)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, utils.Error(utils.LOGINERR, "username repeated"))
		log.Println("username repeated")
		return
	}

	//校验密码长度
	if len(user.UserPassword) < 8 || len(user.UserPassword) > 32 {
		c.JSON(http.StatusBadRequest, utils.Error(utils.LOGINERR, "wrong password length"))
		log.Printf("wrong password length")
		return
	}

	//校验密码合法性
	expr := `^(?![0-9a-zA-Z]+$)(?![a-zA-Z!@#$%^&*]+$)(?![0-9!@#$%^&*]+$)[0-9A-Za-z!@#$%^&*]{8,16}$`
	reg, _ := regexp2.Compile(expr, 0)
	m, _ := reg.MatchString(user.UserPassword)
	if !m {
		c.JSON(http.StatusBadRequest, utils.Error(utils.LOGINERR, "Invalid password, At least one special character, lowercase and uppercase, is required"))
		log.Printf("illegal user password")
		return
	}

	//用户密码脱敏
	user.UserPassword = utils.MD5Crypt(user.UserPassword)

	//生成唯一id
	id, err := utils.IdGenerator()
	if err != nil {
		log.Println(fmt.Sprintf("id generator %v", err))
		return
	}
	user.ID = id.Int64()

	//生成创建时间和更新时间
	user.CreateTime = time.Now().UTC()
	user.UpdateTime = time.Now().UTC()

	//插入数据库
	err = uc.rwdService.Add(user, c)
	if err != nil {
		log.Printf("create user failed")
		return
	}

	var resultUser *model.ResultUser
	resultUser = model.UserProc(user)
	//插入成功
	c.JSON(http.StatusOK, utils.Success(resultUser, "register success"))
	log.Printf("register success")
}

// Login 用户登录
// @Summary User login
// @Description Authenticate user login
// @Tags User
// @Accept json
// @Produce json
// @Param user body string true "User object"
// @Success 200 {object} utils.ApiResponse "Login successful"
// @Failure 400 {object} utils.ApiResponse "Login failed"
// @Router /user/login [post]
func (uc *UserController) Login(c *gin.Context) {
	var user model.User

	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		return
	}

	//判断用户名和密码是否为空
	if user.UserName == "" || user.UserPassword == "" {
		log.Printf("username  or password  is empty")
		c.JSON(http.StatusBadRequest, utils.Error(utils.LOGINERR, "username or password is empty"))
		return
	}

	///用户密码脱敏
	user.UserPassword = utils.MD5Crypt(user.UserPassword)

	//查询用户用户名和密码是否匹配
	ret, err := uc.rwdService.Read(user, c)
	if ret == nil {
		log.Println("The user has not registered yet")
		c.JSON(http.StatusBadRequest, utils.Error(utils.AUTHERR, "The user has not registered yet"))
		return
	}
	re := ret.(model.User)
	if re.UserPassword != user.UserPassword {
		log.Println("username or password is wrong")
		c.JSON(http.StatusBadRequest, utils.Error(utils.AUTHERR, "username or password is wrong"))
		return
	}

	//登录成功, 存储session信息
	userSession := model.UserSession{
		ID:       re.ID,
		UserName: re.UserName,
		Role:     re.UserRole,
		Tags:     re.Tags,
	}
	err = uc.sessionService.NewOrUpdateSession(c, userSession)
	if err != nil {
		log.Println(fmt.Sprintf("session create %v", err))
		return
	}

	var resultUsers *model.ResultUser
	resultUsers = model.UserProc(re)
	log.Printf("login success")
	c.JSON(http.StatusOK, utils.Success(resultUsers, "login success"))
}

// Logout 登出账户
// @Summary User logout
// @Description User logout
// @Tags User
// @Accept json
// @Produce json
// @Param user body string true "User object"
// @Success 200 {object} utils.ApiResponse "Logout successful"
// @Failure 400 {object} utils.ApiResponse "Logout failed"
// @Router /user/logout [get]
func (uc *UserController) Logout(c *gin.Context) {

	//判断用户是否登录
	validity := uc.checkValidity(c)
	if validity == ANONYMOUS {
		log.Printf("you must login first")
		c.JSON(http.StatusBadRequest, utils.Error(utils.LOGINERR, "you must login first"))
		return
	}

	//删除session
	err := uc.sessionService.DeleteSession(c)
	if err != nil {
		log.Println(fmt.Sprintf("session delete %v", err))
		return
	}

	log.Printf("logout success")
	c.JSON(http.StatusOK, utils.Success(nil, "logout success"))
}

// QueryUser 查询用户
// @Summary Query user by username
// @Description Get user information by username
// @Tags User
// @Accept  json
// @Produce  json
// @Param username query string true "Username"
// @Success 200 {object} utils.ApiResponse "Query successful"
// @Failure 400 {object} utils.ApiResponse "Query failed"
// @Router /user/admin/query [get]
func (uc *UserController) QueryUser(c *gin.Context) {
	//判断用户权限
	validity := uc.checkValidity(c)

	if validity != ADMIN {
		c.JSON(http.StatusBadRequest, utils.Error(utils.AUTHERR, "you are not admin"))
		return
	}

	//获取username参数的值
	var user model.User

	username := c.Param("username")
	if username == "" {
		log.Println("not a valid query username")
		c.JSON(http.StatusBadRequest, utils.Error(utils.PARAMSERR, "not a valid query username"))
		return
	}

	user.UserName = username
	res, err := uc.rwdService.Read(user, c)
	if err != nil {
		return
	}
	ret := model.UserProc(res.(model.User))
	log.Println("query user success")
	c.JSON(http.StatusOK, utils.Success(ret, "query user success"))

}

// UpdateUser 更新用户信息
// @Summary User update
// @Description Update user information
// @Tags User
// @Accept json
// @Produce json
// @Param user body string true "User object"
// @Success 200 {object} utils.ApiResponse "Update successful"
// @Failure 400 {object} utils.ApiResponse "Update failed"
// @Router /user/update [post]
func (uc *UserController) UpdateUser(c *gin.Context) {
	//判断用户权限
	validity := uc.checkValidity(c)

	if validity != ADMIN {
		c.JSON(http.StatusBadRequest, utils.Error(utils.AUTHERR, "you are not admin"))
		return
	}

	var user model.User
	//反序列化取出JSON数据
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("JSON unmarshal  %v", err)
		return
	}

	usersession := uc.sessionService.GetSession(c)
	user.UserName = usersession.UserName
	err := uc.rwdService.Update(user, c)
	if err != nil {
		log.Printf("user info update  %v", err)
		return
	}
	log.Printf("update user success")
	c.JSON(http.StatusOK, utils.Success(nil, "update user success"))
}

// DeleteUser 删除用户
// @Summary Delete user by username
// @Description Delete user information by username
// @Tags User
// @Accept  json
// @Produce  json
// @Param username query string true "Username"
// @Success 200 {object} utils.ApiResponse "Delete successful"
// @Failure 400 {object} utils.ApiResponse "Delete failed"
// @Router /user/admin/delete [get]
func (uc *UserController) DeleteUser(c *gin.Context) {

	//判断用户权限
	validity := uc.checkValidity(c)
	if validity != ADMIN {
		c.JSON(http.StatusBadRequest, utils.Error(utils.AUTHERR, "you are not validity user"))
		return
	}

	//获取username参数的值
	var user model.User

	username := c.Param("username")
	if username == "" {
		log.Println("not a valid query username")
		c.JSON(http.StatusBadRequest, utils.Error(utils.PARAMSERR, "not a valid query username"))
		return
	}

	user.UserName = username

	//逻辑删除用户
	err := uc.rwdService.Delete(user, c)
	if err != nil {
		log.Printf("delete user  %v", err)
		c.JSON(http.StatusBadRequest, utils.Error(utils.OPERATIONERR, err.Error()))
		return
	}

	log.Printf("delete user success")
	c.JSON(http.StatusOK, utils.Success(nil, "delete user success"))
}

// 校验用户是否登录以及是否为admin用户
func (uc *UserController) checkValidity(c *gin.Context) int {

	userSession := uc.sessionService.GetSession(c)
	if userSession.Role == ADMIN {
		return ADMIN
	}

	if userSession.Role == LOGIN {
		return LOGIN
	}
	return ANONYMOUS
}

// MatchUsers  查询用户
//func (uc *UserController) MatchUsers(c *gin.Context) {
//
//	//判断用户权限
//	validity := uc.checkValidity(c)
//	if validity != LOGIN {
//		c.JSON(http.StatusBadRequest, utils.Error(utils.AUTHERR, "you must login first"))
//		return
//	}
//
//	var user model.User
//
//	//反序列化取出JSON数据
//	if err := c.ShouldBindJSON(&user); err != nil {
//		log.Printf("JSON unmarshal  %v", err)
//		return
//	}
//
//	if user.Tags == "" {
//		log.Printf("user tags is nil")
//		c.JSON(http.StatusBadRequest, utils.Error(utils.PARAMSERR, "user tags is nil"))
//		return
//	}
//	results, err := uc.userService.GetUsersByTags(user.Tags)
//	if errors.Is(err, gorm.ErrRecordNotFound) {
//		log.Println("No query users found")
//		c.JSON(http.StatusBadRequest, utils.Error(utils.LOGINERR, "No match users"))
//		return
//	}
//
//	//查询结果处理
//	rets := model.UsersProc(results)
//
//	log.Printf("query user success")
//	c.JSON(http.StatusOK, utils.Success(rets, "query user success"))
//
//}
