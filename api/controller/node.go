package controller

import (
	"strconv"
	"github.com/gin-gonic/gin"

	"ultrathreads/render"
	"ultrathreads/form"
	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
)

type NodeController struct {
	BaseController
	nodeSvc service.NodeServicer 
}

func NewNodeController(svc service.NodeServicer) *NodeController {
    return &NodeController{nodeSvc: svc}
}

// List 节点列表
func (c *NodeController) List(ctx *gin.Context) {
	nodes := c.nodeSvc.GetNodes()

	c.Success(ctx, render.ToNodes(nodes))
}

// Show 显示单个节点
func (c *NodeController) Show(ctx *gin.Context) {
	var gDto form.IdentifierDto
	if c.BindAndValidate(ctx, &gDto) {
		var node *model.Node
		if id, parseErr := strconv.ParseInt(gDto.Slug, 10, 64); parseErr == nil {
			node = c.nodeSvc.Get(id)
		} else {
			node = c.nodeSvc.GetBySlug(gDto.Slug)
		}

		c.Success(ctx, render.ToNode(node))
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
