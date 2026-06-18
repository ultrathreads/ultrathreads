package dao

import (
	"gorm.io/gorm"

	"ultrathreads/model"
	"ultrathreads/util/querybuilder"
)

// UserReadStateRepository 用户阅读状态数据访问契约
type UserReadStateRepository interface {
	Get(id int64) *model.UserReadState
	Take(where ...interface{}) *model.UserReadState
	Find(cnd *querybuilder.QueryBuilder) []model.UserReadState
	FindOne(cnd *querybuilder.QueryBuilder) *model.UserReadState
	List(cnd *querybuilder.QueryBuilder) ([]model.UserReadState, *querybuilder.Paging)
	Create(t *model.UserReadState) error
	Update(t *model.UserReadState) error
	Updates(id int64, columns map[string]interface{}) error
	UpdateColumn(id int64, name string, value interface{}) error
	Delete(id int64)
	GetLastReadAt(userID, nodeID int64) int64
	GetAllReadStates(userID int64) map[int64]int64
	Upsert(userID, nodeID, lastReadAt int64) error
}

type userReadStateRepo struct {
	db *gorm.DB
}

func NewUserReadStateDao(db *gorm.DB) UserReadStateRepository {
	return &userReadStateRepo{db: db}
}

func (r *userReadStateRepo) Get(id int64) *model.UserReadState {
	ret := &model.UserReadState{}
	if err := r.db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userReadStateRepo) Take(where ...interface{}) *model.UserReadState {
	ret := &model.UserReadState{}
	if err := r.db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userReadStateRepo) Find(cnd *querybuilder.QueryBuilder) (list []model.UserReadState) {
	cnd.Find(r.db, &list)
	return
}

func (r *userReadStateRepo) FindOne(cnd *querybuilder.QueryBuilder) *model.UserReadState {
	ret := &model.UserReadState{}
	if err := cnd.FindOne(r.db, ret); err != nil {
		return nil
	}
	return ret
}

func (r *userReadStateRepo) List(cnd *querybuilder.QueryBuilder) (list []model.UserReadState, paging *querybuilder.Paging) {
	cnd.Find(r.db, &list)
	count := cnd.Count(r.db, &model.UserReadState{})

	paging = &querybuilder.Paging{
		Page:     cnd.Paging.Page,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
	}
	return
}

func (r *userReadStateRepo) Create(t *model.UserReadState) error {
	return r.db.Create(t).Error
}

func (r *userReadStateRepo) Update(t *model.UserReadState) error {
	return r.db.Save(t).Error
}

func (r *userReadStateRepo) Updates(id int64, columns map[string]interface{}) error {
	return r.db.Model(&model.UserReadState{}).Where("id = ?", id).Updates(columns).Error
}

func (r *userReadStateRepo) UpdateColumn(id int64, name string, value interface{}) error {
	return r.db.Model(&model.UserReadState{}).Where("id = ?", id).UpdateColumn(name, value).Error
}

func (r *userReadStateRepo) Delete(id int64) {
	r.db.Delete(&model.UserReadState{}, "id = ?", id)
}

// GetLastReadAt 获取用户在某节点的最后阅读时间，未找到返回 0
func (r *userReadStateRepo) GetLastReadAt(userID, nodeID int64) int64 {
	var state model.UserReadState
	if err := r.db.Where("user_id = ? AND node_id = ?", userID, nodeID).First(&state).Error; err != nil {
		return 0
	}
	return state.LastReadAt
}

// GetAllReadStates 获取用户所有已读状态，返回 map[nodeID]lastReadAt
func (r *userReadStateRepo) GetAllReadStates(userID int64) map[int64]int64 {
	var states []model.UserReadState
	r.db.Where("user_id = ?", userID).Find(&states)
	result := make(map[int64]int64, len(states))
	for _, s := range states {
		result[s.NodeID] = s.LastReadAt
	}
	return result
}

// Upsert 插入或更新用户阅读状态
func (r *userReadStateRepo) Upsert(userID, nodeID, lastReadAt int64) error {
	var state model.UserReadState
	err := r.db.Where("user_id = ? AND node_id = ?", userID, nodeID).First(&state).Error
	if err == nil {
		return r.db.Model(&state).Update("last_read_at", lastReadAt).Error
	}
	return r.db.Create(&model.UserReadState{
		UserID:     userID,
		NodeID:     nodeID,
		LastReadAt: lastReadAt,
	}).Error
}
