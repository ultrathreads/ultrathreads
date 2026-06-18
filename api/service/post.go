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
	"ultrathreads/domain"
	"ultrathreads/model"
	"ultrathreads/repository"
	"ultrathreads/util"
	"ultrathreads/util/hashid"
	"ultrathreads/util/log"
	"ultrathreads/util/querybuilder"
	"ultrathreads/util/urls"
)

type ScanPostCallback func(posts []domain.Post)

// PostService 帖子业务契约
type PostService interface {
	Get(id int64) *domain.Post
	GetBySlug(slug string) *domain.Post
	Find(cnd *querybuilder.QueryBuilder) []domain.Post
	List(cnd *querybuilder.QueryBuilder) ([]domain.Post, *querybuilder.Paging)
	Count(cnd *querybuilder.QueryBuilder) int64
	GetNodeThreadsFull(page, limit int, nodeSlug string) ([]domain.Post, *querybuilder.Paging)
	GetTagThreadsFull(tagSlug string, page int) ([]domain.Post, *querybuilder.Paging)
	GetPostTree(postSlug string) (*domain.Post, []domain.Post, error)
	GetPostsByThreadId(slug string) ([]domain.Post, error)
	GetUserPosts(userSlug, postType string, page int, pageSize int) ([]domain.Post, *querybuilder.Paging)
	Delete(id int64) error
	Undelete(id int64) error
	CreateRootPost(userID int64, cmd domain.CreatePostCommand) (*domain.Post, error)
	UpdateRootPost(cmd domain.UpdatePostCommand) error
	CreateReply(userID int64, cmd domain.CreateReplyCommand) (*domain.Post, error)
	UpdateReply(cmd domain.UpdateReplyCommand) error
	SetRecommend(postId int64, recommend bool) error
	GetPostTags(postId int64) []model.Tag
	GetPostInIds(postIds []int64) map[int64]domain.Post
	IncrViewCount(postId int64)
	OnComment(postId, lastCommentUserId, lastCommentTime int64)
	GenerateRss()
	ScanDesc(dateFrom, dateTo int64, cb ScanPostCallback)
}

func NewPostService(repo repository.PostRepository, nodeRepo repository.NodeRepository, postTagSvc PostTagService, settingSvc SettingService, tagCache cache.TagCacheInterface, userCache cache.UserCacheInterface, settingCache cache.SettingCacheInterface, db *gorm.DB) PostService {
	return &postService{
		repo:         repo,
		nodeRepo:     nodeRepo,
		postTagSvc:   postTagSvc,
		settingSvc:   settingSvc,
		tagCache:     tagCache,
		userCache:    userCache,
		settingCache: settingCache,
		db:           db,
	}
}

type postService struct {
	repo         repository.PostRepository
	nodeRepo     repository.NodeRepository
	postTagSvc   PostTagService
	settingSvc   SettingService
	tagCache     cache.TagCacheInterface
	userCache    cache.UserCacheInterface
	settingCache cache.SettingCacheInterface
	db           *gorm.DB
}

func (s *postService) Get(id int64) *domain.Post {
	return toDomainPost(s.repo.Get(id))
}

func (s *postService) GetBySlug(slug string) *domain.Post {
	id := hashid.Slug2Id[model.Post](slug)
	return toDomainPost(s.repo.Get(id))
}

func (s *postService) Find(cnd *querybuilder.QueryBuilder) []domain.Post {
	return toDomainPosts(s.repo.Find(cnd))
}

func (s *postService) List(cnd *querybuilder.QueryBuilder) ([]domain.Post, *querybuilder.Paging) {
	posts, paging := s.repo.List(cnd)
	return toDomainPosts(posts), paging
}

func (s *postService) Count(cnd *querybuilder.QueryBuilder) int64 {
	return s.repo.Count(cnd)
}

