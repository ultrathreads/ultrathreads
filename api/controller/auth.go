package controller

import (
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"

	"ultrathreads/render"
	"ultrathreads/form"
	"ultrathreads/service"
	"ultrathreads/util"
)

// AuthController auth controller
type AuthController struct {
	BaseController
}

// Register
func (c *AuthController) Register(ctx *gin.Context) {
	ref := ctx.Request.FormValue("ref")
	var dto form.AuthSignupForm
	if c.BindAndValidate(ctx, &dto) {
		if !captcha.VerifyString(dto.CaptchaID, dto.CaptchaCode) {
			c.Fail(ctx, util.ErrorCaptchaWrong)
			return
		}

		user, err := service.Srv.User.Create(dto.Username, dto.Email, dto.Nickname, dto.Password, dto.RePassword)
		if err != nil {
			c.Fail(ctx, util.FromError(err))
			return
		}
		c.Success(ctx, gin.H{
			"user": render.ToUser(user),
			"ref":  ref,
		})
	}
}
