package dao

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

func NewNotificationDao(db *gorm.DB) *notificationDao {
	return &notificationDao{db: db}
}

type notificationDao struct {
	db *gorm.DB
}

func (d *notificationDao) Get(id int64) *model.Notification {
	ret := &model.Notification{}
	if err := d.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *notificationDao) Take(where ...interface{}) *model.Notification {
	ret := &model.Notification{}
	if err := d.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *notificationDao) Find(cnd *querybuilder.QueryBuilder) (list []model.Notification) {
	cnd.Find(d.db, &list)
	return
}

func (d *notificationDao) FindOne(cnd *querybuilder.QueryBuilder) *model.Notification {
	ret := &model.Notification{}
	if err := cnd.FindOne(d.db, ret); err != nil {
		return nil
	}
	return ret
}

func (d *notificationDao) List(cnd *querybuilder.QueryBuilder) (list []model.Notification, paging *querybuilder.Paging) {
	cnd.Find(d.db, &list)
	count := cnd.Count(d.db, &model.Notification{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (d *notificationDao) Create(t *model.Notification) (err error) {
	err = d.db.Create(t).Error
	return
}

func (d *notificationDao) Update(t *model.Notification) (err error) {
	err = d.db.Save(t).Error
	return
}

func (d *notificationDao) Updates(id int64, columns map[string]interface{}) (err error) {
	err = d.db.Model(&model.Notification{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (d *notificationDao) UpdateColumn(id int64, name string, value interface{}) (err error) {
	err = d.db.Model(&model.Notification{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (d *notificationDao) GetUnReadCount(userId int64) (count int64) {
	d.db.Model(&model.Notification{}).Where("user_id = ? and status = ?", userId, model.NotificationStatusUnread).Count(&count)
	return
}

func (d *notificationDao) UpdateStatusBatch(userId int64) (err error) {
	err = d.db.Model(&model.Notification{}).Where("user_id = ? and status = ?", userId, model.NotificationStatusUnread).Updates(model.Notification{Status: model.NotificationStatusReaded}).Error
	return
}

func (d *notificationDao) Delete(id int64) {
	d.db.Delete(&model.Notification{}, "id = ?", id)
}