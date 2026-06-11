package controller

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/cache"
)

type SiteController struct {
	BaseController
}

func (c *SiteController) Stat(ctx *gin.Context) {
	data := make(map[string]interface{})
	data["userCount"] = cache.StatCache.GetUserCount()
	data["postCount"] = cache.StatCache.GetPostCount()

	c.Success(ctx, data)
}

// Ping 健康检查 - 通过事件总线异步输出 pong
func (c *SiteController) Ping(ctx *gin.Context) {
	c.Success(ctx, gin.H{
		"message": "pong",
	})
}
