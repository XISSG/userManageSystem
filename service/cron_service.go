package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type CronJob struct {
	mysqlService DBService
	redisService CacheService
}

func NewCronJob(m DBService, r CacheService) *CronJob {

	return &CronJob{
		mysqlService: m,
		redisService: r,
	}
}

// 定时任务，开启一个time ticker, 当time ticker为0是开始执行job
// 任务执行结束发送一个done
func (cron *CronJob) Start() {

	names := []string{
		"xissg",
		"111111",
	}
	ctx := gin.Context{}
	for _, name := range names {
		res, err := cron.mysqlService.GetUserByName(name)
		if err != nil {
			log.Println("mysql read error:", err)
			continue
		}
		err = cron.redisService.AddUser(res, &ctx)
		if err != nil {
			log.Println("redis write error:", err)
		}
	}

	log.Println("job done", time.Now().String())
}
