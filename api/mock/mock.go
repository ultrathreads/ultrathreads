package mock

import (
	"math/rand"

	"ultrathreads/cache"
	"ultrathreads/repository"
	"ultrathreads/util/log"

	"gorm.io/gorm"
)

var (
	userDao    repository.UserRepository
	postDao    repository.PostRepository
	tagDao     repository.TagRepository
	postTagDao repository.PostTagRepository
	settingDao repository.SettingRepository
	linkDao    repository.LinkRepository
	rbacDao    repository.RbacRepository
	mockDB     *gorm.DB
	tagCache   cache.TagCacheInterface
)

// SetMockDaos 设置 mock 包需要的 dao 实例（依赖注入）
func SetMockDaos(repos *repository.Repositories, db *gorm.DB) {
	userDao = repos.User
	postDao = repos.Post
	tagDao = repos.Tag
	postTagDao = repos.PostTag
	settingDao = repos.Setting
	linkDao = repos.Link
	rbacDao = repos.Rbac
	mockDB = db
}

// SetMockTagCache 设置标签缓存实例（依赖注入）
func SetMockTagCache(tc cache.TagCacheInterface) {
	tagCache = tc
}

// dropAndCreateTable 清空并重建表
func dropAndCreateTable(table interface{}) {
	migrator := mockDB.Migrator()

	if migrator.HasTable(table) {
		if err := migrator.DropTable(table); err != nil {
			log.Error("mock: drop table failed: %v", err)
			return
		}
	}

	if err := migrator.CreateTable(table); err != nil {
		log.Error("mock: create table failed: %v", err)
	}
}

// Mock 执行所有数据填充
func Mock(db *gorm.DB) {
	UserTableSeeder(true, 10)
	RbacTableSeeder(true)
	NodeTableSeeder(db, true, 4)
	PostTableSeeder(true, 200)
	TagTableSeeder(true, 6)
	LinkTableSeeder(true, 6)
	SettingTableSeeder(true)
	UpdateNodeTopicCount()
}

// RandInt 生成 [min, max) 范围内的随机整数
func RandInt(min, max int) int {
	if min >= max || min < 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}
