package mock

import (
	"math/rand"
	"ultrathreads/dao"
)

// dropAndCreateTable - 清空表
func dropAndCreateTable(table interface{}) {
	dao.DB().DropTable(table)
	dao.DB().CreateTable(table)
}

// Mock -
func Mock() {
	UserTableSeeder(true,10)
	NodeTableSeeder(true,4)
	PostTableSeeder(true, 200)
	LinkTableSeeder(true,6)
	SettingTableSeeder(true)
}

func RandInt(min, max int) int {
	if min >= max || min < 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}