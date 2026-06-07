package mock

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"ultrathreads/dao"
	"ultrathreads/model"

	"github.com/Pallinder/go-randomdata"
)

const TIMESTAMP_MILLI = true

func PostTableSeeder(needCleanTable bool) {
	if needCleanTable {
		dropAndCreateTable(&model.Post{})
	}

	now := time.Now()

	type timeSegment struct {
		count   int
		minSecs int
		maxSecs int
	}
	segments := []timeSegment{
		{count: 5, minSecs: 0, maxSecs: 30 * 60},
		{count: 12, minSecs: 30 * 60, maxSecs: 3 * 3600},
		{count: 13, minSecs: 3 * 3600, maxSecs: 24 * 3600},
		{count: 10, minSecs: 24 * 3600, maxSecs: 7 * 24 * 3600},
		{count: 10, minSecs: 7 * 24 * 3600, maxSecs: 30 * 24 * 3600},
	}

	var rootPosts []*model.Post
	for _, seg := range segments {
		cursor := now.Add(-time.Duration(seg.maxSecs) * time.Second)
		segDuration := seg.maxSecs - seg.minSecs
		avgInterval := segDuration / (seg.count + 1)

		for i := 0; i < seg.count; i++ {
			offset := RandInt(avgInterval/2, avgInterval*2)
			cursor = cursor.Add(time.Duration(offset) * time.Second)

			segmentEnd := now.Add(-time.Duration(seg.minSecs) * time.Second)
			if cursor.After(segmentEnd) {
				cursor = segmentEnd.Add(-time.Duration(RandInt(1, 60)) * time.Second)
			}

			ts := timeToUnix(cursor)
			post := postFactory(0, 0, ts, ts)
			if err := dao.PostDao.Create(post); err != nil {
				fmt.Printf("mock root post error: %v\n", err)
				continue
			}
			post.ThreadId = post.ID
			if err := dao.PostDao.Update(post); err != nil {
				fmt.Printf("mock update thread_id error: %v\n", err)
				continue
			}
			rootPosts = append(rootPosts, post)
		}
	}

	for _, root := range rootPosts {
		rootTime := unixToTime(root.CreateTime)
		ageHours := now.Sub(rootTime).Hours()

		replyCount := calcReplyCount(ageHours)
		if replyCount == 0 {
			continue
		}

		threadCursor := rootTime
		lastCommentTs := root.CreateTime

		for j := 0; j < replyCount; j++ {
			var replyTime time.Time

			if ageHours > 48 && rand.Float64() < 0.15 {
				graveOffset := RandInt(0, 6*3600)
				replyTime = now.Add(-time.Duration(graveOffset) * time.Second)
			} else {
				minGap, maxGap := calcReplyGap(ageHours)
				replyTime = advanceCursorTime(&threadCursor, minGap, maxGap)
			}

			if replyTime.Before(rootTime) {
				replyTime = rootTime.Add(time.Duration(RandInt(60, 300)) * time.Second)
			}
			if replyTime.After(now) {
				replyTime = now.Add(-time.Duration(RandInt(1, 60)) * time.Second)
			}

			replyTs := timeToUnix(replyTime)
			reply := postFactory(root.ID, root.ID, replyTs, replyTs)
			if err := dao.PostDao.Create(reply); err != nil {
				fmt.Printf("mock reply error: %v\n", err)
				continue
			}

			if replyTs > lastCommentTs {
				lastCommentTs = replyTs
			}

			if RandInt(0, 3) == 0 {
				subReplyTime := advanceCursorTime(&replyTime, 1*60, 30*60)
				if subReplyTime.After(now) {
					subReplyTime = now.Add(-time.Duration(RandInt(1, 60)) * time.Second)
				}
				subTs := timeToUnix(subReplyTime)
				subReply := postFactory(reply.ID, root.ID, subTs, subTs)
				if err := dao.PostDao.Create(subReply); err != nil {
					fmt.Printf("mock sub-reply error: %v\n", err)
				}
				if subTs > lastCommentTs {
					lastCommentTs = subTs
				}
			}
		}

		if lastCommentTs != root.CreateTime {
			root.LastCommentTime = lastCommentTs
			if err := dao.PostDao.Update(root); err != nil {
				fmt.Printf("mock update last_comment_time error: %v\n", err)
			}
		}
	}
}

// ========== 🆕 真实感标题生成器 ==========

