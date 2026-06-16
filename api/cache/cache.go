package cache

import (
	"github.com/goburrow/cache"

	"ultrathreads/dao"
	"ultrathreads/util/log"
)

type Caches struct {
	Node NodeCacheInterface
}

func NewCaches(repos *dao.Repositories) *Caches {
	// ✅ 创建真实实例
	nodeCache := NewNodeCache(repos.Node)
	
	// ✅ 同时将实例注入全局变量，让旧代码继续能用 cache.NodeCache
	InitNodeCache(repos.Node)
	
	return &Caches{
		Node: nodeCache,
	}
}

func key2Int64(key cache.Key) int64 {
	return key.(int64)
}

func Setup() {
	log.Info("Cache setup")
}

// Shutdown 是优雅退出的统一接口占位。
// goburrow/cache 为纯内存缓存，无需显式关闭资源。
// 若未来替换为 Redis 等带连接的缓存，在此处实现关闭逻辑即可。
func Shutdown() {
	log.Info("Cache shutdown (no-op for in-memory cache)")
}