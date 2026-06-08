package mock

import (
	"fmt"
	"math/rand"
	"time"

	"ultrathreads/dao"
	"ultrathreads/model"
)

// TagTableSeeder 初始化标签及帖子-标签关联数据
func TagTableSeeder(needCleanTable bool, totalTags int) {
	if needCleanTable {
		dropAndCreateTable(&model.PostTag{})
		dropAndCreateTable(&model.Tag{})
	}

	// 1. 准备基础标签数据
	baseTags := []*model.Tag{
		{Name: "Golang", Status: 0},
		{Name: "React", Status: 0},
		{Name: "Next.js", Status: 0},
		{Name: "TypeScript", Status: 0},
		{Name: "TailwindCSS", Status: 0},
	}

	// 如果要求生成的标签数大于预设，则自动补充通用标签
	for i := len(baseTags); i < totalTags; i++ {
		baseTags = append(baseTags, &model.Tag{
			Name:   fmt.Sprintf("Tag-%d", i+1),
			Status: 0,
		})
	}

	// 2. 批量创建标签并收集 ID
	tagIDs := make([]int64, 0, len(baseTags))
	for _, tag := range baseTags {
		if err := dao.TagDao.Create(tag); err != nil {
			fmt.Printf("[Mock] create tag '%s' error: %v\n", tag.Name, err)
			continue
		}
		tagIDs = append(tagIDs, tag.ID)
	}

	if len(tagIDs) == 0 {
		fmt.Println("[Mock] no tags created, skip post_tag seeding")
		return
	}

	// 3. 获取 10 条根帖 (parentId = 0)
	rootPosts, err := dao.PostDao.GetRootPosts(10)
	if err != nil {
		fmt.Printf("[Mock] get root posts error: %v\n", err)
		return
	}

	if len(rootPosts) == 0 {
		fmt.Println("[Mock] no root posts found, skip post_tag seeding")
		return
	}

	// 4. 为每条根帖随机分配 1~3 个标签，写入 post_tag 表
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	inserted := 0

	for _, post := range rootPosts {
		// 随机决定该帖子关联的标签数量 (1~3)
		count := rng.Intn(3) + 1
		if count > len(tagIDs) {
			count = len(tagIDs)
		}

		// 随机选取不重复的标签
		shuffled := make([]int64, len(tagIDs))
		copy(shuffled, tagIDs)
		rng.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})

		for k := 0; k < count; k++ {
			pt := &model.PostTag{
				PostId: post.ID,
				TagId:  shuffled[k],
			}
			if err := dao.PostTagDao.Create(pt); err != nil {
				fmt.Printf("[Mock] create post_tag (post=%d, tag=%d) error: %v\n",
					post.ID, shuffled[k], err)
				continue
			}
			inserted++
		}
	}

	fmt.Printf("[Mock] tag seeder done: %d tags, %d root posts, %d post_tag relations\n",
		len(tagIDs), len(rootPosts), inserted)
}