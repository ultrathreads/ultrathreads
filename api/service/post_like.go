package service

import (
	"errors"

	"gorm.io/gorm" // ✅ v2 替换 jinzhu/gorm

	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/log"
	"ultrathreads/util/querybuilder"
)

var PostLikeService = newPostLikeService()

func newPostLikeService() *postLikeService {
	return &postLikeService{}
}

type postLikeService struct{}

func (s *postLikeService) Get(id int64) *model.PostLike {
	return dao.PostLikeDao.Get(id)
}

func (s *postLikeService) Take(where ...interface{}) *model.PostLike {
	return dao.PostLikeDao.Take(where...)
}

func (s *postLikeService) Find(cnd *querybuilder.QueryBuilder) []model.PostLike {
	return dao.PostLikeDao.Find(cnd)
}

func (s *postLikeService) FindOne(cnd *querybuilder.QueryBuilder) *model.PostLike {
	return dao.PostLikeDao.FindOne(cnd)
}

func (s *postLikeService) List(cnd *querybuilder.QueryBuilder) (list []model.PostLike, paging *querybuilder.Paging) {
	return dao.PostLikeDao.List(cnd)
}

func (s *postLikeService) Create(t *model.PostLike) error {
	return dao.PostLikeDao.Create(t)
}

func (s *postLikeService) Update(t *model.PostLike) error {
	return dao.PostLikeDao.Update(t)
}

func (s *postLikeService) Updates(id int64, columns map[string]interface{}) error {
	return dao.PostLikeDao.Updates(id, columns)
}

func (s *postLikeService) UpdateColumn(id int64, name string, value interface{}) error {
	return dao.PostLikeDao.UpdateColumn(id, name, value)
}

// Delete 删除点赞记录
func (s *postLikeService) Delete(id int64) error { // ✅ 补充 error 返回值
	return dao.PostLikeDao.Delete(id)
}

// Count 统计数量
func (s *postLikeService) Count(postId int64) int64 {
	var count int64
	// ✅ 补充错误处理
	if err := dao.DB().Model(&model.PostLike{}).Where("post_id = ?", postId).Count(&count).Error; err != nil {
		log.Error("PostLikeService.Count failed: %v", err)
	}
	return count
}

// Recent 最近点赞
func (s *postLikeService) Recent(postId int64, count int) []model.PostLike {
	return s.Find(querybuilder.NewQueryBuilder().Eq("post_id", postId).Desc("id").Limit(count))
}

// Like 点赞
func (s *postLikeService) Like(userId int64, postId int64) error {
	post := dao.PostDao.Get(postId)
	if post == nil || post.Status != model.StatusOk {
		return errors.New("话题不存在")
	}

	// 判断是否已经点赞了
	postLike := dao.PostLikeDao.Take("user_id = ? AND post_id = ?", userId, postId)
	if postLike != nil {
		return errors.New("已点赞")
	}

	// ✅ v2 事务 + 修复事务穿透：全部使用 tx 操作
	return dao.DB().Transaction(func(tx *gorm.DB) error {
		// 1. 使用 tx 创建点赞记录（原代码 dao.PostLikeDao.Create 用的是全局 db，不受 tx 控制）
		newLike := &model.PostLike{
			UserId:     userId,
			PostId:     postId,
			CreateTime: util.NowTimestamp(),
		}
		if err := tx.Create(newLike).Error; err != nil {
			return err
		}

		// 2. 使用 tx 更新帖子点赞数（原代码 dao.DB() 导致此 UPDATE 不在事务内）
		if err := tx.Model(&model.Post{}).Where("id = ?", postId).
			UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error; err != nil {
			return err
		}

		// 3. 发送通知（非 DB 操作，放在事务内即可；若通知失败不应回滚点赞，可移到事务外）
		NotificationService.SendPostLikeNotification(newLike)

		return nil
	})
}