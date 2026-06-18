package app

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/delivery/handler/base"
	"ultrathreads/oauth/gitee"
	"ultrathreads/oauth/github"
	"ultrathreads/oauth/qq"
	"ultrathreads/util"
)

// OAuthHandler oauth controller
type OAuthHandler struct {
	base.BaseHandler
}

// Authorize authorize
func (h *OAuthHandler) Authorize(ctx *gin.Context) {
	ref := util.FormStringDefault(ctx, "ref", "")
	provider := ctx.Param("provider")
	params := map[string]string{"ref": ref}
	var url string
	if provider == "github" {
		url = github.AuthCodeURL(params)
	} else if provider == "gitee" {
		url = gitee.AuthCodeURL(params)
	} else {
		url = qq.AuthorizeUrl(params)
	}

	h.Success(ctx, gin.H{
		"url": url,
	})
}
