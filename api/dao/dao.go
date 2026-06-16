package dao

import "gorm.io/gorm"

var (
	db             *gorm.DB
	NodeDao        *nodeDao
	ArticleDao     *articleDao
	ArticleTagDao  *articleTagDao
	FavoriteDao    *favoriteDao
	LinkDao        *linkDao
	LoginSourceDao *loginSourceDao
	NotificationDao *notificationDao
	PostDao        *postDao
	PostLikeDao    *postLikeDao
	PostTagDao     *postTagDao
	RbacDao        *rbacDao
	SettingDao     *settingDao
	TagDao         *tagDao
	UserDao        *userDao
	UserReadStateDao *userReadStateDao
	UserScoreDao   *userScoreDao
	UserScoreLogDao *userScoreLogDao
	UserWatchDao   *userWatchDao
)

// New 创建所有 DAO 实例，接收外部注入的 *gorm.DB
func Setup(gormDB *gorm.DB) {
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

// DB 获取全局数据库实例（仅用于无法注入的场景）
func DB() *gorm.DB {
	return db
}
