package cache

import (
	"time"

	"github.com/goburrow/cache"
)

type readStateCache struct {
	cache           cache.LoadingCache
	userStatesCache cache.LoadingCache
}

// NewReadStateCache 创建 ReadStateCache 实例
// readStateLoader 负责根据 userID 和 nodeID 加载最后阅读时间
// userStatesLoader 负责根据 userID 加载所有阅读状态
func NewReadStateCache(
	readStateLoader func(userID, nodeID int64) int64,
	userStatesLoader func(userID int64) map[int64]int64,
) *readStateCache {
	c := &readStateCache{}
	c.cache = cache.NewLoadingCache(
		func(key cache.Key) (cache.Value, error) {
			k := key.(readStateKey)
			lastReadAt := readStateLoader(k.UserID, k.NodeID)
			return lastReadAt, nil
		},
		cache.WithMaximumSize(5000),
		cache.WithExpireAfterAccess(10*time.Minute),
	)

	c.userStatesCache = cache.NewLoadingCache(
		func(key cache.Key) (cache.Value, error) {
			userID := key.(int64)
			states := userStatesLoader(userID)
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

func (c *readStateCache) Invalidate(userID, nodeID int64) {
	c.cache.Invalidate(readStateKey{UserID: userID, NodeID: nodeID})
}

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

func (c *readStateCache) InvalidateUserStates(userID int64) {
	c.userStatesCache.Invalidate(userID)
}

// ReadStateCacheInterface 定义 ReadStateCache 的接口
type ReadStateCacheInterface interface {
	Get(userID, nodeID int64) int64
	Invalidate(userID, nodeID int64)
	GetUserStates(userID int64) map[int64]int64
	InvalidateUserStates(userID int64)
}

// 确保 readStateCache 实现接口
var _ ReadStateCacheInterface = (*readStateCache)(nil)

// 为了向后兼容，保留类型别名
type ReadStateCache = readStateCache
