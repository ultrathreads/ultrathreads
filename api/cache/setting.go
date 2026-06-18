package cache

import (
	"time"

	"github.com/goburrow/cache"

	"ultrathreads/model"
)

type settingCache struct {
	cache cache.LoadingCache
}

// NewSettingCache 创建 SettingCache 实例
// loader 函数负责根据 key 加载设置
func NewSettingCache(loader func(key string) *model.Setting) *settingCache {
	c := &settingCache{}
	c.cache = cache.NewLoadingCache(
		func(key cache.Key) (value cache.Value, e error) {
			value = loader(key.(string))
			return
		},
		cache.WithMaximumSize(1000),
		cache.WithExpireAfterAccess(30*time.Minute),
	)
	return c
}

func (c *settingCache) Get(key string) *model.Setting {
	val, err := c.cache.Get(key)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.(*model.Setting)
	}
	return nil
}

func (c *settingCache) GetValue(key string) string {
	sysConfig := c.Get(key)
	if sysConfig == nil {
		return ""
	}
	return sysConfig.Value
}

func (c *settingCache) Invalidate(key string) {
	c.cache.Invalidate(key)
}

// SettingCacheInterface 定义 SettingCache 的接口
type SettingCacheInterface interface {
	Get(key string) *model.Setting
	GetValue(key string) string
	Invalidate(key string)
}

// 确保 settingCache 实现接口
var _ SettingCacheInterface = (*settingCache)(nil)

// 为了向后兼容，保留类型别名
type SettingCache = settingCache
