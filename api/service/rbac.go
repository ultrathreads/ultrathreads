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
	roles     map[string]struct{}
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

// GetUserRoles 获取用户所有角色标识
func (s *rbacService) GetUserRoles(userID int64) []string {
	roles := s.getUserRoles(userID)
	result := make([]string, 0, len(roles))
	for role := range roles {
		result = append(result, role)
	}
	return result
}

// GetUserPermissions 获取用户所有权限码
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

// loadUserRbac 统一从 DB 加载权限和角色，并原子性写入缓存
func (s *rbacService) loadUserRbac(userID int64) (map[string]struct{}, map[string]struct{}) {
	codeList := dao.RbacDao.GetUserPermissionCodes(userID)
	roleList := dao.RbacDao.GetUserRoleCodes(userID)

	codeSet := make(map[string]struct{}, len(codeList))
	for _, c := range codeList {
		codeSet[c] = struct{}{}
	}

	roleSet := make(map[string]struct{}, len(roleList))
	for _, r := range roleList {
		roleSet[r] = struct{}{}
	}

	s.mu.Lock()
	s.cache[userID] = &permCacheEntry{
		codes:     codeSet,
		roles:     roleSet,
		expiredAt: time.Now().Add(s.ttl),
	}
	s.mu.Unlock()

	return codeSet, roleSet
}

// getUserPermCodes 从缓存或 DAO 获取用户权限码集合
func (s *rbacService) getUserPermCodes(userID int64) map[string]struct{} {
	s.mu.RLock()
	entry, ok := s.cache[userID]
	if ok && time.Now().Before(entry.expiredAt) {
		s.mu.RUnlock()
		return entry.codes
	}
	s.mu.RUnlock()

	codes, _ := s.loadUserRbac(userID)
	return codes
}

// getUserRoles 从缓存或 DAO 获取用户角色集合
func (s *rbacService) getUserRoles(userID int64) map[string]struct{} {
	s.mu.RLock()
	entry, ok := s.cache[userID]
	if ok && time.Now().Before(entry.expiredAt) {
		s.mu.RUnlock()
		return entry.roles
	}
	s.mu.RUnlock()

	_, roles := s.loadUserRbac(userID)
	return roles
}