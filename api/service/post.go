package service

import (
	"errors"
	"fmt"
	"math"
	"path"
	"time"

	"github.com/gorilla/feeds"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"ultrathreads/cache"
	"ultrathreads/dao"
	"ultrathreads/form"
	"ultrathreads/dto"
	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/log"
	"ultrathreads/util/hashid"
	"ultrathreads/util/querybuilder"
	"ultrathreads/util/urls"
)

type ScanPostCallback func(posts []model.Post)

func NewPostService(repo dao.PostRepository, nodeRepo dao.NodeRepository) *postService {
    return &postService{repo: repo, nodeRepo: nodeRepo}
}

type postService struct{
	repo dao.PostRepository
	nodeRepo dao.NodeRepository
}

func (s *postService) Get(id int64) *model.Post {
	return dao.PostDao.Get(id)
}

func (s *postService) GetBySlug(slug string) *model.Post {
	id := hashid.Slug2Id[model.Post](slug)
	return dao.PostDao.Get(id)
}

func (s *postService) Find(cnd *querybuilder.QueryBuilder) []model.Post {
	return dao.PostDao.Find(cnd)
}

func (s *postService) List(cnd *querybuilder.QueryBuilder) (list []model.Post, paging *querybuilder.Paging) {
	return dao.PostDao.List(cnd)
}

// Count 统计数量
func (s *postService) Count(cnd *querybuilder.QueryBuilder) int64 {
	return dao.PostDao.Count(cnd)
}

// GetNodeThreadsFull 获取主帖列表（分页）+ 每个主帖下的所有回复（扁平化）
func (s *postService) GetNodeThreadsFull(page, limit int, nodeSlug string) ([]model.Post, *querybuilder.Paging) {
	nodeId := hashid.Slug2Id[model.Node](nodeSlug)
	rootCnd := querybuilder.NewQueryBuilder().
		Eq("parent_id", 0).
		Eq("status", model.StatusOk)

	if nodeId > 0 {
		rootCnd = rootCnd.Eq("node_id", nodeId)
	}

	rootCnd = rootCnd.
		Desc("is_pinned").
		Desc("last_comment_time").
		Page(page, limit)

	rootPosts, paging := dao.PostDao.List(rootCnd)
	if len(rootPosts) == 0 {
		return []model.Post{}, paging
	}

	threadIds := make([]int64, 0, len(rootPosts))
	for _, p := range rootPosts {
		threadIds = append(threadIds, p.ID)
	}

	allCnd := querybuilder.NewQueryBuilder().
		In("thread_id", threadIds).
		Eq("status", model.StatusOk).
		Asc("create_time")

	allPosts := dao.PostDao.Find(allCnd)

	postMap := make(map[int64][]model.Post, len(threadIds))
	for _, p := range allPosts {
		postMap[p.ThreadId] = append(postMap[p.ThreadId], p)
	}

	result := make([]model.Post, 0, len(allPosts))
	for _, tid := range threadIds {
		if posts, ok := postMap[tid]; ok {
			result = append(result, posts...)
		}
	}

	return result, paging
}

// GetTagThreadsFull 获取指定标签下的主帖列表（分页）+ 每个主帖下的所有回复（扁平化）
func (s *postService) GetTagThreadsFull(tagSlug string, page int) (posts []model.Post, paging *querybuilder.Paging) {
	tagId := hashid.Slug2Id[model.Tag](tagSlug)

	// 1. 获取当前页的根帖（通过 IN 子查询过滤标签，保证排序和分页正确）
	subQuery := dao.DB().Model(&model.PostTag{}).
        Select("post_id").
        Where("tag_id = ? AND status = ?", tagId, model.StatusOk)

 	rootCnd := querybuilder.NewQueryBuilder().
        Eq("parent_id", 0).
        Eq("status", model.StatusOk).
        Where("id IN (?)", subQuery). 
        Desc("last_comment_time").
        Page(page, 20)

	rootPosts, paging := dao.PostDao.List(rootCnd)
	if len(rootPosts) == 0 {
		return []model.Post{}, paging
	}

	// 2. 提取当前页所有根帖的 ID
	threadIds := make([]int64, 0, len(rootPosts))
	for _, p := range rootPosts {
		threadIds = append(threadIds, p.ID)
	}

	// 3. 批量查询这些根帖下的所有正常状态的回复（按创建时间正序）
	allCnd := querybuilder.NewQueryBuilder().
		In("thread_id", threadIds).
		Eq("status", model.StatusOk).
		Asc("create_time")

	allPosts := dao.PostDao.Find(allCnd)

	// 4. 将回复按 thread_id 进行分组
	postMap := make(map[int64][]model.Post, len(threadIds))
	for _, p := range allPosts {
		postMap[p.ThreadId] = append(postMap[p.ThreadId], p)
	}

	// 5. 按照根帖的原始顺序，将根帖和对应的回复组装成扁平化切片
	result := make([]model.Post, 0, len(allPosts))
	for _, tid := range threadIds {
		if posts, ok := postMap[tid]; ok {
			result = append(result, posts...)
		}
	}

	return result, paging
}

