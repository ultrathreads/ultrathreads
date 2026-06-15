package service

import (
	"sync"
	"time"

	"ultrathreads/dao"
)

// 后台管理准入权限码（与数据库 permissions.code 对应）
const PermAdminPanelAccess = "admin:panel:access"

var RbacService = newRbacService()

func newRbacService() *rbacService {
	return &rbacService{
		cache: make(map[int64]*permCacheEntry),
		ttl:   5 * time.Minute,
	}
}

type permCacheEntry struct {
	codes     map[string]struct{}
	expiredAt time.Time
}

type rbacService struct {
	mu    sync.RWMutex
	cache map[int64]*permCacheEntry
	ttl   time.Duration
}

// CanAccessAdminPanel 校验用户是否有后台管理准入权限
func (s *rbacService) CanAccessAdminPanel(userID int64) bool {
	return s.HasPermission(userID, PermAdminPanelAccess)
}

// HasPermission 校验用户是否拥有指定权限码
func (s *rbacService) HasPermission(userID int64, code string) bool {
	codes := s.getUserPermCodes(userID)
	_, ok := codes[code]
	return ok
}

// GetUserPermissions 获取用户所有权限码（带缓存）
func (s *rbacService) GetUserPermissions(userID int64) []string {
	codes := s.getUserPermCodes(userID)
	result := make([]string, 0, len(codes))
	for code := range codes {
		result = append(result, code)
	}
	return result
}

// InvalidateUserCache 用户角色/权限变更时主动清除缓存
func (s *rbacService) InvalidateUserCache(userID int64) {
	s.mu.Lock()
	delete(s.cache, userID)
	s.mu.Unlock()
}

// getUserPermCodes 从缓存或 DAO 获取用户权限码集合
func (s *rbacService) getUserPermCodes(userID int64) map[string]struct{} {
	// 读缓存
	s.mu.RLock()
	entry, ok := s.cache[userID]
	if ok && time.Now().Before(entry.expiredAt) {
		s.mu.RUnlock()
		return entry.codes
	}
	s.mu.RUnlock()

	// 缓存未命中，查库
	codeList := dao.RbacDao.GetUserPermissionCodes(userID)
	codeSet := make(map[string]struct{}, len(codeList))
	for _, c := range codeList {
		codeSet[c] = struct{}{}
	}

	// 写缓存
	s.mu.Lock()
	s.cache[userID] = &permCacheEntry{
		codes:     codeSet,
		expiredAt: time.Now().Add(s.ttl),
	}
	s.mu.Unlock()

	return codeSet
}