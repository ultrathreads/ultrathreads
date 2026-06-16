package dao

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/querybuilder"
)

func NewPostTagDao(db *gorm.DB) *postTagDao {
	return &postTagDao{db: db}
}

type postTagDao struct {
	db *gorm.DB
}

func (d *postTagDao) Get(id int64) *model.PostTag {
	ret := &model.PostTag{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *postTagDao) Take(where ...interface{}) *model.PostTag {
	ret := &model.PostTag{}
	if err := d.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *postTagDao) Find(cnd *querybuilder.QueryBuilder) (list []model.PostTag) {
	cnd.Find(d.db, &list)
	return
}

func (d *postTagDao) FindOne(cnd *querybuilder.QueryBuilder) *model.PostTag {
	ret := &model.PostTag{}
	if err := cnd.FindOne(d.db, ret); err != nil {
		return nil
	}
	return ret
}

func (d *postTagDao) List(cnd *querybuilder.QueryBuilder) (list []model.PostTag, paging *querybuilder.Paging) {
	cnd.Find(d.db, &list)
	count := cnd.Count(d.db, &model.PostTag{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (d *postTagDao) Create(t *model.PostTag) (err error) {
	err = d.db.Create(t).Error
	return
}

func (d *postTagDao) Update(t *model.PostTag) (err error) {
	err = d.db.Save(t).Error
	return
}

func (d *postTagDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = d.db.Model(&model.PostTag{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *postTagDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = d.db.Model(&model.PostTag{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *postTagDao) Delete(id int64) {
	d.db.Delete(&model.PostTag{}, "id = ?", id)
}

func (d *postTagDao) AddPostTags(postId int64, tagIds []int64) {
	if postId <= 0 || len(tagIds) == 0 {
		return
	}
	for _, tagId := range tagIds {
		_ = d.Create(&model.PostTag{
			PostId:     postId,
			TagId:      tagId,
			CreateTime: util.NowTimestamp(),
		})
	}
}

func (d *postTagDao) DeletePostTags(postId int64) {
	if postId <= 0 {
		return
	}
	d.db.Where("post_id = ?", postId).Delete(&model.PostTag{})
}