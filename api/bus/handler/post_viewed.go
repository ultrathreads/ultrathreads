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
         PostViewUserReadStateHandler(payload)
         PostViewCountHandler(payload)
    })
}

// PostViewCountHandler 负责增加帖子浏览量
func PostViewCountHandler(payload event.PostViewed) {
    if payload.PostID > 0 {
        post := dao.PostDao.Get(payload.PostID)
        if post != nil && post.ThreadId > 0 {
            if err := dao.PostDao.IncrViewCount(post.ThreadId); err != nil {
                log.Error("incr root post view count failed: %v", err)
            }
        }
    }
}

// PostViewUserReadStateHandler 负责更新阅读状态并清理缓存
func PostViewUserReadStateHandler(payload event.PostViewed) {
    if err := dao.UserReadStateDao.Upsert(payload.UserID, payload.NodeID, payload.ViewedTime); err != nil {
        log.Error("upsert user read state failed: %v", err)
        return
    }
    cache.ReadStateCache.InvalidateUserStates(payload.UserID)
}