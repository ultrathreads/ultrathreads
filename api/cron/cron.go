package cron

import (
	"github.com/robfig/cron"

	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/log"
)

var c *cron.Cron

func Setup(articleSvc service.ArticleService, postSvc service.PostService) {
	if !util.IsProd() {
		log.Info("Not in a production enviroment!")
		return
	}

	log.Info("Cron setup")
	startSchedule(articleSvc, postSvc)
}

func startSchedule(articleSvc service.ArticleService, postSvc service.PostService) {
	c = cron.New()

	addCronFunc(c, "@every 30m", func() {
		articleSvc.GenerateRss()
		postSvc.GenerateRss()
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