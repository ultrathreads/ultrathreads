package cache

import (
	"time"

	"github.com/goburrow/cache"
)

type statCache struct {
	userCountCache cache.LoadingCache
	postCountCache cache.LoadingCache
}

// NewStatCache 创建 StatCache 实例
// userCountLoader 负责加载用户总数
// postCountLoader 负责加载帖子总数
func NewStatCache(
	userCountLoader func() int,
	postCountLoader func() int,
) *statCache {
	c := &statCache{}
	c.userCountCache = cache.NewLoadingCache(
		func(key cache.Key) (value cache.Value, e error) {
			value = userCountLoader()
			return
		},
		cache.WithMaximumSize(10),
		cache.WithExpireAfterAccess(1*time.Hour),
	)
	c.postCountCache = cache.NewLoadingCache(
		func(key cache.Key) (value cache.Value, e error) {
			value = postCountLoader()
			return
		},
		cache.WithMaximumSize(10),
		cache.WithExpireAfterAccess(30*time.Minute),
	)
	return c
}

func (c *statCache) GetUserCount() int {
	val, err := c.userCountCache.Get("data")
	if err != nil {
		return 0
	}
	return val.(int)
}

func (c *statCache) GetPostCount() int {
	val, err := c.postCountCache.Get("data")
	if err != nil {
		return 0
	}
	return val.(int)
}

// StatCacheInterface 定义 StatCache 的接口
type StatCacheInterface interface {
	GetUserCount() int
	GetPostCount() int
}

// 确保 statCache 实现接口
var _ StatCacheInterface = (*statCache)(nil)

// 为了向后兼容，保留类型别名
type StatCache = statCache
