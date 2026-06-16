package dao

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

func NewUserScoreLogDao(db *gorm.DB) *userScoreLogDao {
	return &userScoreLogDao{db: db}
}

type userScoreLogDao struct {
	db *gorm.DB
}

func (d *userScoreLogDao) Get(id int64) *model.UserScoreLog {
	ret := &model.UserScoreLog{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *userScoreLogDao) Take(where ...interface{}) *model.UserScoreLog {
	ret := &model.UserScoreLog{}
	if err := d.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *userScoreLogDao) Find(cnd *querybuilder.QueryBuilder) (list []model.UserScoreLog) {
	cnd.Find(d.db, &list)
	return
}

func (d *userScoreLogDao) FindOne(cnd *querybuilder.QueryBuilder) *model.UserScoreLog {
	ret := &model.UserScoreLog{}
	if err := cnd.FindOne(d.db, ret); err != nil {
		return nil
	}
	return ret
}

func (d *userScoreLogDao) List(cnd *querybuilder.QueryBuilder) (list []model.UserScoreLog, paging *querybuilder.Paging) {
	cnd.Find(d.db, &list)
	count := cnd.Count(d.db, &model.UserScoreLog{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (d *userScoreLogDao) Create(t *model.UserScoreLog) (err error) {
	err = d.db.Create(t).Error
	return
}

func (d *userScoreLogDao) Update(t *model.UserScoreLog) (err error) {
	err = d.db.Save(t).Error
	return
}

func (d *userScoreLogDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = d.db.Model(&model.UserScoreLog{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *userScoreLogDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = d.db.Model(&model.UserScoreLog{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *userScoreLogDao) Delete(id int64) {
	d.db.Delete(&model.UserScoreLog{}, "id = ?", id)
}