package handler

import (
    "ultrathreads/bus/core"
    "ultrathreads/bus/event"
    "ultrathreads/cache"
    "ultrathreads/model"
    "ultrathreads/dao"
    "ultrathreads/util/hashid"
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
    id := hashid.Slug2Id[model.Post](payload.PostSlug)
    if id > 0 {
        post := dao.PostDao.Get(id)
        if post != nil && post.ThreadId > 0 {
            if err := dao.PostDao.IncrViewCount(post.ThreadId); err != nil {
                log.Error("incr root post view count failed: %v", err)
            }
        }
    }
}

// PostViewUserReadStateHandler 负责更新阅读状态并清理缓存
func PostViewUserReadStateHandler(payload event.PostViewed) {
    postID := hashid.Slug2Id[model.Post](payload.PostSlug)
    if(postID <= 0) {
        log.Warn("invalid post slug in PostViewed event: %s", payload.PostSlug)
        return
    }

    // ② 关键修复：通过 PostID 查询其所属的 NodeID
    //    绝不能用 postID 直接当作 nodeID 写入
    post := dao.PostDao.Get(postID)
    if post.ID <= 0 {
        log.Warn("get post id for post %d", post.ID)
        return
    }

    if err := dao.UserReadStateDao.Upsert(payload.UserID, post.NodeId, payload.ViewedTime); err != nil {
        log.Warn("upsert user read state failed: %v", err)
        return
    }
    cache.ReadStateCache.InvalidateUserStates(payload.UserID)
}