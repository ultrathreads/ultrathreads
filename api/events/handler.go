package events

import (
	"ultrathreads/dao"
	"ultrathreads/util/log"
	"ultrathreads/cache"
)

func NodeViewedHandler(mgr *Manager) {
	SubscribeTyped(mgr, func(payload NodeViewedPayload) {
		log.Debug("payload=%v", payload)
		if err := dao.UserReadStateDao.Upsert(payload.UserID, payload.NodeID, payload.ViewedTime); err != nil {
			log.Error("upsert user read state failed: payload=%v, err=%v", payload, err)
		}
		// 写入成功后立即失效缓存，下次 Get 自动通过 LoadingCache 加载最新值
		cache.ReadStateCache.InvalidateUserStates(payload.UserID)

		// 2. PostID > 0 时，直接通过 ThreadId 给根帖浏览数 +1
		if payload.PostID > 0 {
			post := dao.PostDao.Get(payload.PostID)
			if post != nil && post.ThreadId > 0 {
				if err := dao.PostDao.IncrViewCount(post.ThreadId); err != nil {
					log.Error("incr root post view count failed: threadId=%d, err=%v", post.ThreadId, err)
				}
			}
		}
	})
}