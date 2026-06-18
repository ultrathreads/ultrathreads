package dao

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// SettingRepository 设置数据访问契约
type SettingRepository interface {
	Get(id int64) *model.Setting
	Take(where ...interface{}) *model.Setting
	Find(cnd *querybuilder.QueryBuilder) []model.Setting
	FindOne(cnd *querybuilder.QueryBuilder) *model.Setting
	List(cnd *querybuilder.QueryBuilder) ([]model.Setting, *querybuilder.Paging)
	Create(t *model.Setting) error
	Update(t *model.Setting) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64)
	GetByKey(key string) *model.Setting
	Transaction(fc func(tx *gorm.DB) error) error
}

type settingRepo struct {
	db *gorm.DB
}

func NewSettingDao(db *gorm.DB) SettingRepository {
	return &settingRepo{db: db}
}

func (r *settingRepo) Get(id int64) *model.Setting {
	ret := &model.Setting{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *settingRepo) Take(where ...interface{}) *model.Setting {
	ret := &model.Setting{}
	if err := r.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *settingRepo) Find(cnd *querybuilder.QueryBuilder) (list []model.Setting) {
	cnd.Find(r.db, &list)
	return
}

func (r *settingRepo) FindOne(cnd *querybuilder.QueryBuilder) *model.Setting {
	ret := &model.Setting{}
	if err := cnd.FindOne(r.db, ret); err != nil {
		return nil
	}
	return ret
}

func (r *settingRepo) List(cnd *querybuilder.QueryBuilder) (list []model.Setting, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.Setting{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (r *settingRepo) Create(t *model.Setting) error {
	return r.db.Create(t).Error
}

func (r *settingRepo) Update(t *model.Setting) error {
	return r.db.Save(t).Error
}

func (r *settingRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.Setting{}).Where("id = ?", id).Updates(columns).Error
}

func (r *settingRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.Setting{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (r *settingRepo) Delete(id int64) {
	r.db.Delete(&model.Setting{}, "id = ?", id)
}

func (r *settingRepo) GetByKey(key string) *model.Setting {
	if len(key) == 0 {
		return nil
	}
	return r.Take("`key` = ?", key)
}

func (r *settingRepo) Transaction(fc func(tx *gorm.DB) error) error {
	return r.db.Transaction(fc)
}
