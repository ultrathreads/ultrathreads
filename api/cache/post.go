package cache

import (
	"time"

	"github.com/goburrow/cache"

	"ultrathreads/model"
)

type postCache struct {
	recommendCache cache.LoadingCache
}

// NewPostCache 创建 PostCache 实例
// loader 函数负责从数据源加载推荐帖子列表
func NewPostCache(loader func() []model.Post) *postCache {
	c := &postCache{}
	c.recommendCache = cache.NewLoadingCache(
		func(key cache.Key) (value cache.Value, e error) {
			value = loader()
			return
		},
		cache.WithMaximumSize(10),
		cache.WithRefreshAfterWrite(30*time.Minute),
	)
	return c
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

// PostCacheInterface 定义 PostCache 的接口
type PostCacheInterface interface {
	GetRecommendPosts() []model.Post
	InvalidateRecommend()
}

// 确保 postCache 实现接口
var _ PostCacheInterface = (*postCache)(nil)

// 为了向后兼容，保留类型别名
type PostCache = postCache
