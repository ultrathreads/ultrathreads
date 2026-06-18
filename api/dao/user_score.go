package dao

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// UserScoreRepository 用户积分数据访问契约
type UserScoreRepository interface {
	Get(id int64) *model.UserScore
	Take(where ...interface{}) *model.UserScore
	Find(cnd *querybuilder.QueryBuilder) []model.UserScore
	FindOne(cnd *querybuilder.QueryBuilder) *model.UserScore
	List(cnd *querybuilder.QueryBuilder) ([]model.UserScore, *querybuilder.Paging)
	Create(t *model.UserScore) error
	Update(t *model.UserScore) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64)
}

type userScoreRepo struct {
	db *gorm.DB
}

func NewUserScoreDao(db *gorm.DB) UserScoreRepository {
	return &userScoreRepo{db: db}
}

func (r *userScoreRepo) Get(id int64) *model.UserScore {
	ret := &model.UserScore{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userScoreRepo) Take(where ...interface{}) *model.UserScore {
	ret := &model.UserScore{}
	if err := r.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userScoreRepo) Find(cnd *querybuilder.QueryBuilder) (list []model.UserScore) {
	cnd.Find(r.db, &list)
	return
}

func (r *userScoreRepo) FindOne(cnd *querybuilder.QueryBuilder) *model.UserScore {
	ret := &model.UserScore{}
	if err := cnd.FindOne(r.db, ret); err != nil {
		return nil
	}
	return ret
}

func (r *userScoreRepo) List(cnd *querybuilder.QueryBuilder) (list []model.UserScore, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.UserScore{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (r *userScoreRepo) Create(t *model.UserScore) error {
	return r.db.Create(t).Error
}

func (r *userScoreRepo) Update(t *model.UserScore) error {
	return r.db.Save(t).Error
}

func (r *userScoreRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.UserScore{}).Where("id = ?", id).Updates(columns).Error
}

func (r *userScoreRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.UserScore{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (r *userScoreRepo) Delete(id int64) {
	r.db.Delete(&model.UserScore{}, "id = ?", id)
}