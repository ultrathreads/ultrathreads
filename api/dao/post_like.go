package dao

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

func NewPostLikeDao(db *gorm.DB) *postLikeDao {
	return &postLikeDao{db: db}
}

type postLikeDao struct {
	db *gorm.DB
}

// Get 根据 ID 获取点赞记录，未找到返回 nil
func (d *postLikeDao) Get(id int64) *model.PostLike {
	ret := &model.PostLike{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

// Take 按条件获取单条记录，未找到返回 nil
func (d *postLikeDao) Take(where ...interface{}) *model.PostLike {
	ret := &model.PostLike{}
	if err := d.db.Take(ret, where...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *postLikeDao) Find(cnd *querybuilder.QueryBuilder) (list []model.PostLike) {
	cnd.Find(d.db, &list)
	return
}

// FindOne 通过 QueryBuilder 查询单条记录
func (d *postLikeDao) FindOne(cnd *querybuilder.QueryBuilder) *model.PostLike {
	ret := &model.PostLike{}
	if err := cnd.FindOne(d.db, ret); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *postLikeDao) List(cnd *querybuilder.QueryBuilder) (list []model.PostLike, paging *querybuilder.Paging) {
	cnd.Find(d.db, &list)
	count := cnd.Count(d.db, &model.PostLike{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

// Count 统计数量
func (d *postLikeDao) Count(cnd *querybuilder.QueryBuilder) int64 {
	return cnd.Count(d.db, &model.PostLike{})
}

func (d *postLikeDao) Create(t *model.PostLike) error {
	return d.db.Create(t).Error
}

func (d *postLikeDao) Update(t *model.PostLike) error {
	return d.db.Save(t).Error
}

func (d *postLikeDao) Updates(id int64, columns map[string]interface{}) error {
	return d.db.Model(&model.PostLike{}).Where("id = ?", id).Updates(columns).Error
}

func (d *postLikeDao) UpdateColumn(id int64, name string, value interface{}) error {
	return d.db.Model(&model.PostLike{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

// Delete 根据 ID 删除
func (d *postLikeDao) Delete(id int64) error {
	return d.db.Delete(&model.PostLike{}, "id = ?", id).Error
}