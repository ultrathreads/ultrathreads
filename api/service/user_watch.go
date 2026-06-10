package service

import (
	"errors"

	"gorm.io/gorm" // ✅ v2 替换 jinzhu/gorm

	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

var UserWatchService = newUserWatchService()

func newUserWatchService() *userWatchService {
	return &userWatchService{}
}

type userWatchService struct{}

func (s *userWatchService) Get(id int64) *model.UserWatch {
	return dao.UserWatchDao.Get(id)
}

func (s *userWatchService) Take(where ...interface{}) *model.UserWatch {
	return dao.UserWatchDao.Take(where...)
}

func (s *userWatchService) Find(cnd *querybuilder.QueryBuilder) []model.UserWatch {
	return dao.UserWatchDao.Find(cnd)
}

func (s *userWatchService) FindOne(cnd *querybuilder.QueryBuilder) *model.UserWatch {
	return dao.UserWatchDao.FindOne(cnd)
}

func (s *userWatchService) List(cnd *querybuilder.QueryBuilder) (list []model.UserWatch, paging *querybuilder.Paging) {
	return dao.UserWatchDao.List(cnd)
}

func (s *userWatchService) Create(t *model.UserWatch) error {
	return dao.UserWatchDao.Create(t)
}

func (s *userWatchService) Update(t *model.UserWatch) error {
	return dao.UserWatchDao.Update(t)
}

func (s *userWatchService) Updates(id int64, columns map[string]interface{}) error {
	return dao.UserWatchDao.Updates(id, columns)
}

func (s *userWatchService) UpdateColumn(id int64, name string, value interface{}) error {
	return dao.UserWatchDao.UpdateColumn(id, name, value)
}

// Delete 删除关注记录
func (s *userWatchService) Delete(id int64) error { // ✅ 补充 error 返回值
	return dao.UserWatchDao.Delete(id)
}

// GetBy 根据用户ID和关注者ID获取关注记录
func (s *userWatchService) GetBy(userID int64, watcherID int64) *model.UserWatch {
	return dao.UserWatchDao.Take("user_id = ? AND watcher_id = ?", userID, watcherID)
}

// Count 统计某用户的粉丝数量
func (s *userWatchService) Count(userId int64) int64 {
	var count int64
	// ✅ v2 Count 签名变更：不再需要传指针，直接返回 error（此处忽略）
	dao.DB().Model(&model.UserWatch{}).Where("user_id = ?", userId).Count(&count)
	return count
}

// Recent 获取最近关注列表
func (s *userWatchService) Recent(userId int64, count int) []model.UserWatch {
	return s.Find(querybuilder.NewQueryBuilder().Eq("user_id", userId).Desc("id").Limit(count))
}

// Watch 关注用户
func (s *userWatchService) Watch(userID int64, watcherID int64) error {
	if userID == watcherID {
		return errors.New("不能自己关注自己")
	}
	user := dao.UserDao.Get(userID)
	if user == nil || user.Status != model.StatusOk {
		return errors.New("用户不存在")
	}

	// 判断是否已经关注
	userWatch := dao.UserWatchDao.Take("user_id = ? AND watcher_id = ?", userID, watcherID)
	if userWatch != nil {
		return errors.New("已关注")
	}

	// ✅ v2 标准事务 API + 🔴 修复事务穿透
	return dao.DB().Transaction(func(tx *gorm.DB) error {
		newWatch := &model.UserWatch{
			UserID:     userID,
			WatcherID:  watcherID,
			CreateTime: util.NowTimestamp(),
		}

		// 🔴 关键修复：原代码使用 dao.UserWatchDao.Create（全局db），导致事务完全失效
		if err := tx.Create(newWatch).Error; err != nil {
			return err
		}

		// ⚠️ 注意：发送通知属于副作用操作
		// 如果通知失败不应回滚关注记录，因此放在事务提交后更合理
		// 但如果业务要求通知必须成功才视为关注成功，则保留在事务内
		NotificationService.SendUserWatchNotification(newWatch)

		return nil
	})
}