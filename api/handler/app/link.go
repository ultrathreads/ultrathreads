package app

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/handler/base"
	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

type LinkHandler struct {
	base.BaseHandler
	linkSvc service.LinkServicer
}

func NewLinkHandler(linkSvc service.LinkServicer) *LinkHandler {
	return &LinkHandler{linkSvc: linkSvc}
}

// List 列表
func (h *LinkHandler) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)

	links, paging := h.linkSvc.List(querybuilder.NewQueryBuilder().
		Eq("status", model.StatusOk).Page(page, 20).Asc("id"))

	var results []map[string]interface{}
	for _, v := range links {
		results = append(results, h.buildLink(v))
	}
	h.Success(ctx, gin.H{
		"results": results,
		"paging":  paging,
	})
}

// 前10个链接
func (h *LinkHandler) GetToplinks(ctx *gin.Context) {
	links := h.linkSvc.Find(querybuilder.NewQueryBuilder().
		Eq("status", model.StatusOk).Limit(10).Asc("id"))

	var results []map[string]interface{}
	for _, v := range links {
		results = append(results, h.buildLink(v))
	}
	h.Success(ctx, results)
}

func (h *LinkHandler) buildLink(link model.Link) map[string]interface{} {
	return map[string]interface{}{
		"linkId":     link.ID,
		"url":        link.Url,
		"title":      link.Title,
		"summary":    link.Summary,
		"logo":       link.Logo,
		"createTime": link.CreateTime,
	}
}
