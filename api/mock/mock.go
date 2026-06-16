package mock

import (
	"math/rand"

	"gorm.io/gorm"
	"ultrathreads/dao"
	"ultrathreads/util/log"
)

// dropAndCreateTable 清空并重建表
func dropAndCreateTable(table interface{}) {
	db := dao.DB()
	migrator := db.Migrator()

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