package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
)

// CurrentUserReadState 注入当前用户在指定节点的已读时间戳
// ⚠️ 注意：本中间件仅处理携带有效 nodeId 的请求
// 若请求未携带 nodeId，不会注入任何值，Controller 需自行处理缺失情况
func CurrentUserReadState() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 解析 nodeId，缺失或非法直接跳过，不注入任何零值
		nodeIDStr := c.Query("nodeId")
		nodeID, err := strconv.Atoi(nodeIDStr)
		if err != nil || nodeID <= 0 {
			c.Next()
			return
		}

		// 2. 安全获取 CurrentUser，兼容未登录场景
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