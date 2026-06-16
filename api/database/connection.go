package database

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

const (
	DriverMySQL  = "mysql"
	DriverSQLite = "sqlite"
)

// Setup 初始化数据库并执行迁移，返回可用的 *gorm.DB 实例
func Setup() (*gorm.DB, error) {
	gormDB, err := openDatabase()
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err = gormDB.AutoMigrate(model.Models...); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	log.Info("Database migration completed successfully")

	return gormDB, nil
}

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
		host := viper.GetString("database.mysql.host")
		user := viper.GetString("database.mysql.user")
		password := viper.GetString("database.mysql.password")
		name := viper.GetString("database.mysql.name")
		charset := viper.GetString("database.mysql.charset")
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
			user, password, host, name, charset)
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
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)
	}

	return gormDB, nil
}

// Close 安全关闭数据库连接
func Close(db *gorm.DB) error {
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