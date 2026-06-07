package mock

import (
	"fmt"
	"time"
	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/util"

	"github.com/Pallinder/go-randomdata"
)

// PostTableSeeder -
func PostTableSeeder(needCleanTable bool) {
	if needCleanTable {
		dropAndCreateTable(&model.Post{})
	}

	// ✅ 初始化时间游标：从 30 天前开始，保证所有 mock 数据都是"过去"的时间
	timeCursor := time.Now().AddDate(0, 0, -30)

	// 1. 生成主帖并回填 ThreadId
	var rootPosts []*model.Post
	for i := 0; i < 50; i++ {
		// ✅ 每次生成主帖，时间向前推进 1~6 小时
		createTime, lastCommentTime := advanceTime(&timeCursor, 1, 6)

		post := postFactory(0, 0, createTime, lastCommentTime)
		if err := dao.PostDao.Create(post); err != nil {
			fmt.Printf("mock root post error: %v\n", err)
			continue
		}

		// 关键：入库后拿到自增ID，回填 ThreadId = 自身ID
		post.ThreadId = post.ID
		if err := dao.PostDao.Update(post); err != nil {
			fmt.Printf("mock update thread_id error: %v\n", err)
			continue
		}
		rootPosts = append(rootPosts, post)
	}

	// 2. 为主帖生成树形回复
	for _, root := range rootPosts {
		replyCount := RandInt(0, 6)
		for j := 0; j < replyCount; j++ {
			// ✅ 一级回复时间晚于主帖，随机推进 1~30 分钟
			createTime, lastCommentTime := advanceTime(&timeCursor, 0, 1) // 0~1小时(即0~60分钟)

			reply := postFactory(root.ID, root.ID, createTime, lastCommentTime)
			if err := dao.PostDao.Create(reply); err != nil {
				fmt.Printf("mock reply error: %v\n", err)
				continue
			}

			// 约 1/3 概率生成二级嵌套回复
			if RandInt(0, 3) == 0 {
				// ✅ 二级回复时间晚于一级回复，随机推进 1~15 分钟
				subCreateTime, subLastCommentTime := advanceTime(&timeCursor, 0, 1)

				subReply := postFactory(reply.ID, root.ID, subCreateTime, subLastCommentTime)
				if err := dao.PostDao.Create(subReply); err != nil {
					fmt.Printf("mock sub-reply error: %v\n", err)
				}
			}
		}
	}
}

// advanceTime 推进时间游标并返回创建时间和最后评论时间
// minHours/maxHours: 本次推进的最小时数和最大小时数（支持小数级别可通过改参数类型实现）
// 返回值: createTime, lastCommentTime
func advanceTime(cursor *time.Time, minHours, maxHours int) (int64, int64) {
	// 计算随机偏移量（秒级精度）
	minSecs := minHours * 3600
	maxSecs := maxHours * 3600
	offsetSecs := RandInt(minSecs, maxSecs+1)

	// 推进游标
	*cursor = cursor.Add(time.Duration(offsetSecs) * time.Second)

	ts := util.Timestamp(*cursor)
	// 对于新创建的帖子，LastCommentTime 初始等于 CreateTime
	return ts, ts
}

// postFactory 生成 Post 模拟数据
func postFactory(parentId, threadId int64, createTime, lastCommentTime int64) *model.Post {
	return &model.Post{
		Title:           randomdata.Country(randomdata.FullCountry),
		Content:         randomdata.Paragraph(),
		UserId:          int64(RandInt(1, 10)),
		NodeId:          int64(RandInt(1, 4)),
		ParentId:        parentId,
		ThreadId:        threadId,
		CreateTime:      createTime,
		LastCommentTime: lastCommentTime,
	}
}