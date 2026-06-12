package service

import (
	"errors"

	"ultrathreads/cache"
	"ultrathreads/dao"
	"ultrathreads/form"
	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/hashid"
	"ultrathreads/util/log"
	"ultrathreads/util/querybuilder"
)

var NodeService = newNodeService()

func newNodeService() *nodeService {
	return &nodeService{}
}

type nodeService struct{}

func (s *nodeService) Get(id int64) *model.Node {
	return dao.NodeDao.Get(id)
}

func (s *nodeService) GetBySlug(slug string) *model.Node {
	id,_ := hashid.Decode[model.Node](slug)
	return dao.NodeDao.Get(id)
}

func (s *nodeService) List(cnd *querybuilder.QueryBuilder) (list []model.Node, paging *querybuilder.Paging) {
	return dao.NodeDao.List(cnd)
}

func (s *nodeService) Create(dto form.NodeCreateForm) (*model.Node, error) {
	node := &model.Node{
		Name:        dto.Name,
		Description: dto.Description,
		SortNo:      dto.SortNo,
		Status:      dto.Status,
		CreateTime:  util.NowTimestamp(),
	}
	if err := dao.NodeDao.Create(node); err != nil {
		return nil, errors.New("创建节点失败")
	}
	cache.NodeCache.InvalidateAll()
	return node, nil
}

func (s *nodeService) Update(dto form.NodeUpdateForm) error {
	err := dao.NodeDao.Updates(dto.ID, map[string]interface{}{
		"name":        dto.Name,
		"description": dto.Description,
		"sort_no":     dto.SortNo,
		"status":      dto.Status,
		"update_time": util.NowTimestamp(),
	})
	if err != nil {
		return errors.New("更新节点失败")
	}
	cache.NodeCache.InvalidateAll()
	return nil
}

func (s *nodeService) Delete(id int64) error {
	if err := dao.NodeDao.Delete(id); err != nil {
		return errors.New("删除节点失败")
	}
	cache.NodeCache.InvalidateAll()
	return nil
}

// IncrTopicCount 主题数+1（高频写操作，仅失效单条缓存）
func (s *nodeService) IncrTopicCount(nodeId int64) {
	if err := dao.NodeDao.IncrField(nodeId, "topic_count", 1); err != nil {
		log.Error("IncrTopicCount failed: nodeId=%d, err=%v", nodeId, err)
		return
	}
	cache.NodeCache.Invalidate(nodeId)
}

func (s *nodeService) GetRecommendNodes() []model.Node {
	return dao.NodeDao.Find(
		querybuilder.NewQueryBuilder().
			Eq("status", model.StatusOk).
			Asc("sort_no").
			Desc("id").
			Limit(3),
	)
}

func (s *nodeService) GetNodes() []model.Node {
	return dao.NodeDao.Find(
		querybuilder.NewQueryBuilder().
			Eq("status", model.StatusOk).
			Asc("sort_no").
			Desc("id"),
	)
}