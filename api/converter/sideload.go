package converter

import (
	"ultrathreads/dao"
	"ultrathreads/model"
	//"ultrathreads/util"
	"ultrathreads/service"
	"ultrathreads/util/hashid"
)

// ToSimplePostsWithIncluded 适配 slug 为主键的 sideload
func ToSimplePostsWithIncluded(posts []model.Post) (
	[]model.PostItem,
	[]model.UserIncluded,
	[]model.NodeIncluded,
) {
	if len(posts) == 0 {
		return nil, nil, nil
	}

	// 1. 收集【数据库ID】（去重，用于批量查库）
	var (
		userIDs = make(map[int64]struct{})
		nodeIDs = make(map[int64]struct{})
	)
	for _, p := range posts {
		userIDs[p.UserId] = struct{}{}
		nodeIDs[p.NodeId] = struct{}{}

		// 扁平化回帖 同理收集回帖人ID
		// for _, reply := range p.Replies {
		// 	userIDs[reply.UserID] = struct{}{}
		// }
	}

	// 转切片，用于 IN 查询
	uidList := make([]int64, 0, len(userIDs))
	for id := range userIDs {
		uidList = append(uidList, id)
	}
	nidList := make([]int64, 0, len(nodeIDs))
	for id := range nodeIDs {
		nidList = append(nidList, id)
	}

	// 2. 批量查库（依旧用 ID，主键索引最快）
	users := dao.UserDao.FindByIds(uidList)
	nodes := dao.NodeDao.FindByIds(nidList)

	// 3. 构建【slug 为 key】的内存索引（前端/侧载专用）
	var (
		userSlugMap = make(map[string]model.UserIncluded)
		nodeSlugMap = make(map[string]model.NodeIncluded)
	)

	// 用户映射：id → slug → 存入map
	for _, u := range users {
		slug := hashid.Id2Slug[model.User](u.ID)
		inc := model.UserIncluded{
			Slug:   slug,
			Name:   u.Nickname,
			Avatar: u.Avatar,
		}
		userSlugMap[slug] = inc
	}

	// 板块映射
	for _, n := range nodes {
		slug := hashid.Id2Slug[model.Node](n.ID)
		inc := model.NodeIncluded{
			Slug: slug,
			Name: n.Name,
		}
		nodeSlugMap[slug] = inc
	}

	// 4. 转换主列表 DTO
	respList := make([]model.PostItem, 0, len(posts))
	for _, p := range posts {
		postSlug := hashid.Id2Slug[model.Post](p.ID)
		userSlug := hashid.Id2Slug[model.User](p.UserId)
		nodeSlug := hashid.Id2Slug[model.Node](p.NodeId)

		parentSlug := hashid.Id2Slug[model.Post](p.ParentId)
		threadSlug := hashid.Id2Slug[model.Post](p.ThreadId)

		rsp := model.PostItem{
			Slug:       postSlug,
			ParentSlug: parentSlug,
			ThreadSlug: threadSlug,
			UserSlug:   userSlug,
			NodeSlug:   nodeSlug,
			Title:      p.Title,
			User:       ToUserDefaultIfNull(p.UserId),
		}

		rsp.LastCommentTime = p.LastCommentTime
		rsp.CreateTime = p.CreateTime

		if p.NodeId > 0 {
			node := service.NodeService.Get(p.NodeId)
			rsp.Node = ToNode(node)
		}
		respList = append(respList, rsp)
	}

	// 5. slug map 转切片，作为 sideload included
	incUsers := make([]model.UserIncluded, 0, len(userSlugMap))
	for _, u := range userSlugMap {
		incUsers = append(incUsers, u)
	}
	incNodes := make([]model.NodeIncluded, 0, len(nodeSlugMap))
	for _, n := range nodeSlugMap {
		incNodes = append(incNodes, n)
	}

	return respList, incUsers, incNodes
}