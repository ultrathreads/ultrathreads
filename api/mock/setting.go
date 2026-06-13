package mock

import (
	"encoding/json"
	"fmt"

	"ultrathreads/cache"
	"ultrathreads/dao"
	"ultrathreads/model"
)

// SettingTableSeeder -
func SettingTableSeeder(needCleanTable bool) {
	if needCleanTable {
		dropAndCreateTable(&model.Setting{})
	}

	// ✅ GetHot() 返回 []model.Tag，需先提取 Name 字段转为 []string
	hotTagEntities := cache.TagCache.GetHot()
	tagNames := make([]string, 0, len(hotTagEntities))
	for _, t := range hotTagEntities {
		if t.Name != "" { // 过滤空名称，避免存入无效标签
			tagNames = append(tagNames, t.Name)
		}
	}

	recommendTagsJSON, err := json.Marshal(tagNames)
	if err != nil {
		fmt.Printf("marshal recommendTags error: %v\n", err)
		recommendTagsJSON = []byte(`[]`)
	}

	ns := []*model.Setting{
		{Key: "defaultNodeId", Value: "1"},
		{Key: "siteTitle", Value: "UltraThreads"},
		{Key: "siteDescription", Value: "小而美的开发者社区"},
		{Key: "recommendTags", Value: string(recommendTagsJSON)},
	}

	for _, n := range ns {
		if err := dao.SettingDao.Create(n); err != nil {
			fmt.Printf("mock setting error: %v\n", err)
		}
	}
}