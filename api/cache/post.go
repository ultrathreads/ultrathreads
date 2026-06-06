package cache

import (
	"time"

	"github.com/goburrow/cache"

	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var (
	postRecommendCacheKey = "recommend_posts_cache"
)

var PostCache = newPostCache()

type postCache struct {
	recommendCache cache.LoadingCache
}

func newPostCache() *postCache {
	return &postCache{
		recommendCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = dao.PostDao.Find(querybuilder.NewQueryBuilder().Eq("recommend", true).Eq("status", model.StatusOk).Limit(20).Desc("last_comment_time"))
				return
			},
			cache.WithMaximumSize(10),
			cache.WithRefreshAfterWrite(30*time.Minute),
		),
	}
}

func (c *postCache) GetRecommendPosts() []model.Post {
	val, err := c.recommendCache.Get(postRecommendCacheKey)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]model.Post)
	}
	return nil
}

func (c *postCache) InvalidateRecommend() {
	c.recommendCache.Invalidate(postRecommendCacheKey)
}
