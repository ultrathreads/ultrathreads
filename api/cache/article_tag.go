package cache

import (
	"time"

	"github.com/goburrow/cache"

	"ultrathreads/model"
)

type articleTagCache struct {
	cache cache.LoadingCache
}

// NewArticleTagCache 创建 ArticleTagCache 实例
// loader 函数负责从数据源加载数据
func NewArticleTagCache(loader func(articleId int64) []int64) *articleTagCache {
	c := &articleTagCache{}
	c.cache = cache.NewLoadingCache(
		func(key cache.Key) (value cache.Value, e error) {
			articleId := key2Int64(key)
			tagIds := loader(articleId)
			if len(tagIds) > 0 {
				value = tagIds
			}
			return
		},
		cache.WithMaximumSize(1000),
		cache.WithExpireAfterAccess(30*time.Minute),
	)
	return c
}

func (c *articleTagCache) Get(articleId int64) []int64 {
	val, err := c.cache.Get(articleId)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]int64)
	}
	return nil
}

func (c *articleTagCache) Invalidate(articleId int64) {
	c.cache.Invalidate(articleId)
}

// ArticleTagCacheInterface 定义 ArticleTagCache 的接口
type ArticleTagCacheInterface interface {
	Get(articleId int64) []int64
	Invalidate(articleId int64)
}

// 确保 articleTagCache 实现接口
var _ ArticleTagCacheInterface = (*articleTagCache)(nil)

// 为了向后兼容，保留类型别名
type ArticleTagCache = articleTagCache

// 辅助函数：从 model.ArticleTag 列表中提取 tag IDs
func ExtractTagIdsFromArticleTags(articleTags []model.ArticleTag) []int64 {
	var tagIds []int64
	for _, articleTag := range articleTags {
		tagIds = append(tagIds, articleTag.TagId)
	}
	return tagIds
}
