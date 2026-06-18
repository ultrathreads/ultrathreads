package controller

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

type LinkController struct {
	BaseController
	linkSvc service.LinkServicer
}

func NewLinkController(linkSvc service.LinkServicer) *LinkController {
	return &LinkController{linkSvc: linkSvc}
}

// List 列表
func (c *LinkController) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)

	links, paging := c.linkSvc.List(querybuilder.NewQueryBuilder().
		Eq("status", model.StatusOk).Page(page, 20).Asc("id"))

	var results []map[string]interface{}
	for _, v := range links {
		results = append(results, c.buildLink(v))
	}
	c.Success(ctx, gin.H{
		"results": results,
		"paging":  paging,
	})
}

// 前10个链接
func (c *LinkController) GetToplinks(ctx *gin.Context) {
	links := c.linkSvc.Find(querybuilder.NewQueryBuilder().
		Eq("status", model.StatusOk).Limit(10).Asc("id"))

	var results []map[string]interface{}
	for _, v := range links {
		results = append(results, c.buildLink(v))
	}
	c.Success(ctx, results)
}

func (c *LinkController) buildLink(link model.Link) map[string]interface{} {
	return map[string]interface{}{
		"linkId":     link.ID,
		"url":        link.Url,
		"title":      link.Title,
		"summary":    link.Summary,
		"logo":       link.Logo,
		"createTime": link.CreateTime,
	}
}