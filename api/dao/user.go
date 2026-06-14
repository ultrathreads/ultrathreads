package dao

import (
	"errors"

	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

var UserDao = newUserDao()

func newUserDao() *userDao {
	return &userDao{}
}

type userDao struct{}

// Get 根据 ID 获取用户，未找到返回 nil
func (d *userDao) Get(id int64) *model.User {
	ret := &model.User{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

// Take 按条件获取单条记录（无排序保证），未找到返回 nil
func (d *userDao) Take(where ...interface{}) *model.User {
	ret := &model.User{}
	if err := db.Take(ret, where...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *userDao) Find(cnd *querybuilder.QueryBuilder) (list []model.User) {
	cnd.Find(db, &list)
	return
}

// FindOne 通过 QueryBuilder 查询单条记录
func (d *userDao) FindOne(cnd *querybuilder.QueryBuilder) *model.User {
	ret := &model.User{}
	if err := cnd.FindOne(db, ret); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return nil
	}
	return ret
}

func (d *userDao) List(cnd *querybuilder.QueryBuilder) (list []model.User, paging *querybuilder.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.User{})

	paging = &querybuilder.Paging{
		Page:  cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total: count,
	}
	return
}

// Count 统计数量
func (d *userDao) Count(cnd *querybuilder.QueryBuilder) int64 {
	return cnd.Count(db, &model.User{})
}

func (d *userDao) Create(t *model.User) error {
	return db.Create(t).Error
}

func (d *userDao) Update(t *model.User) error {
	return db.Save(t).Error
}

func (d *userDao) Updates(id int64, columns map[string]interface{}) error {
	return db.Model(&model.User{}).Where("id = ?", id).Updates(columns).Error
}

func (d *userDao) UpdateColumn(id int64, name string, value interface{}) error {
	return db.Model(&model.User{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

// Delete 根据 ID 删除
func (d *userDao) Delete(id int64) error { // ✅ 补充 error 返回值
	return db.Delete(&model.User{}, "id = ?", id).Error
}

// GetByEmail 根据邮箱获取用户
func (d *userDao) GetByEmail(email string) *model.User {
	return d.Take("email = ?", email)
}

// GetByUsername 根据用户名获取用户
func (d *userDao) GetByUsername(username string) *model.User {
	return d.Take("username = ?", username)
}