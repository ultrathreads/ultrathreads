package mock

import (
	"fmt"
	"ultrathreads/model"
	"ultrathreads/dao"
	"ultrathreads/util"
	"github.com/Pallinder/go-randomdata"
)

// TopicTableSeeder -
func TopicTableSeeder(needCleanTable bool) {
	if needCleanTable {
		dropAndCreateTable(&model.Topic{})
	}

	// 1. 生成主帖并回填 ThreadId
	var rootTopics []*model.Topic
	for i := 0; i < 10; i++ {
		topic := topicFactory(0, 0) // ParentId=0, ThreadId=0(占位)
		if err := dao.TopicDao.Create(topic); err != nil {
			fmt.Printf("mock root topic error: %v\n", err)
			continue
		}

		// ✅ 关键：入库后拿到自增ID，回填 ThreadId = 自身ID
		topic.ThreadId = topic.ID
		if err := dao.TopicDao.Update(topic); err != nil {
			fmt.Printf("mock update thread_id error: %v\n", err)
			continue
		}
		rootTopics = append(rootTopics, topic)
	}

	// 2. 为主帖生成树形回复
	for _, root := range rootTopics {
		// 每个主帖随机生成 2~4 条一级回复
		replyCount := RandInt(2, 5)
		for j := 0; j < replyCount; j++ {
			// 一级回复：ParentId=主帖ID, ThreadId=主帖ID
			reply := topicFactory(root.ID, root.ID)
			if err := dao.TopicDao.Create(reply); err != nil {
				fmt.Printf("mock reply error: %v\n", err)
				continue
			}

			// 约 1/3 概率生成二级嵌套回复
			if RandInt(0, 3) == 0 {
				// 二级回复：ParentId=一级回复ID, ThreadId=主帖ID（不变）
				subReply := topicFactory(reply.ID, root.ID)
				if err := dao.TopicDao.Create(subReply); err != nil {
					fmt.Printf("mock sub-reply error: %v\n", err)
				}
			}
		}
	}
}

// topicFactory 生成 Topic 模拟数据
// parentId: 直接父级 ID（主帖传 0）
// threadId: 所属主题帖 ID（主帖创建时传 0 占位，回复传主帖 ID）
func topicFactory(parentId, threadId int64) *model.Topic {
	now := util.NowTimestamp()

	return &model.Topic{
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