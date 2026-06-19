package repository

import (
	"errors"
	"fmt"
	"time"

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
	// 以下方法封装事务/多表操作，避免 service 层依赖 *gorm.DB
	UpsertByKey(key, value, name, description string) error
	UpsertAll(configs map[string]interface{}) error
}

type settingRepo struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) SettingRepository {
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

// upsertSingle 在事务内执行单条 upsert 操作
func (r *settingRepo) upsertSingle(tx *gorm.DB, key, value, name, description string) error {
	if len(key) == 0 {
		return errors.New("sys config key is null")
	}

	var sysConfig model.Setting
	err := tx.Where("`key` = ?", key).First(&sysConfig).Error
	notFound := errors.Is(err, gorm.ErrRecordNotFound)

	if err != nil && !notFound {
		return err
	}

	if notFound {
		sysConfig = model.Setting{
			Key:       key,
			Value:     value,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if len(name) > 0 {
			sysConfig.Name = name
		}
		if len(description) > 0 {
			sysConfig.Description = description
		}
		if err := tx.Create(&sysConfig).Error; err != nil {
			return err
		}
	} else {
		updates := map[string]interface{}{
			"value":      value,
			"updated_at": time.Now(),
		}
		if len(name) > 0 {
			updates["name"] = name
		}
		if len(description) > 0 {
			updates["description"] = description
		}
		if err := tx.Model(&sysConfig).Updates(updates).Error; err != nil {
			return err
		}
	}
	return nil
}

// UpsertByKey 单条设置项 upsert
func (r *settingRepo) UpsertByKey(key, value, name, description string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return r.upsertSingle(tx, key, value, name, description)
	})
}

// UpsertAll 批量设置项 upsert（事务）
func (r *settingRepo) UpsertAll(configs map[string]interface{}) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for k, v := range configs {
			if err := r.upsertSingle(tx, k, fmt.Sprintf("%v", v), "", ""); err != nil {
				return err
			}
		}
		return nil
	})
}
