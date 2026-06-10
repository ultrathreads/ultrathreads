package service

import (
	"strconv"
	"sync/atomic"

	"ultrathreads/util/log"
	"ultrathreads/util/querybuilder"
)

var StatisticService = newStatisticService()

func newStatisticService() *statisticService {
	return &statisticService{}
}

type statisticService struct {
	running atomic.Bool // ✅ 修复：bool → atomic.Bool，解决并发数据竞争
}

// GenerateData 生成统计数据
func (s *statisticService) GenerateData() {
	// ✅ 原子 CAS 操作，替代非线程安全的 bool 读写
	if !s.running.CompareAndSwap(false, true) {
		log.Info("statistic is in building")
		return
	}
	defer s.running.Store(false)

	var (
		// ✅ int → int64 适配：Count 返回 int64，strconv.FormatInt 替代 strconv.Itoa
		statUserCount = strconv.FormatInt(UserService.Count(querybuilder.NewQueryBuilder()), 10)
		statPostCount = strconv.FormatInt(PostService.Count(querybuilder.NewQueryBuilder()), 10)
	)

	// 注意：Set 内部已有事务，此处无需额外包装
	if err := SettingService.Set("statUserCount", statUserCount, "社区会员", "社区会员总数"); err != nil {
		log.Error("set statUserCount failed: %v", err)
	}
	if err := SettingService.Set("statPostCount", statPostCount, "帖子数", "主题总数"); err != nil {
		log.Error("set statPostCount failed: %v", err)
	}
}