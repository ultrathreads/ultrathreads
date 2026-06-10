package dao

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	"github.com/spf13/viper"

	"ultrathreads/model"
	"ultrathreads/util/log"
)

var (
	db *gorm.DB
)

const DRIVER_MYSQL = "mysql"
const DRIVER_SQLITE = "sqlite"

// Setup 初始化数据库连接（GORM v2）
func Setup() {
	var err error
	var dialector gorm.Dialector

	// 1. 日志级别映射
	logLevel := gormlog.Warn // 默认静默+警告
	if viper.GetBool("database.log_sql") {
		logLevel = gormlog.Info // 打印所有 SQL
	}

	switch viper.GetString("database.driver") {
	case DRIVER_SQLITE:
		path := viper.GetString("database.sqlite.path")
		dialector = sqlite.Open(path)
		log.Info("Connecting to SQLite3, path: %s", path)

	case DRIVER_MYSQL:
		host := viper.GetString("database.mysql.host")
		user := viper.GetString("database.mysql.user")
		password := viper.GetString("database.mysql.password")
		name := viper.GetString("database.mysql.name")
		charset := viper.GetString("database.mysql.charset")

		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
			user, password, host, name, charset)
		dialector = mysql.Open(dsn)
		log.Info("Connecting to MySQL, database: %s", name)

	default:
		log.Fatal("Unsupported database driver: %s", viper.GetString("database.driver"))
		return
	}

	// 2. GORM v2 初始化
	db, err = gorm.Open(dialector, &gorm.Config{
		Logger: gormlog.Default.LogMode(logLevel),
		// v2 默认使用单数表名，无需 SingularTable(true)
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect database: %v", err))
	}

	// 3. 连接池配置（仅 MySQL 需要）
	if viper.GetString("database.driver") == DRIVER_MYSQL {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(viper.GetInt("database.mysql.pool.min"))
		sqlDB.SetMaxOpenConns(viper.GetInt("database.mysql.pool.max"))
		sqlDB.SetConnMaxLifetime(time.Minute)
	}

	// 4. AutoMigrate
	if err = db.AutoMigrate(model.Models...); err != nil {
		log.Error("Auto migrate tables failed: %v", err)
	} else {
		log.Info("Database migration completed successfully")
	}
}

// Close 关闭数据库连接（替代原 Shutdown）
func Close() error {
	log.Info("Closing database connections")
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// DB 获取全局数据库实例
func DB() *gorm.DB {
	return db
}

// Tx 事务环绕（GORM v2 推荐写法）
func Tx(txFunc func(tx *gorm.DB) error) (err error) {
	return db.Transaction(txFunc)
}