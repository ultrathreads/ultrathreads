package service

import (
	"sync"
	"time"
	"ultrathreads/repository"
)

// 后台管理准入权限码（与数据库 permissions.code 对应）
const PermAdminPanelAccess = "admin:panel:access"

// RbacServicer RBAC 业务契约
type RbacServicer interface {
	CanAccessAdminPanel(userID int64) bool
	HasPermission(userID int64, code string) bool
	GetUserRoles(userID int64) []string
	GetUserPermissions(userID int64) []string
	InvalidateUserCache(userID int64)
}

func NewRbacService(repo repository.RbacRepository) RbacServicer {
	return &rbacService{
		repo:  repo,
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
	repo  repository.RbacRepository
	cache map[int64]*permCacheEntry
	ttl   time.Duration
}

func (s *rbacService) CanAccessAdminPanel(userID int64) bool {
	return s.HasPermission(userID, PermAdminPanelAccess)
}

func (s *rbacService) HasPermission(userID int64, code string) bool {
	codes := s.getUserPermCodes(userID)
	_, ok := codes[code]
	return ok
}

func (s *rbacService) GetUserRoles(userID int64) []string {
	roles := s.getUserRoles(userID)
	result := make([]string, 0, len(roles))
	for role := range roles {
		result = append(result, role)
	}
	return result
}

func (s *rbacService) GetUserPermissions(userID int64) []string {
	codes := s.getUserPermCodes(userID)
	result := make([]string, 0, len(codes))
	for code := range codes {
		result = append(result, code)
	}
	return result
}

func (s *rbacService) InvalidateUserCache(userID int64) {
	s.mu.Lock()
	delete(s.cache, userID)
	s.mu.Unlock()
}

func (s *rbacService) loadUserRbac(userID int64) (map[string]struct{}, map[string]struct{}) {
	codeList := s.repo.GetUserPermissionCodes(userID)
	roleList := s.repo.GetUserRoleCodes(userID)

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

func (s *rbacService) getUserPermCodes(userID int64) map[string]struct{} {
	s.mu.RLock()
	entry, ok := s.cache[userID]
	s.mu.RUnlock()

	if ok && time.Now().Before(entry.expiredAt) {
		return entry.codes
	}

	codes, _ := s.loadUserRbac(userID)
	return codes
}

func (s *rbacService) getUserRoles(userID int64) map[string]struct{} {
	s.mu.RLock()
	entry, ok := s.cache[userID]
	s.mu.RUnlock()

	if ok && time.Now().Before(entry.expiredAt) {
		return entry.roles
	}

	_, roles := s.loadUserRbac(userID)
	return roles
}