var titlePrefixes = []string{
	"请问", "求助", "分享", "讨论", "吐槽", "推荐", "避坑", "实测",
	"有没有人知道", "刚发现", "关于", "为什么", "如何", "大家觉得",
}

var titleSuffixes = []string{
	"，求大佬指点", "，有人遇到过吗？", "，附详细教程", "，亲测有效",
	"，踩坑记录", "，欢迎讨论", "，在线等挺急的", "，新手必看",
	"，建议收藏", "，别划走", "，真的绝了", "，后悔没早知道",
	"，更新后续", "，已解决", "，持续更新中", "",
}

// generateRealisticTitle 生成符合真实论坛语义结构的标题
// 长度分布：~10% 短标题(8-15字)，~70% 中等标题(15-35字)，~20% 长标题(35-55字)
func generateRealisticTitle() string {
	r := rand.Float64()

	var coreLen int
	switch {
	case r < 0.1:
		coreLen = RandInt(8, 15)
	case r < 0.8:
		coreLen = RandInt(15, 35)
	default:
		coreLen = RandInt(35, 55)
	}

	// 用 randomdata 生成核心内容，截取到目标长度保证自然截断
	core := randomdata.SillyName() + " " + randomdata.Paragraph()
	runes := []rune(core)
	if len(runes) > coreLen {
		runes = runes[:coreLen]
	}
	coreText := strings.TrimSpace(string(runes))

	// 按概率组合前缀和后缀
	var builder strings.Builder
	if rand.Float64() < 0.7 {
		builder.WriteString(titlePrefixes[rand.Intn(len(titlePrefixes))])
	}
	builder.WriteString(coreText)
	if suffix := titleSuffixes[rand.Intn(len(titleSuffixes))]; suffix != "" {
		builder.WriteString(suffix)
	}

	return builder.String()
}

// ========== 回帖行为模型 ==========

func calcReplyCount(ageHours float64) int {
	switch {
	case ageHours < 1:
		return RandInt(3, 8)
	case ageHours < 6:
		return RandInt(2, 6)
	case ageHours < 24:
		return RandInt(1, 4)
	case ageHours < 72:
		return RandInt(0, 3)
	case ageHours < 7*24:
		if rand.Float64() < 0.75 {
			return 0
		}
		return RandInt(1, 2)
	default:
		if rand.Float64() < 0.85 {
			return 0
		}
		return RandInt(1, 2)
	}
}

func calcReplyGap(ageHours float64) (minSecs, maxSecs int) {
	switch {
	case ageHours < 1:
		return 1 * 60, 30 * 60
	case ageHours < 6:
		return 10 * 60, 2 * 3600
	case ageHours < 24:
		return 30 * 60, 4 * 3600
	case ageHours < 7*24:
		return 2 * 3600, 12 * 3600
	default:
		return 6 * 3600, 24 * 3600
	}
}

// ========== 基础工具函数 ==========

func advanceCursorTime(cursor *time.Time, minSecs, maxSecs int) time.Time {
	offset := RandInt(minSecs, maxSecs+1)
	*cursor = cursor.Add(time.Duration(offset) * time.Second)
	return *cursor
}

func postFactory(parentId, threadId int64, createTime, lastCommentTime int64) *model.Post {
	minValidTs := int64(1577808000)
	if TIMESTAMP_MILLI {
		minValidTs *= 1000
	}
	if createTime < minValidTs || lastCommentTime < minValidTs {
		panic(fmt.Sprintf(
			"invalid timestamp! create=%d, last_comment=%d. Check TIMESTAMP_MILLI.",
			createTime, lastCommentTime,
		))
	}
	return &model.Post{
		Title:           generateRealisticTitle(), // ✅ 替换为真实感标题生成器
		Content:         randomdata.Paragraph(),
		UserId:          int64(RandInt(1, 10)),
		NodeId:          int64(RandInt(1, 4)),
		ParentId:        parentId,
		ThreadId:        threadId,
		CreateTime:      createTime,
		LastCommentTime: lastCommentTime,
	}
}

func timeToUnix(t time.Time) int64 {
	if TIMESTAMP_MILLI {
		return t.UnixMilli()
	}
	return t.Unix()
}

func unixToTime(ts int64) time.Time {
	if TIMESTAMP_MILLI {
		return time.UnixMilli(ts)
	}
	return time.Unix(ts, 0)
}