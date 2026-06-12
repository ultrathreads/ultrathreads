package controller

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/cache"
	"ultrathreads/model"
	"ultrathreads/util/hashid"
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
	userHashID, _ := hashid.Encode[model.Node](123)

	decodeStr,_ := hashid.Decode[model.Node]("Xrv1ZBEV")
	c.Success(ctx, gin.H{
		"message": "pong",
		"userHashID": userHashID,
		"decodeStr": decodeStr,
	})
}
