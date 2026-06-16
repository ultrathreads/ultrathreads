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

// Controller 层依赖的业务接口
type NodeServicer interface {
    Get(id int64) *model.Node
    GetBySlug(slug string) *model.Node
    List(cnd *querybuilder.QueryBuilder) ([]model.Node, *querybuilder.Paging)
    Create(req form.NodeCreateForm) (*model.Node, error)
    Update(id int64, req form.NodeUpdateForm) error
    Delete(id int64) error
    IncrTopicCount(nodeId int64)
    GetRecommendNodes() []model.Node
    GetNodes() []model.Node
}

func NewNodeService(repo dao.NodeRepository, nodeCache cache.NodeCacheInterface) *nodeService {
    return &nodeService{
    	repo: repo,
    	nodeCache: nodeCache,
    }
}

type nodeService struct{
	repo dao.NodeRepository
	nodeCache cache.NodeCacheInterface
}

func (s *nodeService) Get(id int64) *model.Node {
	if node := s.nodeCache.Get(id); node != nil {
		return node
	}
	return s.repo.Get(id)
}

func (s *nodeService) GetBySlug(slug string) *model.Node {
	id := hashid.Slug2Id[model.Node](slug)
	return s.Get(id)
}

func (s *nodeService) List(cnd *querybuilder.QueryBuilder) (list []model.Node, paging *querybuilder.Paging) {
	return s.repo.List(cnd)
}

func (s *nodeService) Create(req form.NodeCreateForm) (*model.Node, error) {
	node := &model.Node{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		SortNo:      req.SortNo,
		Status:      req.Status,
		CreateTime:  util.NowTimestamp(),
	}
	if err := s.repo.Create(node); err != nil {
		return nil, errors.New("创建节点失败")
	}
	s.nodeCache.InvalidateAll()
	return node, nil
}

func (s *nodeService) Update(id int64, req form.NodeUpdateForm) error {
	err := s.repo.Updates(id, map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"icon": 	   req.Icon,
		"sort_no":     req.SortNo,
		"status":      req.Status,
	})
	if err != nil {
		return errors.New("更新节点失败")
	}
	s.nodeCache.InvalidateAll()
	return nil
}

func (s *nodeService) Delete(id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return errors.New("删除节点失败")
	}
	s.nodeCache.InvalidateAll()
	return nil
}

// IncrTopicCount 主题数+1（高频写操作，仅失效单条缓存）
func (s *nodeService) IncrTopicCount(nodeId int64) {
	if err := s.repo.IncrField(nodeId, "topic_count", 1); err != nil {
		log.Error("IncrTopicCount failed: nodeId=%d, err=%v", nodeId, err)
		return
	}
	s.nodeCache.Invalidate(nodeId)
}

func (s *nodeService) GetRecommendNodes() []model.Node {
	return s.repo.Find(
		querybuilder.NewQueryBuilder().
			Eq("status", model.StatusOk).
			Asc("sort_no").
			Desc("id").
			Limit(3),
	)
}

func (s *nodeService) GetNodes() []model.Node {
	// 获取全量节点可以走 GetAll 缓存
	if nodes := s.nodeCache.GetAll(); len(nodes) > 0 {
		return nodes
	}

	return s.repo.Find(
		querybuilder.NewQueryBuilder().
			Eq("status", model.StatusOk).
			Asc("sort_no").
			Desc("id"),
	)
}