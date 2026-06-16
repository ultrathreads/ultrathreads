package cache

import (
	"time"
	"github.com/goburrow/cache"
	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
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
	repo     dao.NodeRepository
	cache    cache.LoadingCache
	allCache cache.LoadingCache
}

// ⚠️ 【关键】保留全局变量，但不再在包初始化时创建真实实例
// 它只是一个占位符，等待 InitNodeCache 被调用
var NodeCache NodeCacheInterface

// InitNodeCache 由 main.go / app.go 在组装阶段调用
// 将真实的 DI 实例赋值给全局变量，让旧代码无感切换
func InitNodeCache(repo dao.NodeRepository) {
	NodeCache = NewNodeCache(repo)
}

// NewNodeCache 真正的工厂函数，供 Caches 聚合体和测试使用
func NewNodeCache(repo dao.NodeRepository) NodeCacheInterface {
	c := &nodeCache{repo: repo}

	c.cache = cache.NewLoadingCache(
		func(key cache.Key) (cache.Value, error) {
			return c.repo.Get(key2Int64(key)), nil
		},
		cache.WithMaximumSize(1000),
		cache.WithExpireAfterAccess(30*time.Minute),
	)

	c.allCache = cache.NewLoadingCache(
		func(key cache.Key) (cache.Value, error) {
			return c.repo.Find(querybuilder.NewQueryBuilder().
				Eq("status", model.StatusOk).
				Asc("sort_no").Desc("id")), nil
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
