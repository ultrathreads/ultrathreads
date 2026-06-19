package mock

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
	"ultrathreads/model"

	"github.com/Pallinder/go-randomdata"
)

const TIMESTAMP_MILLI = true

// PostTableSeeder 生成模拟数据
// totalRootPosts: 期望生成的根帖总数（如 30, 50, 100）
func PostTableSeeder(needCleanTable bool, totalRootPosts int) {
	if needCleanTable {
		dropAndCreateTable(&model.Post{})
	}

	if totalRootPosts <= 0 {
		fmt.Println("[mock] totalRootPosts must be > 0, skipping.")
		return
	}

	now := time.Now()

	// ✅ 使用权重代替固定数量，各段占比: 10% / 24% / 26% / 20% / 20%
	type timeSegment struct {
		weight  float64
		minSecs int
		maxSecs int
	}
	segments := []timeSegment{
		{weight: 0.10, minSecs: 0, maxSecs: 30 * 60},
		{weight: 0.24, minSecs: 30 * 60, maxSecs: 3 * 3600},
		{weight: 0.26, minSecs: 3 * 3600, maxSecs: 24 * 3600},
		{weight: 0.20, minSecs: 24 * 3600, maxSecs: 7 * 24 * 3600},
		{weight: 0.20, minSecs: 7 * 24 * 3600, maxSecs: 30 * 24 * 3600},
	}

	// ✅ 按权重分配每个时间段的实际数量
	weights := make([]float64, len(segments))
	for i, seg := range segments {
		weights[i] = seg.weight
	}
	counts := distributeCounts(totalRootPosts, weights)

	var rootPosts []*model.Post
	for idx, seg := range segments {
		segCount := counts[idx]
		if segCount == 0 {
			continue
		}

		cursor := now.Add(-time.Duration(seg.maxSecs) * time.Second)
		segDuration := seg.maxSecs - seg.minSecs
		avgInterval := segDuration / (segCount + 1)

		for i := 0; i < segCount; i++ {
			offset := RandInt(avgInterval/2, avgInterval*2)
			cursor = cursor.Add(time.Duration(offset) * time.Second)

			segmentEnd := now.Add(-time.Duration(seg.minSecs) * time.Second)
			if cursor.After(segmentEnd) {
				cursor = segmentEnd.Add(-time.Duration(RandInt(1, 60)) * time.Second)
			}

			nodeId := int64(RandInt(1, 5))
			post := postFactory(0, 0, nodeId, cursor, cursor, false)
			if err := postDao.Create(post); err != nil {
				fmt.Printf("mock root post error: %v\n", err)
				continue
			}
			post.ThreadId = post.ID
			if err := postDao.Update(post); err != nil {
				fmt.Printf("mock update thread_id error: %v\n", err)
				continue
			}
			rootPosts = append(rootPosts, post)
		}
	}

	for _, root := range rootPosts {
		rootTime := root.CreatedAt
		ageHours := now.Sub(rootTime).Hours()

		replyCount := calcReplyCount(ageHours)
		if replyCount == 0 {
			continue
		}

		threadCursor := rootTime
		var lastReplyTs time.Time = rootTime

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

			reply := postFactory(root.ID, root.ID, root.NodeId, replyTime, replyTime, true)
			if err := postDao.Create(reply); err != nil {
				fmt.Printf("mock reply error: %v\n", err)
				continue
			}

			if replyTime.After(lastReplyTs) {
				lastReplyTs = replyTime
			}

			if RandInt(0, 3) == 0 {
				subReplyTime := advanceCursorTime(&replyTime, 1*60, 30*60)
				if subReplyTime.After(now) {
					subReplyTime = now.Add(-time.Duration(RandInt(1, 60)) * time.Second)
				}
				subReply := postFactory(reply.ID, root.ID, root.NodeId, subReplyTime, subReplyTime, true)
				if err := postDao.Create(subReply); err != nil {
					fmt.Printf("mock sub-reply error: %v\n", err)
				}
				if subReplyTime.After(lastReplyTs) {
					lastReplyTs = subReplyTime
				}
			}
		}

		if !lastReplyTs.Equal(rootTime) {
			root.LastRepliedAt = lastReplyTs
			if err := postDao.Update(root); err != nil {
				fmt.Printf("mock update last_replied_at error: %v\n", err)
			}
		}
	}

	// ✅ 随机设置 2 个根帖为置顶
	pinCount := 2
	if len(rootPosts) < pinCount {
		pinCount = len(rootPosts)
	}

	// Fisher-Yates 洗牌后取前 N 个，保证不重复且均匀随机
	rand.Shuffle(len(rootPosts), func(i, j int) {
		rootPosts[i], rootPosts[j] = rootPosts[j], rootPosts[i]
	})

	for i := 0; i < pinCount; i++ {
		rootPosts[i].IsPinned = true
		if err := postDao.Update(rootPosts[i]); err != nil {
			fmt.Printf("mock set is_pinned error: %v\n", err)
		} else {
			fmt.Printf("[mock] pinned root post id=%d title=%q\n", rootPosts[i].ID, rootPosts[i].Title)
		}
	}
}

