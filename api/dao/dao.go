package dao

import "gorm.io/gorm"

// Repositories 新版 DAO 聚合体（与老全局变量完全隔离）
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

func NewRepositories(db *gorm.DB) *Repositories {
	// 最终目标是不要调用setup(db)
	setup(db)

	return &Repositories{
		Node:          NodeDao,
		Post:          PostDao,
		User:          UserDao,
		Article:       ArticleDao,
		ArticleTag:    ArticleTagDao,
		Favorite:      FavoriteDao,
		Link:          LinkDao,
		LoginSource:   LoginSourceDao,
		Notification:  NotificationDao,
		PostLike:      PostLikeDao,
		PostTag:       PostTagDao,
		Rbac:          RbacDao,
		Setting:       SettingDao,
		Tag:           TagDao,
		UserReadState: UserReadStateDao,
		UserScore:     UserScoreDao,
		UserScoreLog:  UserScoreLogDao,
		UserWatch:     UserWatchDao,
	}
}

// DB 返回全局 *gorm.DB 实例（过渡期使用，最终目标是消除此函数）
func DB() *gorm.DB {
	return db
}

var (
	db               *gorm.DB
	NodeDao          NodeRepository
	ArticleDao       ArticleRepository
	ArticleTagDao    ArticleTagRepository
	FavoriteDao      FavoriteRepository
	LinkDao          LinkRepository
	LoginSourceDao   LoginSourceRepository
	NotificationDao  NotificationRepository
	PostDao          PostRepository
	PostLikeDao      PostLikeRepository
	PostTagDao       PostTagRepository
	RbacDao          RbacRepository
	SettingDao       SettingRepository
	TagDao           TagRepository
	UserDao          UserRepository
	UserReadStateDao UserReadStateRepository
	UserScoreDao     UserScoreRepository
	UserScoreLogDao  UserScoreLogRepository
	UserWatchDao     UserWatchRepository
)

// setup 创建所有 DAO 实例，接收外部注入的 *gorm.DB
func setup(gormDB *gorm.DB) {
	// ⚠️ 关键修复：必须将传入的实例赋值给包级变量 db
	// 否则后续 DB() 返回 nil，且无法正确关闭连接
	db = gormDB

	NodeDao = NewNodeDao(gormDB)
	ArticleDao = NewArticleDao(gormDB)
	ArticleTagDao = NewArticleTagDao(gormDB)
	FavoriteDao = NewFavoriteDao(gormDB)
	LinkDao = NewLinkDao(gormDB)
	LoginSourceDao = NewLoginSourceDao(gormDB)
	NotificationDao = NewNotificationDao(gormDB)
	PostDao = NewPostDao(gormDB)
	PostLikeDao = NewPostLikeDao(gormDB)
	PostTagDao = NewPostTagDao(gormDB)
	RbacDao = NewRbacDao(gormDB)
	SettingDao = NewSettingDao(gormDB)
	TagDao = NewTagDao(gormDB)
	UserDao = NewUserDao(gormDB)
	UserReadStateDao = NewUserReadStateDao(gormDB)
	UserScoreDao = NewUserScoreDao(gormDB)
	UserScoreLogDao = NewUserScoreLogDao(gormDB)
	UserWatchDao = NewUserWatchDao(gormDB)
}