// GetPostWithThread 获取帖子详情及其所属主题的所有扁平回帖
func (s *postService) GetPostTree(postSlug string) (*model.Post, []model.Post, error) {
	postId := hashid.Slug2Id[model.Post](postSlug)
	if postId <= 0 {
		return nil, nil, errors.New("invalid post_id")
	}

	currentPost := dao.PostDao.Get(postId)
	if currentPost == nil || currentPost.Status != model.StatusOk {
		return nil, nil, errors.New("post not found")
	}

	var posts []model.Post
	if currentPost.ThreadId > 0 {
		cnd := querybuilder.NewQueryBuilder().
			Eq("thread_id", currentPost.ThreadId).
			Eq("status", model.StatusOk).
			Asc("create_time")
		posts = dao.PostDao.Find(cnd)
	}

	if posts == nil {
		posts = []model.Post{}
	}

	return currentPost, posts, nil
}

func (s *postService) GetPostsByThreadId(slug string) ([]model.Post, error) {
	threadId := hashid.Slug2Id[model.Post](slug)
	if threadId <= 0 {
		return nil, errors.New("invalid thread_id")
	}

	cnd := querybuilder.NewQueryBuilder().
		Eq("thread_id", threadId).
		Eq("status", model.StatusOk).
		Asc("create_time")

	posts := dao.PostDao.Find(cnd)
	if posts == nil {
		posts = []model.Post{}
	}

	return posts, nil
}

// GetUserPosts 获取用户帖子列表
func (s *postService) GetUserPosts(userSlug, postType string, page int, pageSize int) ([]model.Post, *querybuilder.Paging) {
	userID := hashid.Slug2Id[model.User](userSlug)
	// 1. 构建基础查询条件
	qb := querybuilder.NewQueryBuilder().
		Eq("user_id", userID).
		Eq("status", model.StatusOk).
		Page(page, pageSize).
		Desc("id")

	// 2. 根据 type 动态追加过滤条件
	switch postType {
	case "reply":
		qb.NotEq("parent_id", 0) 
	case "root":
		fallthrough 
	default:
		qb.Eq("parent_id", 0) 
	}

	// 3. 执行查询
	return dao.PostDao.List(qb)
}

// Delete 软删除帖子
func (s *postService) Delete(id int64) error {
	err := dao.PostDao.UpdateColumn(id, "status", model.StatusDeleted)
	if err == nil {
		PostTagService.DeleteByPostId(id)
	}
	return err
}

// Undelete 取消删除
func (s *postService) Undelete(id int64) error {
	err := dao.PostDao.UpdateColumn(id, "status", model.StatusOk)
	if err == nil {
		PostTagService.UndeleteByPostId(id)
	}
	return err
}

