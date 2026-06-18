package cache

import (
	"time"
	"github.com/goburrow/cache"
	"ultrathreads/model"
)

var allNodesCacheKey = "all_nodes_cache"

// NodeCacheInterface 定义接口（方便后续解耦和测试）
type NodeCacheInterface interface {
	Get(nodeId int64) *model.Node
	Invalidate(nodeId int64)
	GetAll() []model.Node
	InvalidateAll()
}

type nodeCache struct {
	cache    cache.LoadingCache
	allCache cache.LoadingCache
}

// NewNodeCache 真正的工厂函数，供 Caches 聚合体和测试使用
// nodeLoader 负责根据 nodeId 加载单个节点
// allNodesLoader 负责加载所有节点列表
func NewNodeCache(
	nodeLoader func(nodeId int64) *model.Node,
	allNodesLoader func() []model.Node,
) NodeCacheInterface {
	c := &nodeCache{}

	c.cache = cache.NewLoadingCache(
		func(key cache.Key) (cache.Value, error) {
			return nodeLoader(key2Int64(key)), nil
		},
		cache.WithMaximumSize(1000),
		cache.WithExpireAfterAccess(30*time.Minute),
	)

	c.allCache = cache.NewLoadingCache(
		func(key cache.Key) (cache.Value, error) {
			return allNodesLoader(), nil
		},
		cache.WithMaximumSize(10),
		cache.WithRefreshAfterWrite(30*time.Minute),
	)

	return c
}

func (c *nodeCache) Get(nodeId int64) *model.Node {
	if nodeId <= 0 {
		return nil
	}
	val, err := c.cache.Get(nodeId)
	if err != nil {
		return nil
	}
	return val.(*model.Node)
}

func (c *nodeCache) Invalidate(nodeId int64) {
	c.cache.Invalidate(nodeId)
}

func (c *nodeCache) GetAll() []model.Node {
	val, err := c.allCache.Get(allNodesCacheKey)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]model.Node)
	}
	return nil
}

func (c *nodeCache) InvalidateAll() {
	c.allCache.Invalidate(allNodesCacheKey)
}

// 确保 nodeCache 实现接口
var _ NodeCacheInterface = (*nodeCache)(nil)
