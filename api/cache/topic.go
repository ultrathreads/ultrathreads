package cache

import (
	"time"

	"github.com/goburrow/cache"

	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var (
	topicRecommendCacheKey = "recommend_topics_cache"
)

var TopicCache = newTopicCache()

type topicCache struct {
	recommendCache cache.LoadingCache
}

func newTopicCache() *topicCache {
	return &topicCache{
		recommendCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = dao.TopicDao.Find(querybuilder.NewQueryBuilder().Eq("recommend", true).Eq("status", model.StatusOk).Limit(20).Desc("last_comment_time"))
				return
			},
			cache.WithMaximumSize(10),
			cache.WithRefreshAfterWrite(30*time.Minute),
		),
	}
}

func (c *topicCache) GetRecommendTopics() []model.Topic {
	val, err := c.recommendCache.Get(topicRecommendCacheKey)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]model.Topic)
	}
	return nil
}

func (c *topicCache) InvalidateRecommend() {
	c.recommendCache.Invalidate(topicRecommendCacheKey)
}
