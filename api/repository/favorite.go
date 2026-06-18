package repository

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// FavoriteRepository 收藏数据访问契约
type FavoriteRepository interface {
	Get(id int64) *model.Favorite
	Take(where ...interface{}) *model.Favorite
	Find(cnd *querybuilder.QueryBuilder) []model.Favorite
	FindOne(cnd *querybuilder.QueryBuilder) *model.Favorite
	List(cnd *querybuilder.QueryBuilder) ([]model.Favorite, *querybuilder.Paging)
	Create(t *model.Favorite) error
	Update(t *model.Favorite) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64)
}

type favoriteRepo struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepo{db: db}
}

func (r *favoriteRepo) Get(id int64) *model.Favorite {
	ret := &model.Favorite{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *favoriteRepo) Take(where ...interface{}) *model.Favorite {
	ret := &model.Favorite{}
	if err := r.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *favoriteRepo) Find(cnd *querybuilder.QueryBuilder) (list []model.Favorite) {
	cnd.Find(r.db, &list)
	return
}

func (r *favoriteRepo) FindOne(cnd *querybuilder.QueryBuilder) *model.Favorite {
	ret := &model.Favorite{}
	if err := cnd.FindOne(r.db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *favoriteRepo) List(cnd *querybuilder.QueryBuilder) (list []model.Favorite, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.Favorite{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (r *favoriteRepo) Create(t *model.Favorite) error {
	return r.db.Create(t).Error
}

func (r *favoriteRepo) Update(t *model.Favorite) error {
	return r.db.Save(t).Error
}

func (r *favoriteRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.Favorite{}).Where("id = ?", id).Updates(columns).Error
}

func (r *favoriteRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.Favorite{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (r *favoriteRepo) Delete(id int64) {
	r.db.Delete(&model.Favorite{}, "id = ?", id)
}
