package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	evbus "github.com/asaskevich/EventBus"

	"ultrathreads/events"
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

// GetBus 从请求上下文中获取事件总线
func (c *BaseController) GetBus(ctx *gin.Context) evbus.Bus {
	return ctx.MustGet(events.BusKey).(evbus.Bus)
}

// BindAndValidate bind and validate
func (c *BaseController) BindAndValidate(ctx *gin.Context, obj interface{}) bool {
	if err := form.Bind(ctx, obj); err != nil {
		c.Fail(ctx, &util.CodeError{Code: -1, Message: err.Error()})
		return false
	}
	return true
}

// BindAndValidateUri 专门处理 URL Path 参数绑定
func (c *BaseController) BindAndValidateUri(ctx *gin.Context, obj interface{}) bool {
    if err := ctx.ShouldBindUri(obj); err != nil {
        c.Fail(ctx, util.FromError(err))
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

func (c *BaseController) GetLastReadAt(ctx *gin.Context, nodeID int) int64 {
    val, exists := ctx.Get(util.ReadStateKey(nodeID))
    if !exists {
        return 0
    }
    ts, ok := val.(int64)
    if !ok {
        return 0
    }
    return ts
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
