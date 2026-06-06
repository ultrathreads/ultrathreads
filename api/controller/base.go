package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ultrathreads/form"
	"ultrathreads/model"
	"ultrathreads/util"
)

type R struct {
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// BaseController controller
type BaseController struct {
}

// BindAndValidate bind and validate
func (c *BaseController) BindAndValidate(ctx *gin.Context, obj interface{}) bool {
	if err := form.Bind(ctx, obj); err != nil {
		c.Fail(ctx, &util.CodeError{Code: -1, Message: err.Error()})
		return false
	}
	return true
}

// GetCurrentUser get current user from contex
func (c *BaseController) GetCurrentUser(ctx *gin.Context) *model.User {
	if currentUser, ok := ctx.Get("CurrentUser"); ok {
		return currentUser.(*model.User)
	}
	return nil
}

// Success output json data
func (c *BaseController) Success(ctx *gin.Context, data interface{}) {
	resp := R{Code: 0, Success: true, Message: "ok", Data: data}

	// 仅 debug 模式使用格式化JSON
	if gin.Mode() == gin.DebugMode {
		ctx.IndentedJSON(http.StatusOK, resp)
	} else {
		ctx.JSON(http.StatusOK, resp)
	}
}

// Fail output error
func (c *BaseController) Fail(ctx *gin.Context, error *util.CodeError) {
	resp := R{Code: error.Code, Success: false, Message: error.Message}

	if gin.Mode() == gin.DebugMode {
		ctx.IndentedJSON(http.StatusOK, resp)
	} else {
		ctx.JSON(http.StatusOK, resp)
	}
	ctx.Abort()
}
