package app

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"ultrathreads/dto"
	"ultrathreads/handler/base"
	"ultrathreads/model"
	"ultrathreads/render"
	"ultrathreads/service"
	"ultrathreads/util"
)

type NodeHandler struct {
	base.BaseHandler
	nodeSvc          service.NodeServicer
	userReadStateSvc service.UserReadStateServicer
}

func NewNodeHandler(svc service.NodeServicer, userReadStateSvc service.UserReadStateServicer) *NodeHandler {
	return &NodeHandler{nodeSvc: svc, userReadStateSvc: userReadStateSvc}
}

// List 节点列表
func (h *NodeHandler) List(ctx *gin.Context) {
	nodes := h.nodeSvc.GetNodes()

	h.Success(ctx, render.ToNodes(nodes))
}

// Show 显示单个节点
func (h *NodeHandler) Show(ctx *gin.Context) {
	var gDto dto.SlugRequest
	if h.BindAndValidate(ctx, &gDto) {
		var node *model.Node
		if id, parseErr := strconv.ParseInt(gDto.Slug, 10, 64); parseErr == nil {
			node = h.nodeSvc.Get(id)
		} else {
			node = h.nodeSvc.GetBySlug(gDto.Slug)
		}

		h.Success(ctx, render.ToNode(node))
	}
}

func (h *NodeHandler) MarkAsRead(ctx *gin.Context) {
	var gDto dto.SlugRequest
	if !h.BindAndValidate(ctx, &gDto) {
		return
	}

	user := h.GetCurrentUser(ctx)
	now := util.NowTimestamp()
	if err := h.userReadStateSvc.MarkAsRead(user.ID, gDto.Slug, now); err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}

	h.Success(ctx, nil)
}
