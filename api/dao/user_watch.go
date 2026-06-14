package dao

import (
	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var UserWatchDao = newUserWatchDao()

func newUserWatchDao() *userWatchDao {
	return &userWatchDao{}
}

type userWatchDao struct{}

func (d *userWatchDao) Get(id int64) *model.UserWatch {
	ret := &model.UserWatch{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (d *userWatchDao) Take(where ...interface{}) *model.UserWatch {
	ret := &model.UserWatch{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (d *userWatchDao) Find(cnd *querybuilder.QueryBuilder) []model.UserWatch {
	var list []model.UserWatch
	cnd.Find(db, &list)
	return list
}

func (d *userWatchDao) FindOne(cnd *querybuilder.QueryBuilder) *model.UserWatch {
	ret := &model.UserWatch{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (d *userWatchDao) List(cnd *querybuilder.QueryBuilder) ([]model.UserWatch, *querybuilder.Paging) {
	var list []model.UserWatch
	cnd.Find(db, &list)

	// ✅ v2 Count 不再接受指针参数，改为直接返回 (int64, error) 或通过链式调用赋值
	count := cnd.Count(db, &model.UserWatch{})

	paging := &querybuilder.Paging{
		Page:  cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total: count,
	}
	return list, paging
}

func (d *userWatchDao) Create(t *model.UserWatch) error {
	return db.Create(t).Error
}

func (d *userWatchDao) Update(t *model.UserWatch) error {
	return db.Save(t).Error
}

func (d *userWatchDao) Updates(id int64, columns map[string]interface{}) error {
	return db.Model(&model.UserWatch{}).Where("id = ?", id).Updates(columns).Error
}

func (d *userWatchDao) UpdateColumn(id int64, name string, value interface{}) error {
	return db.Model(&model.UserWatch{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

// Delete 删除关注记录
func (d *userWatchDao) Delete(id int64) error { // ✅ 补充 error 返回值
	return db.Delete(&model.UserWatch{}, "id = ?", id).Error
}