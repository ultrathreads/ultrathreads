package cache

import (
	"time"

	"github.com/goburrow/cache"

	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/util/log"
	"ultrathreads/util/querybuilder"
)

var (
	hotTagsCacheKey = "hot_tags_cache"
)

type tagCache struct {
	cache    cache.LoadingCache // 标签缓存
	hotCache cache.LoadingCache
}

var TagCache = newTagCache()

func newTagCache() *tagCache {
	return &tagCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = dao.TagDao.Get(key2Int64(key))
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
		hotCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = dao.TagDao.Find(querybuilder.NewQueryBuilder().Eq("status", model.StatusOk).Desc("id").Limit(10))
				return
			},
			cache.WithMaximumSize(10),
			cache.WithRefreshAfterWrite(30*time.Minute),
		),
	}
}

func (c *tagCache) Get(tagId int64) *model.Tag {
	val, err := c.cache.Get(tagId)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	if val != nil {
		return val.(*model.Tag)
	}
	return nil
}

func (c *tagCache) GetList(tagIds []int64) (tags []model.Tag) {
	if len(tagIds) == 0 {
		return nil
	}
	for _, tagId := range tagIds {
		tag := c.Get(tagId)
		if tag != nil {
			tags = append(tags, *tag)
		}
	}
	return
}

func (c *tagCache) Invalidate(tagId int64) {
	c.cache.Invalidate(tagId)
}

func (c *tagCache) GetHot() []model.Tag {
	val, err := c.hotCache.Get(hotTagsCacheKey)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]model.Tag)
	}
	return nil
}

func (c *tagCache) InvalidateHot() {
	c.hotCache.Invalidate(hotTagsCacheKey)
}
