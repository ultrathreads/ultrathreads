package repository

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// ArticleRepository 文章数据访问契约
type ArticleRepository interface {
	Get(id int64) *model.Article
	Find(cnd *querybuilder.QueryBuilder) []model.Article
	List(cnd *querybuilder.QueryBuilder) ([]model.Article, *querybuilder.Paging)
	Create(t *model.Article) error
	Update(t *model.Article) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64) error
}

type articleRepo struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepo{db: db}
}

func (r *articleRepo) Get(id int64) *model.Article {
	ret := &model.Article{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (r *articleRepo) Find(cnd *querybuilder.QueryBuilder) (list []model.Article) {
	cnd.Find(r.db, &list)
	return
}

func (r *articleRepo) List(cnd *querybuilder.QueryBuilder) (list []model.Article, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.Article{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (r *articleRepo) Create(t *model.Article) error {
	return r.db.Create(t).Error
}

func (r *articleRepo) Update(t *model.Article) error {
	return r.db.Save(t).Error
}

func (r *articleRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).Updates(columns).Error
}

func (r *articleRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (r *articleRepo) Delete(id int64) error {
	return r.db.Delete(&model.Article{}, "id = ?", id).Error
}
