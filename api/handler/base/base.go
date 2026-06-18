package base

import (
	"net/http"

	evbus "github.com/asaskevich/EventBus"
	"github.com/gin-gonic/gin"

	"ultrathreads/bus"
	"ultrathreads/bus/core"
	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/binding"
	"ultrathreads/util/hashid"
	"ultrathreads/util/log"
)

type SR struct {
	Data     interface{} `json:"data,omitempty"`
	Included interface{} `json:"included,omitempty"` // 新增字段，nil时自动忽略
}

type FR struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// BaseHandler controller
type BaseHandler struct {
}

// PublishEvent 统一的事件发布入口（已封装异步与异常防御）
func (h *BaseHandler) PublishEvent(ctx *gin.Context, payload interface{}) {
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
func (h *BaseHandler) BindAndValidate(ctx *gin.Context, obj interface{}) bool {
	if err := binding.Bind(ctx, obj); err != nil {
		h.Fail(ctx, &util.CodeError{Code: -1, Message: err.Error()})
		return false
	}
	return true
}

// BindAndValidateUri 专门处理 URL Path 参数绑定
func (h *BaseHandler) BindAndValidateUri(ctx *gin.Context, obj interface{}) bool {
	if err := ctx.ShouldBindUri(obj); err != nil {
		h.Fail(ctx, util.FromError(err))
		return false
	}
	return true
}

// GetCurrentUser get current user from contex
func (h *BaseHandler) GetCurrentUser(ctx *gin.Context) *model.User {
	if currentUser, ok := ctx.Get("CurrentUser"); ok {
		return currentUser.(*model.User)
	}
	return nil
}

func (h *BaseHandler) GetLastReadStates(ctx *gin.Context, nodeSlugs ...string) map[string]int64 {
	empty := make(map[string]int64)

	val, exists := ctx.Get("CurrentUserReadStates")
	if !exists || val == nil {
		return empty
	}
	states, ok := val.(map[int64]int64)
	if !ok {
		return empty
	}

	// 无参数 → 全量返回（key 从 ID 转为 Slug）
	if len(nodeSlugs) == 0 {
		result := make(map[string]int64, len(states))
		for nodeID, lastReadAt := range states {
			slug := hashid.Id2Slug[model.Node](nodeID)
			result[slug] = lastReadAt
		}
		return result
	}

	// 有参数 → 构建反向索引，按指定 slug 按需提取
	slugToID := make(map[string]int64, len(states))
	for nodeID := range states {
		slug := hashid.Id2Slug[model.Node](nodeID)
		slugToID[slug] = nodeID
	}

	result := make(map[string]int64, len(nodeSlugs))
	for _, slug := range nodeSlugs {
		if nodeID, ok := slugToID[slug]; ok {
			if ts := states[nodeID]; ts > 0 {
				result[slug] = ts
			}
		}
	}
	return result
}

// Success output json data
func (h *BaseHandler) Success(ctx *gin.Context, data interface{}) {
	resp := SR{Data: data}
	ctx.JSON(http.StatusOK, resp)
}

func (h *BaseHandler) SuccessWithIncluded(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func (h *BaseHandler) RespondOK(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

// Fail output error
func (h *BaseHandler) Fail(ctx *gin.Context, error *util.CodeError) {
	resp := FR{Code: error.Code, Message: error.Message}

	ctx.JSON(error.Code, resp)
	ctx.Abort()
}
