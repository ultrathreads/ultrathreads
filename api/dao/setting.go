package dao

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

func NewSettingDao(db *gorm.DB) *settingDao {
	return &settingDao{db: db}
}

type settingDao struct {
	db *gorm.DB
}

func (d *settingDao) Get(id int64) *model.Setting {
	ret := &model.Setting{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *settingDao) Take(where ...interface{}) *model.Setting {
	ret := &model.Setting{}
	if err := d.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *settingDao) Find(cnd *querybuilder.QueryBuilder) (list []model.Setting) {
	cnd.Find(d.db, &list)
	return
}

func (d *settingDao) FindOne(cnd *querybuilder.QueryBuilder) *model.Setting {
	ret := &model.Setting{}
	if err := cnd.FindOne(d.db, ret); err != nil {
		return nil
	}
	return ret
}

func (d *settingDao) List(cnd *querybuilder.QueryBuilder) (list []model.Setting, paging *querybuilder.Paging) {
	cnd.Find(d.db, &list)
	count := cnd.Count(d.db, &model.Setting{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (d *settingDao) Create(t *model.Setting) (err error) {
	err = d.db.Create(t).Error
	return
}

func (d *settingDao) Update(t *model.Setting) (err error) {
	err = d.db.Save(t).Error
	return
}

func (d *settingDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = d.db.Model(&model.Setting{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *settingDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = d.db.Model(&model.Setting{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *settingDao) Delete(id int64) {
	d.db.Delete(&model.Setting{}, "id = ?", id)
}

func (d *settingDao) GetByKey(key string) *model.Setting {
	if len(key) == 0 {
		return nil
	}
	return d.Take("`key` = ?", key)
}