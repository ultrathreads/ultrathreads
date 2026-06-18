package cache

import (
	"time"

	"github.com/goburrow/cache"

	"ultrathreads/model"
)

type userCache struct {
	cache      cache.LoadingCache
	scoreCache cache.LoadingCache
}

// NewUserCache 创建 UserCache 实例
// userLoader 负责根据 userId 加载用户
// scoreLoader 负责根据 userId 加载用户分数
func NewUserCache(
	userLoader func(userId int64) *model.User,
	scoreLoader func(userId int64) int,
) *userCache {
	c := &userCache{}
	c.cache = cache.NewLoadingCache(
		func(key cache.Key) (value cache.Value, e error) {
			value = userLoader(key2Int64(key))
			return
		},
		cache.WithMaximumSize(1000),
		cache.WithExpireAfterAccess(30*time.Minute),
	)
	c.scoreCache = cache.NewLoadingCache(
		func(key cache.Key) (value cache.Value, err error) {
			value = scoreLoader(key2Int64(key))
			return
		},
		cache.WithMaximumSize(1000),
		cache.WithExpireAfterAccess(30*time.Minute),
	)
	return c
}

func (c *userCache) Get(userId int64) *model.User {
	if userId <= 0 {
		return nil
	}
	val, err := c.cache.Get(userId)
	if err != nil {
		return nil
	}
	return val.(*model.User)
}

func (c *userCache) Invalidate(userId int64) {
	c.cache.Invalidate(userId)
}

func (c *userCache) GetScore(userId int64) int {
	val, err := c.scoreCache.Get(userId)
	if err != nil {
		return 0
	}
	return val.(int)
}

func (c *userCache) InvalidateScore(userId int64) {
	c.scoreCache.Invalidate(userId)
}

// UserCacheInterface 定义 UserCache 的接口
type UserCacheInterface interface {
	Get(userId int64) *model.User
	Invalidate(userId int64)
	GetScore(userId int64) int
	InvalidateScore(userId int64)
}

// 确保 userCache 实现接口
var _ UserCacheInterface = (*userCache)(nil)

// 为了向后兼容，保留类型别名
type UserCache = userCache
