package converter

import (
	"ultrathreads/model"
	"ultrathreads/util/hashid"
)

func ToNode(node *model.Node) *model.NodeResponse {
	if node == nil {
		return nil
	}
	slug, _ := hashid.Encode[model.Node](node.ID)
	return &model.NodeResponse{
		NodeId:      node.ID,
		Slug:		 slug,
		Name:        node.Name,
		Description: node.Description,
		Icon:        node.Icon,
		TopicCount:  node.TopicCount,
	}
}

// ToNodes 返回 []model.NodeResponse（非指针），与 Response 结构体对齐
func ToNodes(nodes []model.Node) []model.NodeResponse {
	if len(nodes) == 0 {
		return []model.NodeResponse{}
	}
	responses := make([]model.NodeResponse, 0, len(nodes))
	for i := range nodes {
		if r := ToNode(&nodes[i]); r != nil {
			responses = append(responses, *r)
		}
	}
	return responses
}