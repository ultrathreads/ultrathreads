package cron

import (
	"github.com/robfig/cron"

	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/log"
)

// 将 cron 实例提升为包级变量，供 Stop() 使用
var c *cron.Cron

func Setup() {
	if !util.IsProd() {
		log.Info("Not in a production enviroment!")
		return
	}

	log.Info("Cron setup")
	startSchedule()
}

func startSchedule() {
	c = cron.New()

	// Generate RSS
	addCronFunc(c, "@every 30m", func() {
		service.ArticleService.GenerateRss()
		service.Srv.Post.GenerateRss()
	})

	c.Start()
}

func addCronFunc(c *cron.Cron, spec string, cmd func()) {
	err := c.AddFunc(spec, cmd)
	if err != nil {
		log.Error(err.Error())
	}
}

// Stop 优雅停止定时任务，等待正在执行的任务完成
func Stop() {
	if c == nil {
		return
	}
	log.Info("Stopping cron jobs...")
	c.Stop()
	log.Info("Cron jobs stopped")
}