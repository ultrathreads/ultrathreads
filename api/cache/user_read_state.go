package cache

import (
	"time"

	"github.com/goburrow/cache"

	"ultrathreads/dao"
)

// readStateKey 复合缓存键，避免 int64 冲突
type readStateKey struct {
	UserID int64
	NodeID int64
}

type readStateCache struct {
	cache 			cache.LoadingCache
	userStatesCache cache.LoadingCache
}

var ReadStateCache = newReadStateCache()

func newReadStateCache() *readStateCache {
	c := &readStateCache{}
	c.cache = cache.NewLoadingCache(
		func(key cache.Key) (cache.Value, error) {
			k := key.(readStateKey)
			// 直接调用 DAO，不经过 Service，杜绝循环依赖
			lastReadAt := dao.UserReadStateDao.GetLastReadAt(k.UserID, k.NodeID)
			// GetLastReadAt 在未找到时返回 0，这里原样返回即可
			// LoadingCache 会缓存这个 0 值，防止恶意/无效请求穿透到 DB
			return lastReadAt, nil
		},
		cache.WithMaximumSize(5000),          // 已读状态条目远多于用户数，适当调大
		cache.WithExpireAfterAccess(10*time.Minute), // 已读状态时效性要求高，缩短过期时间
	)

	c.userStatesCache = cache.NewLoadingCache(
		func(key cache.Key) (cache.Value, error) {
			userID := key.(int64)
			states := dao.UserReadStateDao.GetAllReadStates(userID)
			if states == nil {
				states = make(map[int64]int64)
			}
			return states, nil
		},
		cache.WithMaximumSize(2000),
		cache.WithExpireAfterAccess(10*time.Minute),
	)

	return c
}

// Get 获取已读时间戳，未命中时自动回源
func (c *readStateCache) Get(userID, nodeID int64) int64 {
	if userID <= 0 || nodeID <= 0 {
		return 0
	}
	val, err := c.cache.Get(readStateKey{UserID: userID, NodeID: nodeID})
	if err != nil {
		return 0
	}
	return val.(int64)
}

// Invalidate 写入成功后主动失效缓存
func (c *readStateCache) Invalidate(userID, nodeID int64) {
	c.cache.Invalidate(readStateKey{UserID: userID, NodeID: nodeID})
}

// GetUserStates 获取用户所有已读状态，返回 map[nodeID]lastReadAt
func (c *readStateCache) GetUserStates(userID int64) map[int64]int64 {
	if userID <= 0 {
		return make(map[int64]int64)
	}
	val, err := c.userStatesCache.Get(userID)
	if err != nil {
		return make(map[int64]int64)
	}
	return val.(map[int64]int64)
}

// InvalidateUserStates 清除指定用户的全量已读缓存
func (c *readStateCache) InvalidateUserStates(userID int64) {
	c.userStatesCache.Invalidate(userID)
}
