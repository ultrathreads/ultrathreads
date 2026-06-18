package mock

import (
	"fmt"
	"ultrathreads/model"
)

// UpdateNodeTopicCount 遍历所有节点，逐个更新主题数
// 主题数 = 该 nodeId 下 parent_id == 0 的 post 数量
func UpdateNodeTopicCount() {
	var nodes []*model.Node
	if err := mockDB.Find(&nodes).Error; err != nil {
		fmt.Printf("[mock] failed to fetch nodes: %v\n", err)
		return
	}

	updated := 0
	for _, node := range nodes {
		var count int64
		err := mockDB.Model(&model.Post{}).
			Where("node_id = ? AND parent_id = ?", node.ID, 0).
			Count(&count).Error

		if err != nil {
			fmt.Printf("[mock] failed to count topics for node %d: %v\n", node.ID, err)
			continue
		}

		// ✅ 使用 topic_count
		if err := mockDB.Model(node).Update("topic_count", count).Error; err != nil {
			fmt.Printf("[mock] failed to update topic_count for node %d: %v\n", node.ID, err)
			continue
		}
		updated++
	}

	fmt.Printf("[mock] updated topic_count for %d/%d nodes\n", updated, len(nodes))
}
