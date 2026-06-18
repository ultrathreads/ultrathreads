package service

import (
	"ultrathreads/cache"
	"ultrathreads/dao"
)

// Services 聚合所有服务实例，作为统一的服务访问入口
type Services struct {
	Node          NodeServicer
	Post          PostServicer
	User          UserServicer
	Article       ArticleServicer
	ArticleTag    ArticleTagServicer
	Favorite      FavoriteServicer
	Link          LinkServicer
	LoginSource   LoginSourceServicer
	Notification  NotificationServicer
	PostLike      PostLikeServicer
	PostTag       PostTagServicer
	Rbac          RbacServicer
	Setting       SettingServicer
	Tag           TagServicer
	UserReadState UserReadStateServicer
	UserScore     UserScoreServicer
	UserScoreLog  UserScoreLogServicer
	UserWatch     UserWatchServicer
	Appinfo       AppinfoServicer
	Statistic     StatisticServicer
}

// NewServices 集中初始化所有服务
func NewServices(repos *dao.Repositories, caches *cache.Caches) *Services {
	// 创建无依赖的基础服务
	linkSvc := NewLinkService(repos.Link)
	appinfoSvc := NewAppinfoService()
	articleTagSvc := NewArticleTagService(repos.ArticleTag)
	postTagSvc := NewPostTagService(repos.PostTag)
	scoreLogSvc := NewUserScoreLogService(repos.UserScoreLog)
	loginSourceSvc := NewLoginSourceService(repos.LoginSource)
	notificationSvc := NewNotificationService(repos.Notification, repos.Post)
	postLikeSvc := NewPostLikeService(repos.PostLike)
	rbacSvc := NewRbacService(repos.Rbac)
	settingSvc := NewSettingService(repos.Setting)
	tagSvc := NewTagService(repos.Tag)
	userReadStateSvc := NewUserReadStateService(repos.UserReadState)
	userWatchSvc := NewUserWatchService(repos.UserWatch)
	favoriteSvc := NewFavoriteService(repos.Favorite, repos.Article, repos.Post)
	userScoreSvc := NewUserScoreService(repos.UserScore, scoreLogSvc)

	// 创建依赖其他服务的服务
	postSvc := NewPostService(repos.Post, repos.Node, postTagSvc, settingSvc)
	userSvc := NewUserService(repos.User, repos.Post)
	articleSvc := NewArticleService(repos.Article, repos.Tag, repos.ArticleTag, articleTagSvc)

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
