package dao

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

func NewTagDao(db *gorm.DB) *tagDao {
	return &tagDao{db: db}
}

type tagDao struct {
	db *gorm.DB
}

func (d *tagDao) Get(id int64) *model.Tag {
	ret := &model.Tag{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *tagDao) Take(where ...interface{}) *model.Tag {
	ret := &model.Tag{}
	if err := d.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *tagDao) Find(cnd *querybuilder.QueryBuilder) (list []model.Tag) {
	cnd.Find(d.db, &list)
	return
}

func (d *tagDao) FindOne(cnd *querybuilder.QueryBuilder) *model.Tag {
	ret := &model.Tag{}
	if err := cnd.FindOne(d.db, ret); err != nil {
		return nil
	}
	return ret
}

func (d *tagDao) List(cnd *querybuilder.QueryBuilder) (list []model.Tag, paging *querybuilder.Paging) {
	cnd.Find(d.db, &list)
	count := cnd.Count(d.db, &model.Tag{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (d *tagDao) Create(t *model.Tag) (err error) {
	err = d.db.Create(t).Error
	return
}

func (d *tagDao) Update(t *model.Tag) (err error) {
	err = d.db.Save(t).Error
	return
}

func (d *tagDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = d.db.Model(&model.Tag{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *tagDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = d.db.Model(&model.Tag{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *tagDao) Delete(id int64) error {
	return d.db.Delete(&model.Tag{}, "id = ?", id).Error
}

func (d *tagDao) FindByIds(ids []int64) []model.Tag {
	if len(ids) == 0 {
		return nil
	}

	var tags []model.Tag
	// GORM 的 Where("id IN ?", slice) 会自动处理空切片和单值情况
	d.db.Where("id IN ?", ids).Find(&tags)

	return tags
}

// FindTagIdsByPostIds 通过关联表批量查询
// SELECT post_id, tag_id FROM post_tags WHERE post_id IN (?)
func (d *tagDao) FindTagIdsByPostIds(postIds []int64) map[int64][]int64 {
	if len(postIds) == 0 {
		return nil
	}

	// 直接使用完整的 model.PostTag 结构体
	// GORM 会根据结构体名自动映射到 post_tags 表
	var rows []model.PostTag

	d.db.Select("post_id, tag_id").
		Where("post_id IN ? AND status = ?", postIds, model.StatusOk). // ⚠️ 必须过滤有效状态
		Find(&rows)

	result := make(map[int64][]int64, len(rows))
	for _, r := range rows {
		result[r.PostId] = append(result[r.PostId], r.TagId)
	}
	return result
}

func (d *tagDao) GetByName(name string) *model.Tag {
	if len(name) == 0 {
		return nil
	}
	return d.Take("name = ?", name)
}

func (d *tagDao) GetOrCreate(name string) (*model.Tag, error) {
	if len(name) == 0 {
		return nil, errors.New("标签为空")
	}
	tag := d.GetByName(name)
	if tag != nil {
		return tag, nil
	}
	tag = &model.Tag{
		Name:   name,
		Status: model.StatusOk,
	}
	err := d.Create(tag)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (d *tagDao) GetOrCreates(tags []string) (tagIDs []int64) {
	for _, tagName := range tags {
		tagName = strings.TrimSpace(tagName)
		tag, err := d.GetOrCreate(tagName)
		if err == nil {
			tagIDs = append(tagIDs, tag.ID)
		}
	}
	return
}