package dao

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

func NewUserScoreDao(db *gorm.DB) *userScoreDao {
	return &userScoreDao{db: db}
}

type userScoreDao struct {
	db *gorm.DB
}

func (d *userScoreDao) Get(id int64) *model.UserScore {
	ret := &model.UserScore{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *userScoreDao) Take(where ...interface{}) *model.UserScore {
	ret := &model.UserScore{}
	if err := d.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *userScoreDao) Find(cnd *querybuilder.QueryBuilder) (list []model.UserScore) {
	cnd.Find(d.db, &list)
	return
}

func (d *userScoreDao) FindOne(cnd *querybuilder.QueryBuilder) *model.UserScore {
	ret := &model.UserScore{}
	if err := cnd.FindOne(d.db, ret); err != nil {
		return nil
	}
	return ret
}

func (d *userScoreDao) List(cnd *querybuilder.QueryBuilder) (list []model.UserScore, paging *querybuilder.Paging) {
	cnd.Find(d.db, &list)
	count := cnd.Count(d.db, &model.UserScore{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (d *userScoreDao) Create(t *model.UserScore) (err error) {
	err = d.db.Create(t).Error
	return
}

func (d *userScoreDao) Update(t *model.UserScore) (err error) {
	err = d.db.Save(t).Error
	return
}

func (d *userScoreDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = d.db.Model(&model.UserScore{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *userScoreDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = d.db.Model(&model.UserScore{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *userScoreDao) Delete(id int64) {
	d.db.Delete(&model.UserScore{}, "id = ?", id)
}