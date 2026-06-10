package cache

import (
	"github.com/goburrow/cache"

	"ultrathreads/util/log"
)

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