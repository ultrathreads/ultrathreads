package controller

import (
	"time"
	"github.com/gin-gonic/gin"

	"ultrathreads/converter"
	"ultrathreads/cache"
	"ultrathreads/form"
	"ultrathreads/service"
	"ultrathreads/util"
)

type NodeController struct {
	BaseController
}

// List 节点列表
func (c *NodeController) List(ctx *gin.Context) {
	nodes := cache.NodeCache.GetAll()

	c.Success(ctx, converter.ToNodes(nodes))
}

// Show 显示单个节点
func (c *NodeController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		node := service.NodeService.Get(gDto.ID)
		c.Success(ctx, converter.ToNode(node))
	}
}

func (c *NodeController) MarkAsRead(ctx *gin.Context) {
	nodeId := form.ParamInt64Default(ctx, "id", 0)

    user := c.GetCurrentUser(ctx)
    if err := service.UserReadStateService.MarkAsRead(user.ID, nodeId, time.Now().Unix()); err != nil {
        c.Fail(ctx, util.FromError(err))
        return
    }
    
    c.Success(ctx, nil)
}