func (s *postService) GetNodeThreadsFull(page, limit int, nodeSlug string) ([]domain.Post, *querybuilder.Paging) {
	nodeId := hashid.Slug2Id[model.Node](nodeSlug)
	rootCnd := querybuilder.NewQueryBuilder().
		Eq("parent_id", 0).
		Eq("status", model.StatusOk)

	if nodeId > 0 {
		rootCnd = rootCnd.Eq("node_id", nodeId)
	}

	rootCnd = rootCnd.
		Desc("is_pinned").
		Desc("last_replied_at").
		Page(page, limit)

	rootPosts, paging := s.repo.List(rootCnd)
	if len(rootPosts) == 0 {
		return []domain.Post{}, paging
	}

	return s.expandThreadPosts(rootPosts), paging
}

func (s *postService) GetTagThreadsFull(tagSlug string, page int) ([]domain.Post, *querybuilder.Paging) {
	tagId := hashid.Slug2Id[model.Tag](tagSlug)

	subQuery := s.db.Model(&model.PostTag{}).
		Select("post_id").
		Where("tag_id = ? AND status = ?", tagId, model.StatusOk)

	rootCnd := querybuilder.NewQueryBuilder().
		Eq("parent_id", 0).
		Eq("status", model.StatusOk).
		Where("id IN (?)", subQuery).
		Desc("last_replied_at").
		Page(page, 20)

	rootPosts, paging := s.repo.List(rootCnd)
	if len(rootPosts) == 0 {
		return []domain.Post{}, paging
	}

	return s.expandThreadPosts(rootPosts), paging
}

// expandThreadPosts 根据根帖列表，查询并组装完整的帖子树（含回帖），保持根帖顺序
func (s *postService) expandThreadPosts(rootPosts []model.Post) []domain.Post {
	threadIds := make([]int64, 0, len(rootPosts))
	for _, p := range rootPosts {
		threadIds = append(threadIds, p.ID)
	}

	allPosts := s.repo.Find(querybuilder.NewQueryBuilder().
		In("thread_id", threadIds).
		Eq("status", model.StatusOk).
		Asc("created_at"))

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

	return toDomainPosts(result)
}

func (s *postService) GetPostTree(postSlug string) (*domain.Post, []domain.Post, error) {
	postId := hashid.Slug2Id[model.Post](postSlug)
	if postId <= 0 {
		return nil, nil, errors.New("invalid post_id")
	}

	currentPost := s.repo.Get(postId)
	if currentPost == nil || currentPost.Status != model.StatusOk {
		return nil, nil, errors.New("post not found")
	}

	var posts []model.Post
	if currentPost.ThreadId > 0 {
		cnd := querybuilder.NewQueryBuilder().
			Eq("thread_id", currentPost.ThreadId).
			Eq("status", model.StatusOk).
			Asc("created_at")
		posts = s.repo.Find(cnd)
	}

	if posts == nil {
		posts = []model.Post{}
	}

	return toDomainPost(currentPost), toDomainPosts(posts), nil
}

func (s *postService) GetPostsByThreadId(slug string) ([]domain.Post, error) {
	threadId := hashid.Slug2Id[model.Post](slug)
	if threadId <= 0 {
		return nil, errors.New("invalid thread_id")
	}

	cnd := querybuilder.NewQueryBuilder().
		Eq("thread_id", threadId).
		Eq("status", model.StatusOk).
		Asc("created_at")

	posts := s.repo.Find(cnd)
	if posts == nil {
		posts = []model.Post{}
	}

	return toDomainPosts(posts), nil
}

func (s *postService) GetUserPosts(userSlug, postType string, page int, pageSize int) ([]domain.Post, *querybuilder.Paging) {
	userID := hashid.Slug2Id[model.User](userSlug)
	qb := querybuilder.NewQueryBuilder().
		Eq("user_id", userID).
		Eq("status", model.StatusOk).
		Page(page, pageSize).
		Desc("id")

	switch postType {
	case "reply":
		qb.NotEq("parent_id", 0)
	case "root":
		fallthrough
	default:
		qb.Eq("parent_id", 0)
	}

	posts, paging := s.repo.List(qb)
	return toDomainPosts(posts), paging
}

func (s *postService) Delete(id int64) error {
	err := s.repo.UpdateColumn(id, "status", model.StatusDeleted)
	if err == nil {
		s.postTagSvc.DeleteByPostId(id)
	}
	return err
}

