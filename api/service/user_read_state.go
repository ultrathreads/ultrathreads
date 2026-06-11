package service

import (
	"ultrathreads/cache"
	"ultrathreads/dao"
	"ultrathreads/util/log"
)

var UserReadStateService = newUserReadStateService()

func newUserReadStateService() *userReadStateService {
	return &userReadStateService{}
}

type userReadStateService struct {
}

// GetLastReadAt 获取用户指定节点的已读时间戳
// 优先走 LoadingCache，Miss 时由 cache 层自动回源 DAO 并缓存结果（含零值防穿透）
func (s *userReadStateService) GetLastReadAt(userID, nodeID int64) int64 {
	return cache.ReadStateCache.Get(userID, nodeID)
}

// MarkAsRead 标记节点为已读
// 核心写入入口，保证游标只向前推进 + 缓存即时失效
func (s *userReadStateService) MarkAsRead(userID, nodeID int64, readAt int64) error {
	if userID <= 0 || nodeID <= 0 {
		return nil // 非法参数静默忽略，不阻断主流程
	}

	err := dao.UserReadStateDao.Upsert(userID, nodeID, readAt)
	if err != nil {
		log.Error("MarkAsRead upsert failed: userID=%d, nodeID=%d, err=%v", userID, nodeID, err)
		return err
	}

	// 写入成功后立即失效缓存，下次 Get 自动通过 LoadingCache 加载最新值
	cache.ReadStateCache.Invalidate(userID, nodeID)
	return nil
}

// BatchMarkAsRead 批量标记已读（如"全部已读"功能）
func (s *userReadStateService) BatchMarkAsRead(userID int64, nodeIDs []int64, readAt int64) {
	for _, nodeID := range nodeIDs {
		if err := s.MarkAsRead(userID, nodeID, readAt); err != nil {
			log.Error("BatchMarkAsRead failed: userID=%d, nodeID=%d, err=%v", userID, nodeID, err)
			// 单条失败不中断批量操作，与 IncrTopicCount 容错风格一致
		}
	}
}

// IsUnread 判断指定帖子是否未读
// 供列表接口逐条比对使用
func (s *userReadStateService) IsUnread(userID, nodeID int64, postCreatedAt int64) bool {
	if userID <= 0 {
		return false // 未登录用户不展示未读标记
	}
	lastReadAt := s.GetLastReadAt(userID, nodeID)
	return postCreatedAt > lastReadAt
}

// GetUserReadStates 获取用户所有节点的已读状态
// 直接走全量缓存，一次加载避免 N+1 查询
func (s *userReadStateService) GetUserReadStates(userID int64) map[int64]int64 {
	if userID <= 0 {
		return make(map[int64]int64)
	}
	return cache.ReadStateCache.GetUserStates(userID)
}
