package mock

import (
	"fmt"
	"ultrathreads/model"
	"strconv"
)

func linkFactory(i int) *model.Link {
	index := strconv.Itoa(i)
	return &model.Link{
		Title: "link title " + index,
		Url:  "https://www.baidu.com/" + index,
	}
}

// LinksTableSeeder -
func LinkTableSeeder(needCleanTable bool, totalLinks int) {
	if needCleanTable {
		dropAndCreateTable(&model.Link{})
	}

	for i := 0; i < totalLinks; i++ {
		link := linkFactory(i)
		if err := linkDao.Create(link); err != nil {
			fmt.Printf("mock link error： %v\n", err)
		}
	}
}