func (s *postService) Undelete(id int64) error {
	err := s.repo.UpdateColumn(id, "status", model.StatusOk)
	if err == nil {
		s.postTagSvc.UndeleteByPostId(id)
	}
	return err
}

func (s *postService) CreateRootPost(userID int64, cmd domain.CreatePostCommand) (*domain.Post, error) {
	nodeID := hashid.Slug2Id[model.Node](cmd.NodeSlug)

	if nodeID <= 0 {
		nodeID = s.settingSvc.GetSetting().DefaultNodeId
	}
	if nodeID <= 0 {
		return nil, errors.New("请配置默认节点")
	}
	node := s.nodeRepo.Get(nodeID)
	if node == nil || node.Status != model.StatusOk {
		return nil, errors.New("节点不存在或已禁用")
	}

	now := time.Now()
	post := &model.Post{
		Type:          model.PostTypeNormal,
		UserId:        userID,
		NodeId:        nodeID,
		Title:         cmd.Title,
		Content:       cmd.Content,
		Status:        model.StatusOk,
		LastRepliedAt: now,
		CreatedAt:     now,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(post).Error; err != nil {
			return fmt.Errorf("创建帖子失败: %w", err)
		}

		if err := tx.Model(post).UpdateColumn("thread_id", post.ID).Error; err != nil {
			return fmt.Errorf("更新ThreadId失败: %w", err)
		}
		post.ThreadId = post.ID

		return nil
	})

	if err != nil {
		return nil, err
	}

	return toDomainPost(post), nil
}

