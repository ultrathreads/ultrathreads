package controller

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/converter"
	"ultrathreads/cache"
	"ultrathreads/bus/event"
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
	nodeId := util.ParamInt64Default(ctx, "id", 0)

	user := c.GetCurrentUser(ctx)
	now := util.NowTimestamp()
	if err := service.UserReadStateService.MarkAsRead(user.ID, nodeId, now); err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}

	c.Success(ctx, nil)
}

func (c *NodeController) ViewPost(ctx *gin.Context) {
	nodeId := util.ParamInt64Default(ctx, "id", 0)
	postId := util.QueryInt64Default(ctx, "postId", 0)
	user := c.GetCurrentUser(ctx)

    c.PublishEvent(ctx, event.PostViewed{
        UserID:     user.ID,
        NodeID:     nodeId,
        PostID:     postId,
        ViewedTime: util.NowTimestamp(),
    })

	c.Success(ctx, nil)
}