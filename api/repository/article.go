package repository

import (
	"errors"
	"time"

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
	// 以下方法封装事务/多表操作，避免 service 层依赖 *gorm.DB
	FindByIds(ids []int64) []model.Article
	CreateWithTags(article *model.Article, tagIds []int64) error
	UpdateWithTags(id int64, updates map[string]interface{}, tagIds []int64) error
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

// FindByIds 根据ID列表批量查询文章
func (r *articleRepo) FindByIds(ids []int64) []model.Article {
	if len(ids) == 0 {
		return nil
	}
	var articles []model.Article
	r.db.Where("id IN ?", ids).Find(&articles)
	return articles
}

// CreateWithTags 创建文章并关联标签（事务：创建文章 + 添加标签关联）
func (r *articleRepo) CreateWithTags(article *model.Article, tagIds []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(article).Error; err != nil {
			return err
		}
		if len(tagIds) > 0 {
			articleTags := make([]model.ArticleTag, 0, len(tagIds))
			now := time.Now()
			for _, tagId := range tagIds {
				articleTags = append(articleTags, model.ArticleTag{
					ArticleId: article.ID,
					TagId:     tagId,
					Status:    model.StatusOk,
					CreatedAt: now,
				})
			}
			if err := tx.Create(&articleTags).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateWithTags 更新文章并重建标签关联（事务：更新文章 + 删除旧关联 + 添加新关联）
func (r *articleRepo) UpdateWithTags(id int64, updates map[string]interface{}, tagIds []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Article{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}
		// 删除旧关联
		if err := tx.Where("article_id = ?", id).Delete(&model.ArticleTag{}).Error; err != nil {
			return err
		}
		// 添加新关联
		if len(tagIds) > 0 {
			articleTags := make([]model.ArticleTag, 0, len(tagIds))
			now := time.Now()
			for _, tagId := range tagIds {
				articleTags = append(articleTags, model.ArticleTag{
					ArticleId: id,
					TagId:     tagId,
					Status:    model.StatusOk,
					CreatedAt: now,
				})
			}
			if err := tx.Create(&articleTags).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