func (s *postService) UpdateRootPost(cmd domain.UpdatePostCommand) error {
	nodeID := hashid.Slug2Id[model.Node](*cmd.NodeSlug)
	postID := hashid.Slug2Id[model.Post](cmd.Slug)
	node := s.nodeRepo.Get(nodeID)
	if node == nil || node.Status != model.StatusOk {
		return util.NewErrorMsg("节点不存在")
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Post{}).Where("id = ?", postID).Updates(map[string]interface{}{
			"node_id":    node.ID,
			"title":      cmd.Title,
			"content":    cmd.Content,
			"updated_at": util.NowTimestamp(),
		}).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *postService) CreateReply(userID int64, cmd domain.CreateReplyCommand) (*domain.Post, error) {
	parentID := hashid.Slug2Id[model.Post](cmd.ParentSlug)

	if parentID <= 0 {
		return nil, errors.New("无效的父级帖子")
	}

	title := util.ExtractReplyTitle(cmd.Content, 20)
	now := time.Now()
	post := &model.Post{
		Type:      model.PostTypeNormal,
		UserId:    userID,
		Title:     title,
		Content:   cmd.Content,
		Status:    model.StatusOk,
		CreatedAt: now,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var parentPost model.Post
		if err := tx.Where("id = ? AND status = ?", parentID, model.StatusOk).First(&parentPost).Error; err != nil {
			return fmt.Errorf("父级帖子不存在或已删除: %w", err)
		}

		threadId := parentPost.ThreadId
		if threadId == 0 {
			threadId = parentPost.ID
		}
		post.ParentId = parentID
		post.ThreadId = threadId
		post.NodeId = parentPost.NodeId

		if err := tx.Create(post).Error; err != nil {
			return fmt.Errorf("创建回复失败: %w", err)
		}

		if err := tx.Model(&model.Post{}).
			Where("id = ?", threadId).
			UpdateColumn("last_replied_at", now).Error; err != nil {
			return fmt.Errorf("更新根帖最后评论时间失败: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return toDomainPost(post), nil
}

func (s *postService) UpdateReply(cmd domain.UpdateReplyCommand) error {
	postID := hashid.Slug2Id[model.Post](cmd.Slug)

	title := util.ExtractReplyTitle(*cmd.Content, 20)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Post{}).Where("id = ?", postID).Updates(map[string]interface{}{
			"title":      title,
			"content":    cmd.Content,
			"updated_at": util.NowTimestamp(),
		}).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *postService) SetRecommend(postId int64, recommend bool) error {
	return s.repo.UpdateColumn(postId, "recommend", recommend)
}

func (s *postService) GetPostTags(postId int64) []model.Tag {
	return s.tagCache.GetPostTags(postId)
}

func (s *postService) GetPostInIds(postIds []int64) map[int64]domain.Post {
	if len(postIds) == 0 {
		return nil
	}
	var posts []model.Post
	if err := s.db.Where("id IN (?)", postIds).Find(&posts).Error; err != nil {
		log.Error("GetPostInIds failed: %v", err)
		return nil
	}

	postsMap := make(map[int64]domain.Post, len(posts))
	for _, post := range posts {
		postsMap[post.ID] = *toDomainPost(&post)
	}
	return postsMap
}

// toDomainPost 将 model.Post 转换为 domain.Post
func toDomainPost(m *model.Post) *domain.Post {
	if m == nil {
		return nil
	}
	return &domain.Post{
		ID:              m.ID,
		ThreadId:        m.ThreadId,
		ParentId:        m.ParentId,
		Type:            m.Type,
		NodeId:          m.NodeId,
		UserId:          m.UserId,
		Title:           m.Title,
		Content:         m.Content,
		ImageList:       m.ImageList,
		IsPinned:        m.IsPinned,
		Recommend:       m.Recommend,
		ViewCount:       m.ViewCount,
		LikeCount:       m.LikeCount,
		Status:          m.Status,
		LastReplyUserId: m.LastReplyUserId,
		LastRepliedAt:   m.LastRepliedAt,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		ExtraData:       m.ExtraData,
	}
}

// toDomainPosts 将 []model.Post 转换为 []domain.Post
func toDomainPosts(models []model.Post) []domain.Post {
	if len(models) == 0 {
		return []domain.Post{}
	}
	result := make([]domain.Post, len(models))
	for i, m := range models {
		d := toDomainPost(&m)
		if d != nil {
			result[i] = *d
		}
	}
	return result
}

func (s *postService) IncrViewCount(postId int64) {
	if err := s.db.Model(&model.Post{}).Where("id = ?", postId).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error; err != nil {
		log.Error("IncrViewCount failed: %v", err)
	}
}

func (s *postService) OnComment(postId, lastCommentUserId, lastCommentTime int64) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Post{}).Where("id = ?", postId).Updates(map[string]interface{}{
			"comment_count":      gorm.Expr("comment_count + ?", 1),
			"last_reply_user_id": lastCommentUserId,
			"last_replied_at":    lastCommentTime,
		}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.PostTag{}).Where("post_id = ?", postId).Updates(map[string]interface{}{
			"last_replied_at": lastCommentTime,
		}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error("OnComment failed: %v", err)
	}
}

func (s *postService) GenerateRss() {
	posts := s.repo.Find(querybuilder.NewQueryBuilder().
		Where("status = ?", model.StatusOk).Desc("id").Limit(1000))

	var items []*feeds.Item
	for _, post := range posts {
		postSlug := hashid.Id2Slug[model.Post](post.ID)
		postUrl := urls.PostUrl(postSlug)
		user := s.userCache.Get(post.UserId)
		if user == nil {
			continue
		}
		item := &feeds.Item{
			Title:       post.Title,
			Link:        &feeds.Link{Href: postUrl},
			Description: util.GetMarkdownSummary(post.Content),
			Author:      &feeds.Author{Name: user.Avatar, Email: user.Email.String},
			Created:     post.CreatedAt,
		}
		items = append(items, item)
	}

	siteTitle := s.settingCache.GetValue(model.SettingSiteTitle)
	siteDescription := s.settingCache.GetValue(model.SettingSiteDescription)
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

func (s *postService) ScanDesc(dateFrom, dateTo int64, cb ScanPostCallback) {
	var cursor int64 = math.MaxInt64
	for {
		list := s.repo.Find(querybuilder.NewQueryBuilder().
			Lt("id", cursor).
			Gte("created_at", dateFrom).
			Lt("created_at", dateTo).
			Desc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		cb(toDomainPosts(list))
	}
}
