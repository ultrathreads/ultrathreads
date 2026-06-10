package events

import (
	"fmt"

	"ultrathreads/dao"
	"ultrathreads/util/log"
)

func NodeViewedHandler(mgr *Manager) {
	// ✅ 无需传事件名字符串，无需手动类型断言
	SubscribeTyped(mgr, func(payload NodeViewedPayload) {
		fmt.Printf("node.viewed: userID=%d, nodeID=%d, viewedTime=%d!\n",
			payload.UserID, payload.NodeID, payload.ViewedTime)

		if err := dao.UserReadStateDao.Upsert(payload.UserID, payload.NodeID, payload.ViewedTime); err != nil {
			log.Error("upsert user read state failed: userID=%d, nodeID=%d, err=%v",
				payload.UserID, payload.NodeID, err)
		}
	})
}