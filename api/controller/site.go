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
