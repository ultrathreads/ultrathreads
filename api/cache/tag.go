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
	cache         cache.LoadingCache
	hotCache 	  cache.LoadingCache
	postTagsCache cache.LoadingCache
}

var TagCache = newTagCache()

func newTagCache() *tagCache {
	c := &tagCache{}

	c.cache = cache.NewLoadingCache(
		func(key cache.Key) (value cache.Value, e error) {
			value = dao.TagDao.Get(key2Int64(key))
			return
		},
		cache.WithMaximumSize(1000),
		cache.WithExpireAfterAccess(30*time.Minute),
	)

	c.hotCache = cache.NewLoadingCache(
		func(key cache.Key) (value cache.Value, e error) {
			value = dao.TagDao.Find(querybuilder.NewQueryBuilder().Eq("status", model.StatusOk).Desc("id").Limit(10))
			return
		},
		cache.WithMaximumSize(10),
		cache.WithRefreshAfterWrite(30*time.Minute),
	)

	c.postTagsCache = cache.NewLoadingCache(
		func(key cache.Key) (value cache.Value, e error) {
			postId := key2Int64(key)
			postTags := dao.PostTagDao.Find(
				querybuilder.NewQueryBuilder().Where("post_id = ?", postId),
			)
			var tags []model.Tag
			for _, pt := range postTags {
				// 通过局部实例 c 访问单标签缓存
				if tag := c.Get(pt.TagId); tag != nil {
					tags = append(tags, *tag)
				}
			}
			value = tags
			return
		},
		cache.WithMaximumSize(5000),
		cache.WithExpireAfterAccess(30*time.Minute),
	)

	return c
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

// GetPostTags 获取帖子的标签列表（带缓存）
func (c *tagCache) GetPostTags(postId int64) []model.Tag {
	val, err := c.postTagsCache.Get(postId)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	if val != nil {
		return val.([]model.Tag)
	}
	return nil
}

// InvalidatePostTags 清除指定帖子的标签缓存
func (c *tagCache) InvalidatePostTags(postId int64) {
	c.postTagsCache.Invalidate(postId)
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
