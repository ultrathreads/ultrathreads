package app

import (
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"

	"ultrathreads/dto"
	"ultrathreads/delivery/handler/base"
	"ultrathreads/render"
	"ultrathreads/service"
	"ultrathreads/util"
)

// AuthHandler auth controller
type AuthHandler struct {
	base.BaseHandler
	userSvc service.UserServicer
}

func NewAuthHandler(userSvc service.UserServicer) *AuthHandler {
	return &AuthHandler{userSvc: userSvc}
}

// Register
func (h *AuthHandler) Register(ctx *gin.Context) {
	ref := ctx.Request.FormValue("ref")
	var req dto.AuthSignupForm
	if h.BindAndValidate(ctx, &req) {
		if !captcha.VerifyString(req.CaptchaID, req.CaptchaCode) {
			h.Fail(ctx, util.ErrorCaptchaWrong)
			return
		}

		user, err := h.userSvc.Create(req.Username, req.Email, req.Nickname, req.Password, req.RePassword)
		if err != nil {
			h.Fail(ctx, util.FromError(err))
			return
		}
		h.Success(ctx, gin.H{
			"user": render.ToUser(user),
			"ref":  ref,
		})
	}
}
