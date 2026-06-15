package converter

import (
	"ultrathreads/dao"
	"ultrathreads/model"
	//"ultrathreads/util"
	//"ultrathreads/service"
	"ultrathreads/util/hashid"
)

// ToSimplePostsWithIncluded 适配 slug 为主键的 sideload
func ToSimplePostsWithIncluded(posts []model.Post) (
	[]model.PostItem,
	[]model.UserIncluded,
	[]model.NodeIncluded,
	[]model.TagIncluded,
) {
	if len(posts) == 0 {
		return nil, nil, nil, nil
	}

	// 1. 收集【数据库ID】（去重，用于批量查库）
	var (
		userIDs = make(map[int64]struct{})
		nodeIDs = make(map[int64]struct{})
		postIDs = make([]int64, 0, len(posts)) // 【关键】收集帖子ID用于查帖子和标签的关联表
	)
	for _, p := range posts {
		userIDs[p.UserId] = struct{}{}
		nodeIDs[p.NodeId] = struct{}{}
		postIDs = append(postIDs, p.ID)

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

	// 【关键】通过关联表批量查出 postId → tagIds 的映射
	// 返回 map[int64][]int64 即 postId → []tagId
	postTagMap := dao.TagDao.FindTagIdsByPostIds(postIDs)

	// 从映射中收集所有去重的 tagID，再批量查 Tag 详情
	tagIDs := make(map[int64]struct{})
	for _, tids := range postTagMap {
		for _, tid := range tids {
			tagIDs[tid] = struct{}{}
		}
	}
	tags := dao.TagDao.FindByIds(mapKeys(tagIDs))

	// ========== 3. 构建 slug 索引 ==========
	userSlugMap := make(map[string]model.UserIncluded, len(users))
	for _, u := range users {
		slug := hashid.Id2Slug[model.User](u.ID)
		userSlugMap[slug] = model.UserIncluded{
			Slug: slug, Username: u.Username.String,
			Nickname: u.Nickname, Avatar: u.Avatar,
		}
	}

	nodeSlugMap := make(map[string]model.NodeIncluded, len(nodes))
	for _, n := range nodes {
		slug := hashid.Id2Slug[model.Node](n.ID)
		nodeSlugMap[slug] = model.NodeIncluded{Slug: slug, Name: n.Name}
	}

	tagSlugMap := make(map[string]model.TagIncluded, len(tags))
	tagIdToSlug := make(map[int64]string, len(tags)) // 【新增】tagID→slug 反查
	for _, t := range tags {
		slug := hashid.Id2Slug[model.Tag](t.ID)
		tagSlugMap[slug] = model.TagIncluded{Slug: slug, Name: t.Name}
		tagIdToSlug[t.ID] = slug
	}

	// 4. 转换主列表 DTO
	respList := make([]model.PostItem, 0, len(posts))
	for _, p := range posts {
		postSlug := hashid.Id2Slug[model.Post](p.ID)
		userSlug := hashid.Id2Slug[model.User](p.UserId)
		nodeSlug := hashid.Id2Slug[model.Node](p.NodeId)

		parentSlug := hashid.Id2Slug[model.Post](p.ParentId)
		threadSlug := hashid.Id2Slug[model.Post](p.ThreadId)

		// 通过 postTagMap + tagIdToSlug 组装该帖子的 TagSlugs
		var tagSlugs []string
		if tids, ok := postTagMap[p.ID]; ok {
			tagSlugs = make([]string, 0, len(tids))
			for _, tid := range tids {
				if slug, exists := tagIdToSlug[tid]; exists {
					tagSlugs = append(tagSlugs, slug)
				}
			}
		} else {
			tagSlugs = []string{} // 保证返回 [] 而非 null
		}

		rsp := model.PostItem{
			Slug:       postSlug,
			ParentSlug: parentSlug,
			ThreadSlug: threadSlug,
			UserSlug:   userSlug,
			NodeSlug:   nodeSlug,
			TagSlugs:        tagSlugs,
			CreateTime:      p.CreateTime,
			Title:      p.Title,
			LastCommentTime: p.LastCommentTime,

			IsRoot:   p.IsRoot(),
			IsPinned: p.IsPinned,
		}

		rsp.LastCommentTime = p.LastCommentTime
		rsp.CreateTime = p.CreateTime
		respList = append(respList, rsp)
	}

	// 5. slug map 转切片，作为 sideload included
	return respList,
		mapValues(userSlugMap),
		mapValues(nodeSlugMap),
		mapValues(tagSlugMap)
}

func mapKeys(m map[int64]struct{}) []int64 {
	keys := make([]int64, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func mapValues[T any](m map[string]T) []T {
	vals := make([]T, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}
