package dao

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// UserWatchRepository 用户关注数据访问契约
type UserWatchRepository interface {
	Get(id int64) *model.UserWatch
	Take(where ...interface{}) *model.UserWatch
	Find(cnd *querybuilder.QueryBuilder) []model.UserWatch
	FindOne(cnd *querybuilder.QueryBuilder) *model.UserWatch
	List(cnd *querybuilder.QueryBuilder) ([]model.UserWatch, *querybuilder.Paging)
	Create(t *model.UserWatch) error
	Update(t *model.UserWatch) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64)
}

type userWatchRepo struct {
	db *gorm.DB
}

func NewUserWatchDao(db *gorm.DB) UserWatchRepository {
	return &userWatchRepo{db: db}
}

func (r *userWatchRepo) Get(id int64) *model.UserWatch {
	ret := &model.UserWatch{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userWatchRepo) Take(where ...interface{}) *model.UserWatch {
	ret := &model.UserWatch{}
	if err := r.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userWatchRepo) Find(cnd *querybuilder.QueryBuilder) (list []model.UserWatch) {
	cnd.Find(r.db, &list)
	return
}

func (r *userWatchRepo) FindOne(cnd *querybuilder.QueryBuilder) *model.UserWatch {
	ret := &model.UserWatch{}
	if err := cnd.FindOne(r.db, ret); err != nil {
		return nil
	}
	return ret
}

func (r *userWatchRepo) List(cnd *querybuilder.QueryBuilder) (list []model.UserWatch, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.UserWatch{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (r *userWatchRepo) Create(t *model.UserWatch) error {
	return r.db.Create(t).Error
}

func (r *userWatchRepo) Update(t *model.UserWatch) error {
	return r.db.Save(t).Error
}

func (r *userWatchRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.UserWatch{}).Where("id = ?", id).Updates(columns).Error
}

func (r *userWatchRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.UserWatch{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (r *userWatchRepo) Delete(id int64) {
	r.db.Delete(&model.UserWatch{}, "id = ?", id)
}