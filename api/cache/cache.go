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