// ========== 权重分配工具函数 ==========

// distributeCounts 使用最大余数法(Largest Remainder Method)按权重分配总数
// 保证分配结果之和严格等于 total，且各段比例尽可能接近原始权重
func distributeCounts(total int, weights []float64) []int {
	n := len(weights)
	result := make([]int, n)
	if total <= 0 || n == 0 {
		return result
	}

	var totalW float64
	for _, w := range weights {
		totalW += w
	}

	type rem struct {
		idx int
		val float64
	}
	rems := make([]rem, n)
	assigned := 0
	for i, w := range weights {
		exact := float64(total) * w / totalW
		floor := int(exact)
		result[i] = floor
		assigned += floor
		rems[i] = rem{idx: i, val: exact - float64(floor)}
	}

	sort.Slice(rems, func(a, b int) bool {
		return rems[a].val > rems[b].val
	})
	for i := 0; i < total-assigned && i < n; i++ {
		result[rems[i].idx]++
	}
	return result
}

// ========== 真实感标题生成器 ==========

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

	core := randomdata.SillyName() + " " + randomdata.Paragraph()
	runes := []rune(core)
	if len(runes) > coreLen {
		runes = runes[:coreLen]
	}
	coreText := strings.TrimSpace(string(runes))

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

// ========== 真实感内容生成器 ==========

// 段落连接词，让多段内容有逻辑衔接感
var paragraphConnectors = []string{
	"", "", "", // 70% 概率不加连接词（自然换行）
	"另外，", "补充一下，", "还有一点，", "顺便说下，",
	"话说回来，", "对了，", "再说说", "除此之外，",
}

// 常见论坛口语化结尾
var contentEndings = []string{
	"", "", "", "", // 40% 概率无特殊结尾
	"\n\n以上就是我的个人看法，仅供参考。",
	"\n\n大家有什么想法欢迎在评论区讨论～",
	"\n\n先写到这，后面有更新再补。",
	"\n\n纯手打，码字不易，觉得有用的话点个赞吧！",
	"\n\n有没有遇到过类似情况的朋友？求分享经验。",
	"\n\n暂时想到这么多，有问题可以留言问我。",
	"\n\n希望能帮到有同样困惑的人。",
}

// generateRealisticContent 生成符合真实论坛风格的多段帖子正文
// 长度分布: ~15% 短帖(1-2段), ~60% 中等(3-5段), ~25% 长帖(6-10段)
func generateRealisticContent() string {
	r := rand.Float64()

	var paraCount int
	switch {
	case r < 0.15:
		paraCount = RandInt(1, 3) // 1-2段
	case r < 0.75:
		paraCount = RandInt(3, 6) // 3-5段
	default:
		paraCount = RandInt(6, 11) // 6-10段
	}

	var builder strings.Builder
	for i := 0; i < paraCount; i++ {
		if i > 0 {
			builder.WriteString("\n\n")
			// 按概率添加段落间连接词
			if connector := paragraphConnectors[rand.Intn(len(paragraphConnectors))]; connector != "" {
				builder.WriteString(connector)
			}
		}

		// 每段 40~180 字，模拟真实段落长度波动
		paraLen := RandInt(40, 181)
		raw := randomdata.Paragraph()
		runes := []rune(raw)
		if len(runes) > paraLen {
			runes = runes[:paraLen]
		}
		builder.WriteString(strings.TrimSpace(string(runes)))
	}

	// 按概率追加口语化结尾
	if ending := contentEndings[rand.Intn(len(contentEndings))]; ending != "" {
		builder.WriteString(ending)
	}

	return builder.String()
}

// generateReplyContent 生成回帖内容
// 70% 短回复(1段 15-80字)，30% 中等回复(2-3段)
func generateReplyContent() string {
	if rand.Float64() < 0.7 {
		paraLen := RandInt(15, 80)
		raw := randomdata.Paragraph()
		runes := []rune(raw)
		if len(runes) > paraLen {
			runes = runes[:paraLen]
		}
		return strings.TrimSpace(string(runes))
	}

	// 中等回复：取主帖生成器的前 2-3 段
	parts := strings.SplitN(generateRealisticContent(), "\n\n", 4)
	limit := RandInt(2, 4)
	if len(parts) > limit {
		parts = parts[:limit]
	}
	return strings.Join(parts, "\n\n")
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

func postFactory(parentId, threadId, nodeId int64, createdAt, lastRepliedAt time.Time, isReply bool) *model.Post {
	content := generateRealisticContent()
	if isReply {
		content = generateReplyContent()
	}

	return &model.Post{
		Title:         generateRealisticTitle(),
		Content:       content,
		UserId:        int64(RandInt(1, 10)),
		NodeId:        nodeId,
		ParentId:      parentId,
		ThreadId:      threadId,
		CreatedAt:     createdAt,
		LastRepliedAt: lastRepliedAt,
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
