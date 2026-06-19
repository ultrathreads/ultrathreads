package service

import (
	"github.com/spf13/viper"

	"ultrathreads/cache"
	"ultrathreads/repository"
)

// Services 聚合所有服务实例，作为统一的服务访问入口
type Services struct {
	Node          NodeService
	Post          PostService
	User          UserService
	Article       ArticleService
	ArticleTag    ArticleTagService
	Favorite      FavoriteService
	Link          LinkService
	LoginSource   LoginSourceService
	Notification  NotificationService
	PostLike      PostLikeService
	PostTag       PostTagService
	Rbac          RbacService
	Setting       SettingService
	Tag           TagService
	UserReadState UserReadStateService
	UserScore     UserScoreService
	UserScoreLog  UserScoreLogService
	UserWatch     UserWatchService
	Appinfo       AppinfoService
	Statistic     StatisticService
}

// NewServices 集中初始化所有服务
func NewServices(repos *repository.Repositories, caches *cache.Caches) *Services {
	// 创建无依赖的基础服务
	linkSvc := NewLinkService(repos.Link)
	appinfoSvc := NewAppinfoService()
	articleTagSvc := NewArticleTagService(repos.ArticleTag)
	postTagSvc := NewPostTagService(repos.PostTag)
	scoreLogSvc := NewUserScoreLogService(repos.UserScoreLog)
	loginSourceSvc := NewLoginSourceService(repos.LoginSource)
	notificationSvc := NewNotificationService(repos.Notification, repos.Post, caches.User, caches.Setting)
	postLikeSvc := NewPostLikeService(repos.PostLike)
	rbacSvc := NewRbacService(repos.Rbac)
	settingSvc := NewSettingService(repos.Setting, caches.Setting)
	tagSvc := NewTagService(repos.Tag, caches.Tag)
	userReadStateSvc := NewUserReadStateService(repos.UserReadState)
	userWatchSvc := NewUserWatchService(repos.UserWatch)
	favoriteSvc := NewFavoriteService(repos.Favorite, repos.Article, repos.Post)
	userScoreSvc := NewUserScoreService(repos.UserScore, scoreLogSvc, caches.User)

	// 创建依赖其他服务
	rssCfg := RssConfig{
		BaseURL:    viper.GetString("base.baseUrl"),
		StaticPath: viper.GetString("base.static_path"),
	}
	postSvc := NewPostService(repos.Post, repos.Node, postTagSvc, settingSvc, caches.Tag, caches.User, caches.Setting, rssCfg)
	userSvc := NewUserService(repos.User, repos.Post, caches.User, repos.LoginSource)
	articleSvc := NewArticleService(repos.Article, repos.Tag, repos.ArticleTag, articleTagSvc, caches.ArticleTag, caches.Tag, caches.User, caches.Setting, rssCfg)

	// 创建依赖多种服务的聚合服务
	statisticSvc := NewStatisticService(userSvc, postSvc, settingSvc)

	return &Services{
		Node:          NewNodeService(repos.Node, caches.Node),
		Post:          postSvc,
		User:          userSvc,
		Article:       articleSvc,
		ArticleTag:    articleTagSvc,
		Favorite:      favoriteSvc,
		Link:          linkSvc,
		LoginSource:   loginSourceSvc,
		Notification:  notificationSvc,
		PostLike:      postLikeSvc,
		PostTag:       postTagSvc,
		Rbac:          rbacSvc,
		Setting:       settingSvc,
		Tag:           tagSvc,
		UserReadState: userReadStateSvc,
		UserScore:     userScoreSvc,
		UserScoreLog:  scoreLogSvc,
		UserWatch:     userWatchSvc,
		Appinfo:       appinfoSvc,
		Statistic:     statisticSvc,
	}
}
