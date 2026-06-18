package repository

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// UserScoreLogRepository 用户积分日志数据访问契约
type UserScoreLogRepository interface {
	Get(id int64) *model.UserScoreLog
	Take(where ...interface{}) *model.UserScoreLog
	Find(cnd *querybuilder.QueryBuilder) []model.UserScoreLog
	FindOne(cnd *querybuilder.QueryBuilder) *model.UserScoreLog
	List(cnd *querybuilder.QueryBuilder) ([]model.UserScoreLog, *querybuilder.Paging)
	Create(t *model.UserScoreLog) error
	Update(t *model.UserScoreLog) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64)
}

type userScoreLogRepo struct {
	db *gorm.DB
}

func NewUserScoreLogRepository(db *gorm.DB) UserScoreLogRepository {
	return &userScoreLogRepo{db: db}
}

func (r *userScoreLogRepo) Get(id int64) *model.UserScoreLog {
	ret := &model.UserScoreLog{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userScoreLogRepo) Take(where ...interface{}) *model.UserScoreLog {
	ret := &model.UserScoreLog{}
	if err := r.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userScoreLogRepo) Find(cnd *querybuilder.QueryBuilder) (list []model.UserScoreLog) {
	cnd.Find(r.db, &list)
	return
}

func (r *userScoreLogRepo) FindOne(cnd *querybuilder.QueryBuilder) *model.UserScoreLog {
	ret := &model.UserScoreLog{}
	if err := cnd.FindOne(r.db, ret); err != nil {
		return nil
	}
	return ret
}

func (r *userScoreLogRepo) List(cnd *querybuilder.QueryBuilder) (list []model.UserScoreLog, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.UserScoreLog{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (r *userScoreLogRepo) Create(t *model.UserScoreLog) error {
	return r.db.Create(t).Error
}

func (r *userScoreLogRepo) Update(t *model.UserScoreLog) error {
	return r.db.Save(t).Error
}

func (r *userScoreLogRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.UserScoreLog{}).Where("id = ?", id).Updates(columns).Error
}

func (r *userScoreLogRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.UserScoreLog{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (r *userScoreLogRepo) Delete(id int64) {
	r.db.Delete(&model.UserScoreLog{}, "id = ?", id)
}
