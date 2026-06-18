// bus/handler/post_updated.go
package handler

import (
	"ultrathreads/bus/core"
	"ultrathreads/bus/event"
	"ultrathreads/util/log"
)

func PostUpdatedHandler(sub core.SafeSubscriber) {
	core.SubscribeTyped(sub, func(payload event.PostUpdated) {
		log.Debug("payload=%v", payload)

		//标签耗时长，放这里实现
		if payload.IsRoot && len(payload.Tags) > 0 {
			tagIds := tagDao.GetOrCreates(payload.Tags)
			postTagDao.DeletePostTags(payload.PostID)
			postTagDao.AddPostTags(payload.PostID, tagIds)

			//清除Tag缓存
			tagCache.InvalidatePostTags(payload.PostID)
		}
	})
}
