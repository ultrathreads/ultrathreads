// bus/handler/post_created.go
package handler

import (
	"ultrathreads/bus/core" 
	"ultrathreads/bus/event"
	"ultrathreads/util/log"
	"ultrathreads/dao"
	"ultrathreads/cache"
)

func PostCreatedHandler(sub core.SafeSubscriber) {
	core.SubscribeTyped(sub, func(payload event.PostCreated) {
		log.Debug("payload=%v", payload)

		//标签耗时长，放这里实现
		if payload.IsRoot && len(payload.Tags) > 0 {
			tagIds := dao.TagDao.GetOrCreates(payload.Tags)
			//dao.PostTagDao.DeletePostTags(payload.PostID)
			dao.PostTagDao.AddPostTags(payload.PostID, tagIds)

			//清除Tag缓存
			cache.TagCache.InvalidatePostTags(payload.PostID)
		}
	})
}