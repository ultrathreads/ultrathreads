package app

import (
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"

	"ultrathreads/handler/base"
	"ultrathreads/util/urls"
)

// CaptchaHandler captcha controller
type CaptchaHandler struct {
	base.BaseHandler
}

// GetRequest request captcha id and url
func (h *CaptchaHandler) GetRequest(ctx *gin.Context) {

	captchaID := captcha.NewLen(4)
	captchaURL := urls.AbsUrl("/api/captcha/show/" + captchaID)

	data := make(map[string]interface{})
	data["captchaId"] = captchaID
	data["captchaUrl"] = captchaURL

	h.Success(ctx, data)
}

// Show show captcha image
func (h *CaptchaHandler) Show(ctx *gin.Context) {
	captchaID := ctx.Param("captchaId")
	if captchaID == "" {
		return
	}
	if !captcha.Reload(captchaID) {
		return
	}

	ctx.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Writer.Header().Set("Pragma", "no-cache")
	ctx.Writer.Header().Set("Expires", "10")
	ctx.Writer.Header().Set("Content-Type", "image/png")
	captcha.WriteImage(ctx.Writer, captchaID, captcha.StdWidth, captcha.StdHeight)
}
