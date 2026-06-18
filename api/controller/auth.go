package controller

import (
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"

	"ultrathreads/render"
	"ultrathreads/dto"
	"ultrathreads/service"
	"ultrathreads/util"
)

// AuthController auth controller
type AuthController struct {
	BaseController
	userSvc service.UserServicer
}

func NewAuthController(userSvc service.UserServicer) *AuthController {
	return &AuthController{userSvc: userSvc}
}

// Register
func (c *AuthController) Register(ctx *gin.Context) {
	ref := ctx.Request.FormValue("ref")
	var req dto.AuthSignupForm
	if c.BindAndValidate(ctx, &req) {
		if !captcha.VerifyString(req.CaptchaID, req.CaptchaCode) {
			c.Fail(ctx, util.ErrorCaptchaWrong)
			return
		}

		user, err := c.userSvc.Create(req.Username, req.Email, req.Nickname, req.Password, req.RePassword)
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