package dao

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

func NewUserWatchDao(db *gorm.DB) *userWatchDao {
	return &userWatchDao{db: db}
}

type userWatchDao struct {
	db *gorm.DB
}

func (d *userWatchDao) Get(id int64) *model.UserWatch {
	ret := &model.UserWatch{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *userWatchDao) Take(where ...interface{}) *model.UserWatch {
	ret := &model.UserWatch{}
	if err := d.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *userWatchDao) Find(cnd *querybuilder.QueryBuilder) []model.UserWatch {
	var list []model.UserWatch
	cnd.Find(d.db, &list)
	return list
}

func (d *userWatchDao) FindOne(cnd *querybuilder.QueryBuilder) *model.UserWatch {
	ret := &model.UserWatch{}
	if err := cnd.FindOne(d.db, ret); err != nil {
		return nil
	}
	return ret
}

func (d *userWatchDao) List(cnd *querybuilder.QueryBuilder) ([]model.UserWatch, *querybuilder.Paging) {
	var list []model.UserWatch
	cnd.Find(d.db, &list)

	count := cnd.Count(d.db, &model.UserWatch{})

	paging := &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return list, paging
}

func (d *userWatchDao) Create(t *model.UserWatch) error {
	return d.db.Create(t).Error
}

func (d *userWatchDao) Update(t *model.UserWatch) error {
	return d.db.Save(t).Error
}

func (d *userWatchDao) Updates(id int64, columns map[string]interface{}) error {
	return d.db.Model(&model.UserWatch{}).Where("id = ?", id).Updates(columns).Error
}

func (d *userWatchDao) UpdateColumn(id int64, name string, value interface{}) error {
	return d.db.Model(&model.UserWatch{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

// Delete 删除关注记录
func (d *userWatchDao) Delete(id int64) error {
	return d.db.Delete(&model.UserWatch{}, "id = ?", id).Error
}