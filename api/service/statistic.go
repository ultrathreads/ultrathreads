package service

import (
	"strconv"
	"sync/atomic"

	"ultrathreads/util/log"
	"ultrathreads/util/querybuilder"
)

// StatisticService 统计业务契约
type StatisticService interface {
	GenerateData()
}

func NewStatisticService(userSvc UserService, postSvc PostService, settingSvc SettingService) StatisticService {
	return &statisticService{
		userSvc:    userSvc,
		postSvc:    postSvc,
		settingSvc: settingSvc,
	}
}

type statisticService struct {
	running    atomic.Bool
	userSvc    UserService
	postSvc    PostService
	settingSvc SettingService
}

func (s *statisticService) GenerateData() {
	if !s.running.CompareAndSwap(false, true) {
		log.Info("statistic is in building")
		return
	}
	defer s.running.Store(false)

	var (
		statUserCount = strconv.FormatInt(s.userSvc.Count(querybuilder.NewQueryBuilder()), 10)
		statPostCount = strconv.FormatInt(s.postSvc.Count(querybuilder.NewQueryBuilder()), 10)
	)

	if err := s.settingSvc.Set("statUserCount", statUserCount, "社区会员", "社区会员总数"); err != nil {
		log.Error("set statUserCount failed: %v", err)
	}
	if err := s.settingSvc.Set("statPostCount", statPostCount, "帖子数", "主题总数"); err != nil {
		log.Error("set statPostCount failed: %v", err)
	}
}