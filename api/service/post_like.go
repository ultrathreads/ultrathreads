package service

import (
	"errors"

	"github.com/jinzhu/gorm"

	"ultrathreads/dao"
	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

var PostLikeService = newPostLikeService()

func newPostLikeService() *postLikeService {
	return &postLikeService{}
}

type postLikeService struct {
}

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

func (s *postLikeService) Delete(id int64) {
	dao.PostLikeDao.Delete(id)
}

// 统计数量
func (s *postLikeService) Count(postId int64) int64 {
	var count int64 = 0
	dao.DB().Model(&model.PostLike{}).Where("post_id = ?", postId).Count(&count)
	return count
}

// 最近点赞
func (s *postLikeService) Recent(postId int64, count int) []model.PostLike {
	return s.Find(querybuilder.NewQueryBuilder().Eq("post_id", postId).Desc("id").Limit(count))
}

func (s *postLikeService) Like(userId int64, postId int64) error {
	post := dao.PostDao.Get(postId)
	if post == nil || post.Status != model.StatusOk {
		return errors.New("话题不存在")
	}

	// 判断是否已经点赞了
	postLike := dao.PostLikeDao.Take("user_id = ? and post_id = ?", userId, postId)
	if postLike != nil {
		return errors.New("已点赞")
	}

	return dao.Tx(dao.DB(), func(tx *gorm.DB) error {
		// 点赞
		postLike := &model.PostLike{
			UserId:     userId,
			PostId:    postId,
			CreateTime: util.NowTimestamp(),
		}
		err := dao.PostLikeDao.Create(postLike)
		if err != nil {
			return err
		}
		// 发送点赞通知
		NotificationService.SendPostLikeNotification(postLike)

		return dao.DB().Model(&post).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
	})
}
