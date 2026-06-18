package dao

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// NotificationRepository 通知数据访问契约
type NotificationRepository interface {
	Get(id int64) *model.Notification
	Take(where ...interface{}) *model.Notification
	Find(cnd *querybuilder.QueryBuilder) []model.Notification
	FindOne(cnd *querybuilder.QueryBuilder) *model.Notification
	List(cnd *querybuilder.QueryBuilder) ([]model.Notification, *querybuilder.Paging)
	Create(t *model.Notification) error
	Update(t *model.Notification) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	GetUnReadCount(userId int64) int64
	UpdateStatusBatch(userId int64) error
	Delete(id int64)
}

type notificationRepo struct {
	db *gorm.DB
}

func NewNotificationDao(db *gorm.DB) NotificationRepository {
	return &notificationRepo{db: db}
}

func (r *notificationRepo) Get(id int64) *model.Notification {
	ret := &model.Notification{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *notificationRepo) Take(where ...interface{}) *model.Notification {
	ret := &model.Notification{}
	if err := r.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *notificationRepo) Find(cnd *querybuilder.QueryBuilder) (list []model.Notification) {
	cnd.Find(r.db, &list)
	return
}

func (r *notificationRepo) FindOne(cnd *querybuilder.QueryBuilder) *model.Notification {
	ret := &model.Notification{}
	if err := cnd.FindOne(r.db, ret); err != nil {
		return nil
	}
	return ret
}

func (r *notificationRepo) List(cnd *querybuilder.QueryBuilder) (list []model.Notification, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.Notification{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (r *notificationRepo) Create(t *model.Notification) error {
	return r.db.Create(t).Error
}

func (r *notificationRepo) Update(t *model.Notification) error {
	return r.db.Save(t).Error
}

func (r *notificationRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.Notification{}).Where("id = ?", id).Updates(columns).Error
}

func (r *notificationRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.Notification{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (r *notificationRepo) GetUnReadCount(userId int64) (count int64) {
	r.db.Model(&model.Notification{}).Where("user_id = ? and status = ?", userId, model.NotificationStatusUnread).Count(&count)
	return
}

func (r *notificationRepo) UpdateStatusBatch(userId int64) error {
	return r.db.Model(&model.Notification{}).Where("user_id = ? and status = ?", userId, model.NotificationStatusUnread).Updates(model.Notification{Status: model.NotificationStatusReaded}).Error
}

func (r *notificationRepo) Delete(id int64) {
	r.db.Delete(&model.Notification{}, "id = ?", id)
}