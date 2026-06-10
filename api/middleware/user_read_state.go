package middleware

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
)

// CurrentUserReadState 注入当前用户在指定节点的已读时间戳
func CurrentUserReadState() gin.HandlerFunc {
	return func(c *gin.Context) {
		nodeID := util.QueryIntDefault(c, "nodeId", 0)
		if nodeID <= 0 {
			c.Next()
			return
		}

		var userID int64
		if userVal, exists := c.Get("CurrentUser"); exists && userVal != nil {
			if user, ok := userVal.(*model.User); ok {
				userID = user.ID
			}
		}

		// 3. 查询并注入（未登录时 GetLastReadAt 应返回 0，由 Service 层保证）
		lastReadAt := service.UserReadStateService.GetLastReadAt(userID, int64(nodeID))
		c.Set(util.ReadStateKey(nodeID), lastReadAt)

		c.Next()
	}
}