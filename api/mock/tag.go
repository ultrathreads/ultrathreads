package mock

import (
	"fmt"

	"ultrathreads/model"
	"ultrathreads/dao"
)

// TagTableSeeder -
func TagTableSeeder(needCleanTable bool, totalTags int) {
	if needCleanTable {
		dropAndCreateTable(&model.Tag{})
	}

	ns := []*model.Tag{
		{
			Name:        "Golang",
			Status:      0,
		},
		{
			Name:        "React",
			Status:      0,
		},
		{
			Name:        "Next.js",
			Status:      0,
		},
		{
			Name:        "TypeScript",
			Status:      0,
		},
		{
			Name:        "TailwindCSS",
			Status:      0,
		},
	}

	for _, n := range ns {
		if err := dao.TagDao.Create(n); err != nil {
			fmt.Printf("mock tag error： %v\n", err)
		}
	}
}