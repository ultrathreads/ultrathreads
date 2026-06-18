package middleware

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/model"
	"ultrathreads/service"
)

var lastReadWhiteList = map[string]struct{}{
	"/api/threads":             {},
	"/api/nodes/:slug/threads": {},
	"/api/tags/:slug/threads":  {},
}

// ShouldReturnLastRead 判断当前接口是否需要返回已读状态
func ShouldReturnLastRead(c *gin.Context) bool {
	_, exists := lastReadWhiteList[c.FullPath()]
	return exists
}

// CurrentUserReadState 仅当请求包含 nodeId 或 tagId（路径参数或查询参数）时，
// 才注入当前用户该节点的已读时间戳，避免无关接口产生无谓的 DB 查询
func CurrentUserReadState(userReadStateSvc service.UserReadStateService) gin.HandlerFunc {
	const readStatesKey = "CurrentUserReadStates"

	return func(c *gin.Context) {
		if !ShouldReturnLastRead(c) {
			c.Next()
			return
		}

		var userID int64
		if userVal, exists := c.Get("CurrentUser"); exists && userVal != nil {
			if user, ok := userVal.(*model.User); ok {
				userID = user.ID
			}
		}

		states := userReadStateSvc.GetUserReadStates(userID)
		c.Set(readStatesKey, states)

		c.Next()
	}
}
