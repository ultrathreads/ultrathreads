package handler

import (
	"ultrathreads/bus/core"
	"ultrathreads/bus/event"
	"ultrathreads/cache"
	"ultrathreads/model"
	"ultrathreads/repository"
	"ultrathreads/util/hashid"
	"ultrathreads/util/log"
)

var (
	readStateCache   cache.ReadStateCacheInterface
	postDao          repository.PostRepository
	userReadStateDao repository.UserReadStateRepository
)

// SetReadStateCache 设置阅读状态缓存实例（依赖注入）
func SetReadStateCache(rsc cache.ReadStateCacheInterface) {
	readStateCache = rsc
}

// SetBusPostDao 设置帖子 DAO 实例（依赖注入）
func SetBusPostDao(d repository.PostRepository) {
	postDao = d
}

// SetBusUserReadStateDao 设置用户阅读状态 DAO 实例（依赖注入）
func SetBusUserReadStateDao(d repository.UserReadStateRepository) {
	userReadStateDao = d
}

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
		post := postDao.Get(id)
		if post != nil && post.ThreadId > 0 {
			if err := postDao.IncrViewCount(post.ThreadId); err != nil {
				log.Error("incr root post view count failed: %v", err)
			}
		}
	}
}

// PostViewUserReadStateHandler 负责更新阅读状态并清理缓存
func PostViewUserReadStateHandler(payload event.PostViewed) {
	postID := hashid.Slug2Id[model.Post](payload.PostSlug)
	if postID <= 0 {
		log.Warn("invalid post slug in PostViewed event: %s", payload.PostSlug)
		return
	}

	nodeID := hashid.Slug2Id[model.Node](payload.NodeSlug)
	if nodeID <= 0 {
		log.Warn("invalid node slug in PostViewed event: %s", payload.NodeSlug)
		return
	}

	if err := userReadStateDao.Upsert(payload.UserID, nodeID, payload.ViewedTime); err != nil {
		log.Warn("upsert user read state failed: %v", err)
		return
	}
	readStateCache.InvalidateUserStates(payload.UserID)
}
