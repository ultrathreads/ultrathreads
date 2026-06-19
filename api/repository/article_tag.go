package repository

import (
	"time"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// ArticleTagRepository 文章标签关联数据访问契约
type ArticleTagRepository interface {
	Get(id int64) *model.ArticleTag
	Take(where ...interface{}) *model.ArticleTag
	Find(cnd *querybuilder.QueryBuilder) []model.ArticleTag
	List(cnd *querybuilder.QueryBuilder) ([]model.ArticleTag, *querybuilder.Paging)
	Create(t *model.ArticleTag) error
	Update(t *model.ArticleTag) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64)
	AddArticleTags(articleId int64, tagIds []int64)
	DeleteArticleTags(articleId int64)
	FindByArticleId(articleId int64) []model.ArticleTag
	FindByTagIds(tagIds []int64, limit int) []model.ArticleTag
}

type articleTagRepo struct {
	db *gorm.DB
}

func NewArticleTagRepository(db *gorm.DB) ArticleTagRepository {
	return &articleTagRepo{db: db}
}

func (r *articleTagRepo) Get(id int64) *model.ArticleTag {
	ret := &model.ArticleTag{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *articleTagRepo) Take(where ...interface{}) *model.ArticleTag {
	ret := &model.ArticleTag{}
	if err := r.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *articleTagRepo) Find(cnd *querybuilder.QueryBuilder) (list []model.ArticleTag) {
	cnd.Find(r.db, &list)
	return
}

func (r *articleTagRepo) List(cnd *querybuilder.QueryBuilder) (list []model.ArticleTag, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.ArticleTag{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (r *articleTagRepo) Create(t *model.ArticleTag) error {
	return r.db.Create(t).Error
}

func (r *articleTagRepo) Update(t *model.ArticleTag) error {
	return r.db.Save(t).Error
}

func (r *articleTagRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.ArticleTag{}).Where("id = ?", id).Updates(columns).Error
}

func (r *articleTagRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.ArticleTag{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (r *articleTagRepo) Delete(id int64) {
	r.db.Delete(&model.ArticleTag{}, "id = ?", id)
}

func (r *articleTagRepo) AddArticleTags(articleId int64, tagIds []int64) {
	if articleId <= 0 || len(tagIds) == 0 {
		return
	}

	for _, tagId := range tagIds {
		_ = r.Create(&model.ArticleTag{
			ArticleId: articleId,
			TagId:     tagId,
			CreatedAt: time.Now(),
		})
	}
}

func (r *articleTagRepo) DeleteArticleTags(articleId int64) {
	if articleId <= 0 {
		return
	}
	r.db.Where("article_id = ?", articleId).Delete(model.ArticleTag{})
}

func (r *articleTagRepo) FindByArticleId(articleId int64) []model.ArticleTag {
	return r.Find(querybuilder.NewQueryBuilder().Where("article_id = ?", articleId))
}

func (r *articleTagRepo) FindByTagIds(tagIds []int64, limit int) []model.ArticleTag {
	if len(tagIds) == 0 {
		return nil
	}
	var articleTags []model.ArticleTag
	q := r.db.Where("tag_id IN ?", tagIds)
	if limit > 0 {
		q = q.Limit(limit)
	}
	q.Find(&articleTags)
	return articleTags
}