// CreateRootPost 创建根帖（主帖）
func (s *postService) CreateRootPost(userID int64, dto dto.CreateRootPostRequest) (*model.Post, error) {
	nodeID := hashid.Slug2Id[model.Node](dto.NodeSlug)

	// ✅ 节点校验（仅根帖需要）
	if nodeID <= 0 {
		nodeID = SettingService.GetSetting().DefaultNodeId
	}
	if nodeID <= 0 {
		return nil, errors.New("请配置默认节点")
	}
	node := s.nodeRepo.Get(nodeID)
	if node == nil || node.Status != model.StatusOk {
		return nil, errors.New("节点不存在或已禁用")
	}

	now := util.NowTimestamp()
	post := &model.Post{
		Type:            model.PostTypeNormal,
		UserId:          userID,
		NodeId:          nodeID,
		Title:           dto.Title,
		Content:         dto.Content,
		Status:          model.StatusOk,
		LastCommentTime: now,
		CreateTime:      now,
	}

	err := dao.DB().Transaction(func(tx *gorm.DB) error {
		// 1. 创建帖子
		if err := tx.Create(post).Error; err != nil {
			return fmt.Errorf("创建帖子失败: %w", err)
		}

		// 2. 回填 threadId = 自身ID
		if err := tx.Model(post).UpdateColumn("thread_id", post.ID).Error; err != nil {
			return fmt.Errorf("更新ThreadId失败: %w", err)
		}
		post.ThreadId = post.ID

		return nil
	})

	return post, err
}

// Update 编辑帖子
func (s *postService) UpdateRootPost(req dto.UpdateRootPostRequest) error {
	nodeID := hashid.Slug2Id[model.Node](*req.NodeSlug)
	postID := hashid.Slug2Id[model.Post](req.Slug)
	node := s.nodeRepo.Get(nodeID)
	if node == nil || node.Status != model.StatusOk {
		return util.NewErrorMsg("节点不存在")
	}

	// 事务：Transaction + 闭包，全部使用 tx 操作
	err := dao.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Post{}).Where("id = ?", postID).Updates(map[string]interface{}{
			"node_id":    node.ID,
			"title":      req.Title,
			"content":    req.Content,
			"updated_at": util.NowTimestamp(),
		}).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

// CreateReply 创建回复
func (s *postService) CreateReply(userID int64, dto form.ReplyCreateForm) (*model.Post, error) {
	parentID := hashid.Slug2Id[model.Post](dto.ParentSlug)

	if parentID <= 0 {
		return nil, errors.New("无效的父级帖子")
	}

	now := util.NowTimestamp()
	post := &model.Post{
		Type:       model.PostTypeNormal,
		UserId:     userID,
		Title:		dto.Title,
		Content:    dto.Content,
		ImageList:  dto.ImageList,
		Status:     model.StatusOk,
		CreateTime: now,
	}

	err := dao.DB().Transaction(func(tx *gorm.DB) error {
		// 1. 查询并校验父级帖子
		var parentPost model.Post
		if err := tx.Where("id = ? AND status = ?", parentID, model.StatusOk).First(&parentPost).Error; err != nil {
			return fmt.Errorf("父级帖子不存在或已删除: %w", err)
		}

		// 2. 确定 threadId 和 nodeId（继承自父级）
		threadId := parentPost.ThreadId
		if threadId == 0 {
			threadId = parentPost.ID
		}
		post.ParentId = parentID
		post.ThreadId = threadId
		post.NodeId = parentPost.NodeId

		// 3. 创建回复
		if err := tx.Create(post).Error; err != nil {
			return fmt.Errorf("创建回复失败: %w", err)
		}

		// 4. 更新根帖最后评论时间
		if err := tx.Model(&model.Post{}).
			Where("id = ?", threadId).
			UpdateColumn("last_comment_time", now).Error; err != nil {
			return fmt.Errorf("更新根帖最后评论时间失败: %w", err)
		}

		return nil
	})

	return post, err
}

// Update 编辑帖子
func (s *postService) UpdateReply(dto form.ReplyUpdateForm) error {
	postID := hashid.Slug2Id[model.Post](dto.Slug)

	// 事务：Transaction + 闭包，全部使用 tx 操作
	err := dao.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Post{}).Where("id = ?", postID).Updates(map[string]interface{}{
			"title":      dto.Title,
			"content":    dto.Content,
			"updated_at": util.NowTimestamp(),
		}).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

// SetRecommend 设置推荐
func (s *postService) SetRecommend(postId int64, recommend bool) error {
	return dao.PostDao.UpdateColumn(postId, "recommend", recommend)
}

// GetPostTags 获取话题标签
func (s *postService) GetPostTags(postId int64) []model.Tag {
	return cache.TagCache.GetPostTags(postId)
}

