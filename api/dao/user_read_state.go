package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"ultrathreads/model"
)

var UserReadStateDao = newUserReadStateDao()

func newUserReadStateDao() *userReadStateDao {
	return &userReadStateDao{}
}

type userReadStateDao struct {
}

// GetLastReadAt 获取用户在指定节点的已读时间戳
func (d *userReadStateDao) GetLastReadAt(userID, nodeID int64) int64 {
	var record model.UserReadState
	err := db.Where("user_id = ? AND node_id = ?", userID, nodeID).First(&record).Error
	if err != nil {
		// 未找到记录或其他错误均返回 0，由上层 LoadingCache 缓存该零值防穿透
		return 0
	}
	return record.LastReadAt
}

// Upsert 插入或更新已读状态（仅向前推进游标）
// 兼容 GORM v1，不使用 Clauses
func (d *userReadStateDao) Upsert(userID, nodeID int64, readAt int64) error {
	var record model.UserReadState

	record.UserID = userID
    record.NodeID = nodeID

	err := db.Where("user_id = ? AND node_id = ?", userID, nodeID).First(&record).Error

	if err == gorm.ErrRecordNotFound {
        newRecord := model.UserReadState{
            UserID:     userID,
            NodeID:     nodeID,
            LastReadAt: readAt,
        }
        return db.Create(&newRecord).Error
    }
    if err != nil {
        return fmt.Errorf("query read state failed: %w", err)
    }

	if readAt > record.LastReadAt {
        return db.Model(&model.UserReadState{}).
            Where("user_id = ? AND node_id = ?", userID, nodeID).
            Updates(map[string]interface{}{"last_read_at": readAt}).Error
    }

	return nil
}

// DeleteByUser 清除用户所有已读状态（注销/重置时使用）
func (d *userReadStateDao) DeleteByUser(userID int64) error {
	return db.Where("user_id = ?", userID).Delete(&model.UserReadState{}).Error
}

// GetUnreadNodeIDs 获取用户在指定节点列表中仍有未读内容的节点ID
// 供首页/列表页批量判断使用
func (d *userReadStateDao) GetUnreadNodeIDs(userID int64, nodeIDs []int64, postCreatedAtMap map[int64]int64) []int64 {
	if len(nodeIDs) == 0 {
		return nil
	}

	var records []model.UserReadState
	db.Where("user_id = ? AND node_id IN (?)", userID, nodeIDs).Find(&records)

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
	return unreadNodeIDs
}