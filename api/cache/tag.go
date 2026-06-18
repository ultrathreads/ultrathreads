package cache

import (
	"time"

	"github.com/goburrow/cache"

	"ultrathreads/model"
	"ultrathreads/util/log"
)

type tagCache struct {
	cache         cache.LoadingCache
	hotCache      cache.LoadingCache
	postTagsCache cache.LoadingCache
}

// NewTagCache 创建 TagCache 实例
// tagLoader 负责根据 tagId 加载单个 tag
// hotTagsLoader 负责加载热门 tags 列表
// postTagsLoader 负责根据 postId 加载该帖子的所有 tags
func NewTagCache(
	tagLoader func(tagId int64) *model.Tag,
	hotTagsLoader func() []model.Tag,
	postTagsLoader func(postId int64) []model.Tag,
) *tagCache {
	c := &tagCache{}

	c.cache = cache.NewLoadingCache(
		func(key cache.Key) (value cache.Value, e error) {
			value = tagLoader(key2Int64(key))
			return
		},
		cache.WithMaximumSize(1000),
		cache.WithExpireAfterAccess(30*time.Minute),
	)

	c.hotCache = cache.NewLoadingCache(
		func(key cache.Key) (value cache.Value, e error) {
			value = hotTagsLoader()
			return
		},
		cache.WithMaximumSize(10),
		cache.WithRefreshAfterWrite(30*time.Minute),
	)

	c.postTagsCache = cache.NewLoadingCache(
		func(key cache.Key) (value cache.Value, e error) {
			postId := key2Int64(key)
			value = postTagsLoader(postId)
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

// TagCacheInterface 定义 TagCache 的接口
type TagCacheInterface interface {
	Get(tagId int64) *model.Tag
	GetList(tagIds []int64) []model.Tag
	Invalidate(tagId int64)
	GetPostTags(postId int64) []model.Tag
	InvalidatePostTags(postId int64)
	GetHot() []model.Tag
	InvalidateHot()
}

// 确保 tagCache 实现接口
var _ TagCacheInterface = (*tagCache)(nil)

// 为了向后兼容，保留类型别名
type TagCache = tagCache
