package service

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/userManageSystem/entity/modeluser"
	mysql2 "github.com/xissg/userManageSystem/service/mysql"
	redis2 "github.com/xissg/userManageSystem/service/redis"
	"log"
	"time"
)

type CronJob struct {
	mysqlService *mysql2.UserService
	redisService *redis2.UserService
}

func NewCronJob(m *mysql2.UserService, r *redis2.UserService) *CronJob {

	return &CronJob{
		mysqlService: m,
		redisService: r,
	}
}

// Start 定时任务，开启一个time ticker, 当time ticker为0是开始执行job
// 任务执行结束发送一个done
func (cron *CronJob) Start() {
	ctx := gin.Context{}
	var queryModel modeluser.AdminUserQueryRequest
	res, err := cron.mysqlService.GetUserList(queryModel)
	if err != nil {
		log.Println("mysql read error:", err)
		return
	}

	err = cron.redisService.AddUsers(res, &ctx)
	if err != nil {
		log.Println("redis write error:", err)
	}
	log.Println("job done", time.Now().String())
}
