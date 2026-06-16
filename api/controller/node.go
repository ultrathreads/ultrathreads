package controller

import (
	"strconv"
	"github.com/gin-gonic/gin"

	"ultrathreads/converter"
	"ultrathreads/cache"
	"ultrathreads/form"
	"ultrathreads/model"
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
	var gDto form.IdentifierDto
	if c.BindAndValidate(ctx, &gDto) {
		var node *model.Node
		if id, parseErr := strconv.ParseInt(gDto.Slug, 10, 64); parseErr == nil {
			node = service.Srv.NodeService.Get(id)
		} else {
			node = service.Srv.NodeService.GetBySlug(gDto.Slug)
		}

		c.Success(ctx, converter.ToNode(node))
	}
}

func (c *NodeController) MarkAsRead(ctx *gin.Context) {
	var gDto form.IdentifierDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}

	user := c.GetCurrentUser(ctx)
	now := util.NowTimestamp()
	if err := service.UserReadStateService.MarkAsRead(user.ID, gDto.Slug, now); err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}

	c.Success(ctx, nil)
}
