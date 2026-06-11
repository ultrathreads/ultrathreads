package middleware

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/model"
	"ultrathreads/service"
)

// CurrentUserReadStates 注入当前用户所有节点的已读时间戳
// 不再依赖 nodeId 参数，一次加载全量供下游按需取用
func CurrentUserReadState() gin.HandlerFunc {
	const readStatesKey = "CurrentUserReadStates"

	return func(c *gin.Context) {
		var userID int64
		if userVal, exists := c.Get("CurrentUser"); exists && userVal != nil {
			if user, ok := userVal.(*model.User); ok {
				userID = user.ID
			}
		}

		// 未登录时 GetUserReadStates 返回空 map，下游按 key 取值自然得到零值
		states := service.UserReadStateService.GetUserReadStates(userID)
		c.Set(readStatesKey, states)

		c.Next()
	}
}