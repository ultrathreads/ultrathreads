package dao

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ultrathreads/model"
)

func NewUserReadStateDao(db *gorm.DB) *userReadStateDao {
	return &userReadStateDao{db: db}
}

type userReadStateDao struct {
	db *gorm.DB
}

// GetLastReadAt 获取用户在指定节点的已读时间戳
func (d *userReadStateDao) GetLastReadAt(userID, nodeID int64) int64 {
	var record model.UserReadState
	err := d.db.Where("user_id = ? AND node_id = ?", userID, nodeID).First(&record).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			// TODO: 建议接入日志框架记录真实 DB 错误
			// log.Error("GetLastReadAt failed: userId=%d, nodeId=%d, err=%v", userID, nodeID, err)
		}
		// 未找到记录或其他错误均返回 0，由上层 LoadingCache 缓存该零值防穿透
		return 0
	}
	return record.LastReadAt
}

// Upsert 插入或更新已读状态（仅向前推进游标）
func (d *userReadStateDao) Upsert(userID, nodeID int64, readAt int64) error {
	record := model.UserReadState{
		UserID:     userID,
		NodeID:     nodeID,
		LastReadAt: readAt,
	}

	// GREATEST + COALESCE 确保游标只增不减，且兼容 last_read_at 为 NULL 的情况
	result := d.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}, {Name: "node_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"last_read_at": gorm.Expr("GREATEST(COALESCE(last_read_at, 0), ?)", readAt),
		}),
	}).Create(&record)

	if result.Error != nil {
		return fmt.Errorf("upsert user read state failed: %w", result.Error)
	}
	return nil
}

// GetAllReadStates 获取用户所有已读状态，返回 map[nodeID]lastReadAt
// 未找到记录时返回空 map（非 nil），便于上层缓存直接使用
func (d *userReadStateDao) GetAllReadStates(userID int64) map[int64]int64 {
	var records []model.UserReadState
	if err := d.db.Where("user_id = ?", userID).Find(&records).Error; err != nil {
		// TODO: 建议接入日志框架记录真实 DB 错误
		// log.Error("GetAllReadStates failed: userId=%d, err=%v", userID, err)
		return make(map[int64]int64)
	}

	states := make(map[int64]int64, len(records))
	for _, r := range records {
		states[r.NodeID] = r.LastReadAt
	}
	return states
}

// DeleteByUser 清除用户所有已读状态（注销/重置时使用）
func (d *userReadStateDao) DeleteByUser(userID int64) error {
	if err := d.db.Where("user_id = ?", userID).Delete(&model.UserReadState{}).Error; err != nil {
		return fmt.Errorf("delete user read state failed: %w", err)
	}
	return nil
}

// GetUnreadNodeIDs 获取用户在指定节点列表中仍有未读内容的节点ID
func (d *userReadStateDao) GetUnreadNodeIDs(userID int64, nodeIDs []int64, postCreatedAtMap map[int64]int64) ([]int64, error) {
	if len(nodeIDs) == 0 {
		return nil, nil
	}

	var records []model.UserReadState
	if err := d.db.Where("user_id = ? AND node_id IN (?)", userID, nodeIDs).Find(&records).Error; err != nil {
		return nil, fmt.Errorf("query user read states failed: %w", err)
	}

	readMap := make(map[int64]int64, len(records))
	for _, r := range records {
		readMap[r.NodeID] = r.LastReadAt
	}

	var unreadNodeIDs []int64
	for _, nodeID := range nodeIDs {
		lastReadAt, exists := readMap[nodeID]
		postCreatedAt := postCreatedAtMap[nodeID]
		// 无已读记录 或 帖子发布时间晚于已读时间 → 未读
		if !exists || postCreatedAt > lastReadAt {
			unreadNodeIDs = append(unreadNodeIDs, nodeID)
		}
	}
	return unreadNodeIDs, nil
}