package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"ultrathreads/dto"
	"ultrathreads/delivery/handler/base"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

// NodeHandler node controller
type NodeHandler struct {
	base.BaseHandler
	nodeSvc service.NodeService
}

func NewNodeHandler(nodeSvc service.NodeService) *NodeHandler {
	return &NodeHandler{nodeSvc: nodeSvc}
}

// Show show node
func (h *NodeHandler) Show(ctx *gin.Context) {
	var gDto dto.IdRequest
	if h.BindAndValidate(ctx, &gDto) {
		node := h.nodeSvc.Get(gDto.ID)
		if node == nil {
			h.Fail(ctx, util.NewErrorMsg("Node not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		h.SuccessWithIncluded(ctx, node)
	}
}

// Store create a node
func (h *NodeHandler) Store(ctx *gin.Context) {
	var nodeForm dto.NodeCreateForm
	if !h.BindAndValidate(ctx, &nodeForm) {
		return
	}
	node, err := h.nodeSvc.Create(nodeForm)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	h.Success(ctx, node)
}

// Update update a node
func (h *NodeHandler) Update(ctx *gin.Context) {
	var nodeForm dto.NodeUpdateForm
	if !h.BindAndValidate(ctx, &nodeForm) {
		return
	}
	node := h.nodeSvc.Get(nodeForm.ID)
	if node == nil {
		h.Fail(ctx, util.NewErrorMsg("Node not found, id="+strconv.FormatInt(nodeForm.ID, 10)))
		return
	}

	err := h.nodeSvc.Update(nodeForm.ID, nodeForm)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	h.Success(ctx, node)
}

func (h *NodeHandler) Sort(ctx *gin.Context) {
	var req dto.SortRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"code": 40001, "message": "参数错误: " + err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"data": gin.H{"updated": len(req.Items)},
	})
}

// Delete delete node
func (h *NodeHandler) Delete(ctx *gin.Context) {
	var gDto dto.IdRequest
	if !h.BindAndValidate(ctx, &gDto) {
		return
	}
	h.nodeSvc.Delete(gDto.ID)
	h.Success(ctx, nil)
}

// List list nodes
func (h *NodeHandler) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	limit := util.FormIntDefault(ctx, "limit", 20)
	id := ctx.Request.FormValue("id")
	name := ctx.Request.FormValue("name")

	conditions := querybuilder.NewQueryBuilder()
	if len(id) > 0 {
		conditions.Eq("id", id)
	}
	if len(name) > 0 {
		conditions.Like("name", name)
	}
	list, paging := h.nodeSvc.List(conditions.Page(page, limit).Asc("sort_no"))
	var results []map[string]interface{}
	for _, node := range list {
		item := util.StructToMap(node)
		results = append(results, item)
	}

	h.SuccessWithIncluded(ctx, &querybuilder.PageResult{Results: results, Page: paging})
}
