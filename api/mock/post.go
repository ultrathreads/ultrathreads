package mock

import (
	"fmt"
	"ultrathreads/model"
	"ultrathreads/dao"
	"ultrathreads/util"
	"github.com/Pallinder/go-randomdata"
)

// PostTableSeeder -
func PostTableSeeder(needCleanTable bool) {
	if needCleanTable {
		dropAndCreateTable(&model.Post{})
	}

	// 1. 生成主帖并回填 ThreadId
	var rootPosts []*model.Post
	for i := 0; i < 10; i++ {
		post := postFactory(0, 0) // ParentId=0, ThreadId=0(占位)
		if err := dao.PostDao.Create(post); err != nil {
			fmt.Printf("mock root post error: %v\n", err)
			continue
		}

		// ✅ 关键：入库后拿到自增ID，回填 ThreadId = 自身ID
		post.ThreadId = post.ID
		if err := dao.PostDao.Update(post); err != nil {
			fmt.Printf("mock update thread_id error: %v\n", err)
			continue
		}
		rootPosts = append(rootPosts, post)
	}

	// 2. 为主帖生成树形回复
	for _, root := range rootPosts {
		// 每个主帖随机生成 2~4 条一级回复
		replyCount := RandInt(2, 5)
		for j := 0; j < replyCount; j++ {
			// 一级回复：ParentId=主帖ID, ThreadId=主帖ID
			reply := postFactory(root.ID, root.ID)
			if err := dao.PostDao.Create(reply); err != nil {
				fmt.Printf("mock reply error: %v\n", err)
				continue
			}

			// 约 1/3 概率生成二级嵌套回复
			if RandInt(0, 3) == 0 {
				// 二级回复：ParentId=一级回复ID, ThreadId=主帖ID（不变）
				subReply := postFactory(reply.ID, root.ID)
				if err := dao.PostDao.Create(subReply); err != nil {
					fmt.Printf("mock sub-reply error: %v\n", err)
				}
			}
		}
	}
}

// postFactory 生成 Post 模拟数据
// parentId: 直接父级 ID（主帖传 0）
// threadId: 所属主题帖 ID（主帖创建时传 0 占位，回复传主帖 ID）
func postFactory(parentId, threadId int64) *model.Post {
	now := util.NowTimestamp()

	return &model.Post{
		Title:           randomdata.Country(randomdata.FullCountry),
		Content:         randomdata.Paragraph(),
		UserId:          int64(RandInt(1, 10)),
		NodeId:          int64(RandInt(1, 4)),
		ParentId:        parentId,
		ThreadId:        threadId,
		CreateTime:      now,
		LastCommentTime: now,
	}
}