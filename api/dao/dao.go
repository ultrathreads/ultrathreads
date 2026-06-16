package dao

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"ultrathreads/model"
	"ultrathreads/util/log"
)

// 数据库驱动常量
const (
	DriverMySQL  = "mysql"
	DriverSQLite = "sqlite"
)

// 全局数据库实例及 DAO（保留全局变量以兼容现有代码，但通过 Setup 统一初始化）
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

// Setup 初始化数据库连接及所有 DAO
func Setup() {
	var err error
	db, err = openDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database: %v", err)
	}

	// AutoMigrate
	if err = db.AutoMigrate(model.Models...); err != nil {
		log.Error("Auto migrate tables failed: %v", err)
	} else {
		log.Info("Database migration completed successfully")
	}

	// 初始化所有 DAO
	initDaos(db)
}

// openDatabase 根据配置创建 GORM 实例并配置连接池
func openDatabase() (*gorm.DB, error) {
	driver := viper.GetString("database.driver")
	logLevel := gormlog.Warn
	if viper.GetBool("database.log_sql") {
		logLevel = gormlog.Info
	}

	var dialector gorm.Dialector
	switch driver {
	case DriverSQLite:
		path := viper.GetString("database.sqlite.path")
		if path == "" {
			return nil, fmt.Errorf("sqlite path is empty")
		}
		dialector = sqlite.Open(path)
		log.Info("Connecting to SQLite3, path: %s", path)

	case DriverMySQL:
		dsn, err := buildMySQLDSN()
		if err != nil {
			return nil, err
		}
		dialector = mysql.Open(dsn)
		log.Info("Connecting to MySQL, database: %s", viper.GetString("database.mysql.name"))

	default:
		return nil, fmt.Errorf("unsupported database driver: %q", driver)
	}

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormlog.Default.LogMode(logLevel),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "ut_",
			SingularTable: false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("gorm open failed: %w", err)
	}

	// 连接池配置
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	switch driver {
	case DriverMySQL:
		sqlDB.SetMaxIdleConns(viper.GetInt("database.mysql.pool.min"))
		sqlDB.SetMaxOpenConns(viper.GetInt("database.mysql.pool.max"))
		sqlDB.SetConnMaxLifetime(time.Minute)
	case DriverSQLite:
		// SQLite 单写连接，避免并发写入锁竞争
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)
	}

	return gormDB, nil
}

// buildMySQLDSN 构建 MySQL DSN，校验必要参数
func buildMySQLDSN() (string, error) {
	host := viper.GetString("database.mysql.host")
	user := viper.GetString("database.mysql.user")
	password := viper.GetString("database.mysql.password")
	name := viper.GetString("database.mysql.name")
	charset := viper.GetString("database.mysql.charset")

	if host == "" || user == "" || name == "" {
		return "", fmt.Errorf("mysql config incomplete: host=%q, user=%q, name=%q", host, user, name)
	}
	if charset == "" {
		charset = "utf8mb4"
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
		user, password, host, name, charset), nil
}

// initDaos 集中初始化所有 DAO 实例
func initDaos(gormDB *gorm.DB) {
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

// Close 关闭数据库连接
func Close() error {
	if db == nil {
		return nil
	}
	log.Info("Closing database connections")
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB failed: %w", err)
	}
	return sqlDB.Close()
}

// DB 获取全局数据库实例（仅用于无法注入的场景）
func DB() *gorm.DB {
	return db
}

// Tx 执行事务，支持无返回值场景
func Tx(txFunc func(tx *gorm.DB) error) error {
	return db.Transaction(txFunc)
}

// TxResult 执行事务并返回结果，适用于需要返回值的业务逻辑
func TxResult[T any](txFunc func(tx *gorm.DB) (T, error)) (T, error) {
	var zero T
	var result T
	err := db.Transaction(func(tx *gorm.DB) error {
		var txErr error
		result, txErr = txFunc(tx)
		return txErr
	})
	if err != nil {
		return zero, err
	}
	return result, nil
}