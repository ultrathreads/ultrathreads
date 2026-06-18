package cache

import (
	"github.com/goburrow/cache"

	"ultrathreads/model"
	"ultrathreads/util/log"
)

// 缓存键常量
const (
	postRecommendCacheKey = "post_recommend"
	hotTagsCacheKey       = "hot_tags"
)

// readStateKey 阅读状态的复合缓存键
type readStateKey struct {
	UserID int64
	NodeID int64
}

// key2Int64 将 cache.Key 转换为 int64
func key2Int64(key cache.Key) int64 {
	if k, ok := key.(int64); ok {
		return k
	}
	return 0
}

// Caches 聚合所有缓存实例
// CacheLoaders 提供数据加载函数，避免 cache 层直接依赖 dao 层
type Caches struct {
	Node       NodeCacheInterface
	Post       PostCacheInterface
	Tag        TagCacheInterface
	User       UserCacheInterface
	ReadState  ReadStateCacheInterface
	Stat       StatCacheInterface
	Setting    SettingCacheInterface
	ArticleTag ArticleTagCacheInterface
}

// CacheLoaders 定义所有缓存需要的数据加载函数
// 由组装层（cmd/app）负责提供，打破 cache → dao 的依赖
type CacheLoaders struct {
	NodeLoader       func(nodeId int64) *model.Node
	AllNodesLoader   func() []model.Node
	PostRecLoader    func() []model.Post
	TagLoader        func(tagId int64) *model.Tag
	HotTagsLoader    func() []model.Tag
	PostTagsLoader   func(postId int64) []model.Tag
	UserLoader       func(userId int64) *model.User
	UserScoreLoader  func(userId int64) int
	ReadStateLoader  func(userID, nodeID int64) int64
	UserStatesLoader func(userID int64) map[int64]int64
	UserCountLoader  func() int
	PostCountLoader  func() int
	SettingLoader    func(key string) *model.Setting
	ArticleTagLoader func(articleId int64) []int64
}

// NewCaches 创建所有缓存实例
func NewCaches(loaders *CacheLoaders) *Caches {
	return &Caches{
		Node:       NewNodeCache(loaders.NodeLoader, loaders.AllNodesLoader),
		Post:       NewPostCache(loaders.PostRecLoader),
		Tag:        NewTagCache(loaders.TagLoader, loaders.HotTagsLoader, loaders.PostTagsLoader),
		User:       NewUserCache(loaders.UserLoader, loaders.UserScoreLoader),
		ReadState:  NewReadStateCache(loaders.ReadStateLoader, loaders.UserStatesLoader),
		Stat:       NewStatCache(loaders.UserCountLoader, loaders.PostCountLoader),
		Setting:    NewSettingCache(loaders.SettingLoader),
		ArticleTag: NewArticleTagCache(loaders.ArticleTagLoader),
	}
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
