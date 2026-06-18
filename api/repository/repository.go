package repository

import "gorm.io/gorm"

// Repositories 聚合体
type Repositories struct {
	Node          NodeRepository
	Post          PostRepository
	User          UserRepository
	Article       ArticleRepository
	ArticleTag    ArticleTagRepository
	Favorite      FavoriteRepository
	Link          LinkRepository
	LoginSource   LoginSourceRepository
	Notification  NotificationRepository
	PostLike      PostLikeRepository
	PostTag       PostTagRepository
	Rbac          RbacRepository
	Setting       SettingRepository
	Tag           TagRepository
	UserReadState UserReadStateRepository
	UserScore     UserScoreRepository
	UserScoreLog  UserScoreLogRepository
	UserWatch     UserWatchRepository
}

func NewRepositories(database *gorm.DB) *Repositories {
	return &Repositories{
		Node:          NewNodeRepository(database),
		Post:          NewPostRepository(database),
		User:          NewUserRepository(database),
		Article:       NewArticleRepository(database),
		ArticleTag:    NewArticleTagRepository(database),
		Favorite:      NewFavoriteRepository(database),
		Link:          NewLinkRepository(database),
		LoginSource:   NewLoginSourceRepository(database),
		Notification:  NewNotificationRepository(database),
		PostLike:      NewPostLikeRepository(database),
		PostTag:       NewPostTagRepository(database),
		Rbac:          NewRbacRepository(database),
		Setting:       NewSettingRepository(database),
		Tag:           NewTagRepository(database),
		UserReadState: NewUserReadStateRepository(database),
		UserScore:     NewUserScoreRepository(database),
		UserScoreLog:  NewUserScoreLogRepository(database),
		UserWatch:     NewUserWatchRepository(database),
	}
}
