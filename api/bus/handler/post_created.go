// bus/handler/post_created.go
package handler

import (
	"ultrathreads/bus/core"
	"ultrathreads/bus/event"
	"ultrathreads/cache"
	"ultrathreads/repository"
	"ultrathreads/util/log"
)

var (
	tagCache   cache.TagCacheInterface
	tagDao     repository.TagRepository
	postTagDao repository.PostTagRepository
)

// SetTagCache 设置标签缓存实例（依赖注入）
func SetTagCache(tc cache.TagCacheInterface) {
	tagCache = tc
}

// SetBusDaos 设置 bus handler 需要的 dao 实例（依赖注入）
func SetBusDaos(tag repository.TagRepository, postTag repository.PostTagRepository) {
	tagDao = tag
	postTagDao = postTag
}

func PostCreatedHandler(sub core.SafeSubscriber) {
	core.SubscribeTyped(sub, func(payload event.PostCreated) {
		log.Debug("payload=%v", payload)

		//标签耗时长，放这里实现
		if payload.IsRoot && len(payload.Tags) > 0 {
			tagIds := tagDao.GetOrCreates(payload.Tags)
			//postTagDao.DeletePostTags(payload.PostID)
			postTagDao.AddPostTags(payload.PostID, tagIds)

			//清除Tag缓存
			tagCache.InvalidatePostTags(payload.PostID)
		}
	})
}
