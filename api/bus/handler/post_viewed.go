// bus/handler/post_viewed.go
package handler

import (
	"ultrathreads/bus/core" 
	"ultrathreads/bus/event"
	"ultrathreads/cache"
	"ultrathreads/dao"
	"ultrathreads/util/log"
)

func PostViewedHandler(sub core.SafeSubscriber) {
	core.SubscribeTyped(sub, func(payload event.PostViewed) {
		log.Debug("payload=%v", payload)
		if err := dao.UserReadStateDao.Upsert(payload.UserID, payload.NodeID, payload.ViewedTime); err != nil {
			log.Error("upsert user read state failed: payload=%v, err=%v", payload, err)
		}
		cache.ReadStateCache.InvalidateUserStates(payload.UserID)

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