// GetPostInIds 根据编号批量获取主题
func (s *postService) GetPostInIds(postIds []int64) map[int64]model.Post {
	if len(postIds) == 0 {
		return nil
	}
	var posts []model.Post
	// ✅ 补充错误处理
	if err := dao.DB().Where("id IN (?)", postIds).Find(&posts).Error; err != nil {
		log.Error("GetPostInIds failed: %v", err)
		return nil
	}

	postsMap := make(map[int64]model.Post, len(posts))
	for _, post := range posts {
		postsMap[post.ID] = post
	}
	return postsMap
}

// IncrViewCount 浏览数+1
func (s *postService) IncrViewCount(postId int64) {
	// ✅ gorm.Expr 在 v2 中用法不变，但需补充错误日志
	if err := dao.DB().Model(&model.Post{}).Where("id = ?", postId).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
		log.Error("IncrViewCount failed: %v", err)
	}
}

// OnComment 评论时更新最后回复时间和回复数量
func (s *postService) OnComment(postId, lastCommentUserId, lastCommentTime int64) {
	// ✅🔴 关键修复：原代码在 tx 闭包内使用了 dao.DB()，导致事务完全失效
	// 升级为全部使用 tx，确保两条 UPDATE 在同一事务中原子执行
	err := dao.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Post{}).Where("id = ?", postId).Updates(map[string]interface{}{
			"comment_count":        gorm.Expr("comment_count + ?", 1),
			"last_comment_user_id": lastCommentUserId,
			"last_comment_time":    lastCommentTime, // ✅ 修复字段名：lastCommentTime → last_comment_time（snake_case）
		}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.PostTag{}).Where("post_id = ?", postId).Updates(map[string]interface{}{
			"last_comment_time": lastCommentTime,
		}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error("OnComment failed: %v", err)
	}
}

// GenerateRss 生成 RSS / Atom
func (s *postService) GenerateRss() {
	posts := dao.PostDao.Find(querybuilder.NewQueryBuilder().
		Where("status = ?", model.StatusOk).Desc("id").Limit(1000))

	var items []*feeds.Item
	for _, post := range posts {
		postSlug := hashid.Id2Slug[model.Post](post.ID)
		postUrl := urls.PostUrl(postSlug)
		user := cache.UserCache.Get(post.UserId)
		if user == nil {
			continue
		}
		item := &feeds.Item{
			Title:       post.Title,
			Link:        &feeds.Link{Href: postUrl},
			Description: util.GetMarkdownSummary(post.Content),
			Author:      &feeds.Author{Name: user.Avatar, Email: user.Email.String},
			Created:     util.TimeFromTimestamp(post.CreateTime),
		}
		items = append(items, item)
	}

	siteTitle := cache.SettingCache.GetValue(model.SettingSiteTitle)
	siteDescription := cache.SettingCache.GetValue(model.SettingSiteDescription)
	feed := &feeds.Feed{
		Title:       siteTitle,
		Link:        &feeds.Link{Href: viper.GetString("base.baseUrl")},
		Description: siteDescription,
		Author:      &feeds.Author{Name: siteTitle},
		Created:     time.Now(),
		Items:       items,
	}

	staticPath := viper.GetString("base.static_path")

	atom, err := feed.ToAtom()
	if err != nil {
		log.Error("GenerateRss ToAtom failed: %v", err)
	} else {
		_ = util.WriteString(path.Join(staticPath, "post_atom.xml"), atom, false)
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.Error("GenerateRss ToRss failed: %v", err)
	} else {
		_ = util.WriteString(path.Join(staticPath, "post_rss.xml"), rss, false)
	}
}

// ScanDesc 倒序扫描
func (s *postService) ScanDesc(dateFrom, dateTo int64, cb ScanPostCallback) {
	var cursor int64 = math.MaxInt64
	for {
		list := dao.PostDao.Find(querybuilder.NewQueryBuilder().
			Lt("id", cursor).
			Gte("create_time", dateFrom).
			Lt("create_time", dateTo).
			Desc("id").Limit(1000))
		if len(list) == 0 { // ✅ 移除冗余的 nil 判断
			break
		}
		cursor = list[len(list)-1].ID
		cb(list)
	}
}