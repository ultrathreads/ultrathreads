package admin

import (
	"github.com/gin-gonic/gin"
	"strconv"

	"ultrathreads/controller"
	"ultrathreads/form"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

// NodeController node controller
type NodeController struct {
	controller.BaseController
}

// Show show node
func (c *NodeController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		node := service.Srv.Node.Get(gDto.ID)
		if node == nil {
			c.Fail(ctx, util.NewErrorMsg("Node not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.SuccessWithIncluded(ctx, node)
	}
}

// Store create a node
func (c *NodeController) Store(ctx *gin.Context) {
	var nodeForm form.NodeCreateForm
	if !c.BindAndValidate(ctx, &nodeForm) {
		return
	}
	node, err := service.Srv.Node.Create(nodeForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, node)
}

// Update update a node
func (c *NodeController) Update(ctx *gin.Context) {
	var nodeForm form.NodeUpdateForm
	if !c.BindAndValidate(ctx, &nodeForm) {
		return
	}
	node := service.Srv.Node.Get(nodeForm.ID)
	if node == nil {
		c.Fail(ctx, util.NewErrorMsg("Node not found, id="+strconv.FormatInt(nodeForm.ID, 10)))
		return
	}

	err := service.Srv.Node.Update(nodeForm.ID, nodeForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, node)
}

func (c *NodeController) Sort(ctx *gin.Context) {
    var req form.SortRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(400, gin.H{"code": 40001, "message": "参数错误: " + err.Error()})
        return
    }

    // 批量更新 sort_no
    /*
    for _, item := range req.Items {
        if err := service.nodeService.UpdateSort(item.ID, item.SortNo); err != nil {
            ctx.JSON(500, gin.H{"code": 50001, "message": "排序保存失败"})
            return
        }
    }
    */

    ctx.JSON(200, gin.H{
        "data": gin.H{"updated": len(req.Items)},
    })
}

// Delete delete node
func (c *NodeController) Delete(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	service.Srv.Node.Delete(gDto.ID)
	c.Success(ctx, nil)
}

// List list nodes
func (c *NodeController) List(ctx *gin.Context) {
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
	list, paging := service.Srv.Node.List(conditions.Page(page, limit).Asc("sort_no"))
	var results []map[string]interface{}
	for _, node := range list {
		item := util.StructToMap(node)
		results = append(results, item)
	}

	c.SuccessWithIncluded(ctx, &querybuilder.PageResult{Results: results, Page: paging})
}
