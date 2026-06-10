package mock

import (
	"math/rand"

	"ultrathreads/dao"
	"ultrathreads/util/log"
)

// dropAndCreateTable 清空并重建表
// ⚠️ v2 不再有 DropTable/CreateTable，使用 Migrator 接口替代
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
func Mock() {
	UserTableSeeder(true, 10)
	NodeTableSeeder(true, 4)
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