package dao

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

func NewArticleTagDao(db *gorm.DB) *articleTagDao {
	return &articleTagDao{db: db}
}

type articleTagDao struct {
	db *gorm.DB
}

func (d *articleTagDao) Get(id int64) *model.ArticleTag {
	ret := &model.ArticleTag{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *articleTagDao) Take(where ...interface{}) *model.ArticleTag {
	ret := &model.ArticleTag{}
	if err := d.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *articleTagDao) Find(cnd *querybuilder.QueryBuilder) (list []model.ArticleTag) {
	cnd.Find(d.db, &list)
	return
}

func (d *articleTagDao) List(cnd *querybuilder.QueryBuilder) (list []model.ArticleTag, paging *querybuilder.Paging) {
	cnd.Find(d.db, &list)
	count := cnd.Count(d.db, &model.ArticleTag{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (d *articleTagDao) Create(t *model.ArticleTag) (err error) {
	err = d.db.Create(t).Error
	return
}

func (d *articleTagDao) Update(t *model.ArticleTag) (err error) {
	err = d.db.Save(t).Error
	return
}

func (d *articleTagDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = d.db.Model(&model.ArticleTag{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *articleTagDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = d.db.Model(&model.ArticleTag{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *articleTagDao) Delete(id int64) {
	d.db.Delete(&model.ArticleTag{}, "id = ?", id)
}

func (d *articleTagDao) AddArticleTags(articleId int64, tagIds []int64) {
	if articleId <= 0 || len(tagIds) == 0 {
		return
	}

	for _, tagId := range tagIds {
		_ = d.Create(&model.ArticleTag{
			ArticleId:  articleId,
			TagId:      tagId,
			CreateTime: util.NowTimestamp(),
		})
	}
}

func (d *articleTagDao) DeleteArticleTags(articleId int64) {
	if articleId <= 0 {
		return
	}
	d.db.Where("article_id = ?", articleId).Delete(model.ArticleTag{})
}

func (d *articleTagDao) FindByArticleId(articleId int64) []model.ArticleTag {
	return d.Find(querybuilder.NewQueryBuilder().Where("article_id = ?", articleId))
}