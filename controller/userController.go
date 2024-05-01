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
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Register 用户注册
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
	_, err := uc.userService.QueryUser(user)
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
	err = uc.userService.AddUser(user)
	if err != nil {
		log.Printf("create user failed")
		return
	}

	var resultUsers []*model.ResultUser
	resultUsers[0] = model.UserProc(user)
	//插入成功
	c.JSON(http.StatusOK, utils.Success(resultUsers, "register success"))
	log.Printf("register success")
}

// Login 用户登录
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
	re, err := uc.userService.QueryUser(user)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("The user has not registered yet")
		c.JSON(http.StatusBadRequest, utils.Error(utils.AUTHERR, "The user has not registered yet"))
		return
	}

	if re.UserPassword != user.UserPassword {
		log.Println("username or password is wrong")
		c.JSON(http.StatusBadRequest, utils.Error(utils.AUTHERR, "username or password is wrong"))
		return
	}

	//登录成功, 存储session信息
	userSession := model.UserSession{
		ID:     re.ID,
		Expire: time.Now().Add(time.Hour * 24).Unix(),
		Role:   re.UserRole,
	}

	err = uc.userService.NewSession(c, userSession)
	if err != nil {
		log.Println(fmt.Sprintf("session create %v", err))
		return
	}

	var resultUsers []*model.ResultUser
	resultUsers[0] = model.UserProc(re)
	log.Printf("login success")
	c.JSON(http.StatusOK, utils.Success(resultUsers, "login success"))

}

// Logout 登出账户
func (uc *UserController) Logout(c *gin.Context) {

	//判断用户是否登录
	validity := uc.checkValidity(c)
	if validity == ANONYMOUS {
		log.Printf("you must login first")
		c.JSON(http.StatusBadRequest, utils.Error(utils.LOGINERR, "you must login first"))
		return
	}

	//删除session
	err := uc.userService.DeleteSession(c)
	if err != nil {
		log.Println(fmt.Sprintf("session delete %v", err))
		return
	}

	log.Printf("logout success")
	c.JSON(http.StatusOK, utils.Success(nil, "logout success"))
}

// QueryUsers 查询用户
func (uc *UserController) QueryUsers(c *gin.Context) {

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
		log.Println("username is empty")
		c.JSON(http.StatusBadRequest, utils.Error(utils.LOGINERR, "username is empty"))
		return
	}
	user.UserName = username

	results, err := uc.userService.QueryUsers(user)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("No query users found")
		c.JSON(http.StatusBadRequest, utils.Error(utils.LOGINERR, "No query users found"))
		return
	}

	//查询结果处理
	rets := model.UsersProc(results)

	log.Printf("query user success")
	c.JSON(http.StatusOK, utils.Success(rets, "query user success"))

}

// DeleteUser 删除用户
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
	err := uc.userService.LogicDeleteUser(user)
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

	userSession := uc.userService.GetSession(c)
	if userSession.Role == ADMIN {
		return ADMIN
	}

	if userSession.Role == LOGIN {
		return LOGIN
	}
	return ANONYMOUS
}
