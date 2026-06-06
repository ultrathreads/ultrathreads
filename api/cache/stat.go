package cache

import (
	"time"
	"ultrathreads/dao"
	"ultrathreads/util/querybuilder"
	"github.com/goburrow/cache"
)

type statCache struct {
	userCountCache    cache.LoadingCache
	topicCountCache   cache.LoadingCache
}

var StatCache = newStatCache()

func newStatCache() *statCache {
	return &statCache{
		userCountCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = dao.UserDao.Count(querybuilder.NewQueryBuilder())
				return
			},
			cache.WithMaximumSize(10),
			cache.WithExpireAfterAccess(1*time.Hour),
		),
		topicCountCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = dao.TopicDao.Count(querybuilder.NewQueryBuilder())
				return
			},
			cache.WithMaximumSize(10),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (c *statCache) GetUserCount() int {
	val, err := c.userCountCache.Get("data")
	if err != nil {
		return 0
	}
	return val.(int)
}

func (c *statCache) GetTopicCount() int {
	val, err := c.topicCountCache.Get("data")
	if err != nil {
		return 0
	}
	return val.(int)
}
