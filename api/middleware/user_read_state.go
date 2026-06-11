package middleware

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
)

// CurrentUserReadState 仅当请求包含 nodeId（路径参数或查询参数）时，
// 才注入当前用户该节点的已读时间戳，避免无关接口产生无效 DB 查询
func CurrentUserReadState() gin.HandlerFunc {
	const readStatesKey = "CurrentUserReadStates"

	return func(c *gin.Context) {
		// ✅ 优先取路径参数，为空则回退到查询参数，两者均无时 nodeID 保持 ""
		nodeID := util.ParamStringDefault(c, "nodeId", "")
		if nodeID == "" {
			nodeID = util.QueryStringDefault(c, "nodeId", "")
		}

		// 两者都缺失时直接跳过，不执行任何 DB 查询
		if nodeID == "" {
			c.Next()
			return
		}

		var userID int64
		if userVal, exists := c.Get("CurrentUser"); exists && userVal != nil {
			if user, ok := userVal.(*model.User); ok {
				userID = user.ID
			}
		}

		states := service.UserReadStateService.GetUserReadStates(userID)
		c.Set(readStatesKey, states)

		c.Next()
	}
}