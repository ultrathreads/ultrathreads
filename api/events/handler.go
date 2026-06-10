package events

import (
	"ultrathreads/dao"
	"ultrathreads/util/log"
)

func NodeViewedHandler(mgr *Manager) {
	SubscribeTyped(mgr, func(payload NodeViewedPayload) {
		if err := dao.UserReadStateDao.Upsert(payload.UserID, payload.NodeID, payload.ViewedTime); err != nil {
			log.Error("upsert user read state failed: payload=%v, err=%v", payload, err)
		}
	})
}