package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	evbus "github.com/asaskevich/EventBus"

	"ultrathreads/bus"
	"ultrathreads/bus/core"
	"ultrathreads/form"
	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/log"
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

// PublishEvent 统一的事件发布入口（已封装异步与异常防御）
func (c *BaseController) PublishEvent(ctx *gin.Context, payload interface{}) {
    // 1. 提取必要的依赖（在同步阶段完成，避免异步后获取失败）
    busCtx := ctx.MustGet(bus.BusKey).(evbus.Bus)
    
    // 2. 启动异步协程处理事件
    go func() {
        // 防御性编程：捕获异步协程中的 panic，防止整个进程崩溃
        defer func() {
            if r := recover(); r != nil {
                log.Error("Async event publish panic recovered: %v", r)
            }
        }()
        
        // 3. 执行核心分发逻辑
        core.PublishTyped(busCtx, payload)
    }()
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

// GetLastReadStates 获取当前用户的已读状态
// - 不传参：返回全量 map（nodeId=0 场景）
// - 传一个或多个 nodeID：仅返回指定节点的已读状态，自动过滤零值
func (c *BaseController) GetLastReadStates(ctx *gin.Context, nodeIDs ...int64) map[int64]int64 {
	val, exists := ctx.Get("CurrentUserReadStates")
	if !exists || val == nil {
		return make(map[int64]int64)
	}
	states, ok := val.(map[int64]int64)
	if !ok {
		return make(map[int64]int64)
	}

	// 无参数 → 全量返回（只读场景直接返回引用，避免拷贝开销）
	if len(nodeIDs) == 0 {
		return states
	}

	// 有参数 → 按需提取指定节点
	result := make(map[int64]int64, len(nodeIDs))
	for _, id := range nodeIDs {
		if ts := states[id]; ts > 0 {
			result[id] = ts
		}
	}
	return result
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
