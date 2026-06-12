package controller

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/oauth/gitee"
	"ultrathreads/oauth/github"
	"ultrathreads/oauth/qq"
	"ultrathreads/util"
)

// OAuthController oauth controller
type OAuthController struct {
	BaseController
}

// Authorize authorize
func (c *OAuthController) Authorize(ctx *gin.Context) {
	ref := util.FormStringDefault(ctx, "ref","")
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

	c.Success(ctx, gin.H{
		"url": url,
	})
}
