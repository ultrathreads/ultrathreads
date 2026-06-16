package dao

import "gorm.io/gorm"

// Daos 新版 DAO 聚合体（与老全局变量完全隔离）
type Daos struct {
    Node NodeRepository
    Post *postDao
    User *userDao
}

func NewDaos(db *gorm.DB) *Daos {
    return &Daos{
        Node: NewNodeDao(db),
        Post: NewPostDao(db),
        User: NewUserDao(db),
    }
}

var (
	db              *gorm.DB
	NodeDao         NodeRepository
	ArticleDao      *articleDao
	ArticleTagDao   *articleTagDao
	FavoriteDao     *favoriteDao
	LinkDao         *linkDao
	LoginSourceDao  *loginSourceDao
	NotificationDao *notificationDao
	PostDao         *postDao
	PostLikeDao     *postLikeDao
	PostTagDao      *postTagDao
	RbacDao         *rbacDao
	SettingDao      *settingDao
	TagDao          *tagDao
	UserDao         *userDao
	UserReadStateDao *userReadStateDao
	UserScoreDao    *userScoreDao
	UserScoreLogDao *userScoreLogDao
	UserWatchDao    *userWatchDao
)

// Setup 创建所有 DAO 实例，接收外部注入的 *gorm.DB
func Setup(gormDB *gorm.DB) {
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

// DB 获取全局数据库实例（仅用于无法注入的场景）
func DB() *gorm.DB {
	return db
}

// Close 关闭全局数据库连接（配合 cmd/web.go 优雅退出使用）
func Close() error {
	if sqlDB, err := db.DB(); err == nil {
		return sqlDB.Close()
	}
	return nil
}