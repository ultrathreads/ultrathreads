package service

import (
	"errors"
	"math"
	"path"
	"time"

	"github.com/gorilla/feeds"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	"ultrathreads/cache"
	"ultrathreads/dao"
	"ultrathreads/form"
	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/log"
	"ultrathreads/util/querybuilder"
	"ultrathreads/util/urls"
)

type ScanPostCallback func(posts []model.Post)

var PostService = newPostService()

func newPostService() *postService {
	return &postService{}
}

type postService struct{}

func (s *postService) Get(id int64) *model.Post {
	return dao.PostDao.Get(id)
}

func (s *postService) Find(cnd *querybuilder.QueryBuilder) []model.Post {
	return dao.PostDao.Find(cnd)
}

func (s *postService) List(cnd *querybuilder.QueryBuilder) (list []model.Post, paging *querybuilder.Paging) {
	return dao.PostDao.List(cnd)
}

func (s *postService) Count(cnd *querybuilder.QueryBuilder) int {
	return dao.PostDao.Count(cnd)
}

// ListThreadsWithReplies 获取主帖列表（分页）+ 每个主帖下的所有回复（扁平化）
// 返回的列表已按 create_time ASC 排序，前端可直接根据 parent_id 组装树
func (s *postService) ListThreadsWithReplies(page, limit, nodeId int) ([]model.Post, *querybuilder.Paging) {
	// ========== 第一步：分页查出主帖ID ==========
	rootCnd := querybuilder.NewQueryBuilder().
    	Eq("parent_id", 0).
    	Eq("status", model.StatusOk)

	if nodeId > 0 {
		rootCnd = rootCnd.Eq("node_id", nodeId)
	}

	rootCnd = rootCnd.
		Desc("last_comment_time").
		Page(page, limit)

	rootPosts, paging := dao.PostDao.List(rootCnd)
	if len(rootPosts) == 0 {
		return []model.Post{}, paging
	}

	// 提取主帖ID列表
	threadIds := make([]int64, 0, len(rootPosts))
	for _, p := range rootPosts {
		threadIds = append(threadIds, p.ID)
	}

	// ========== 第二步：批量拉取这些主题下的所有帖子 ==========
	allCnd := querybuilder.NewQueryBuilder().
		In("thread_id", threadIds).
		Eq("status", model.StatusOk).
		Asc("create_time")

	allPosts := dao.PostDao.Find(allCnd)

	// ========== 第三步：按主帖原始分页顺序分组，再组内按时间排序 ==========
	// 保持主帖的分页排序（last_comment_time DESC），同时每个主题内部按时间正序
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

// GetPostWithThread 获取帖子详情及其所属主题的所有扁平回帖
// replies 已按 create_time ASC 排序，前端可直接根据 parent_id 组装树
func (s *postService) GetPostWithThread(postId int64) (*model.Post, []model.Post, error) {
	if postId <= 0 {
		return nil, nil, errors.New("invalid post_id")
	}

	// 1. 获取帖子详情
	post := dao.PostDao.Get(postId)
	if post == nil || post.Status != model.StatusOk {
		return nil, nil, errors.New("post not found")
	}

	// 2. 获取同主题下所有扁平回帖（含主帖自身）
	var replies []model.Post
	if post.ThreadId > 0 {
		cnd := querybuilder.NewQueryBuilder().
			Eq("thread_id", post.ThreadId).
			Eq("status", model.StatusOk).
			Asc("create_time")

		replies = dao.PostDao.Find(cnd)
	}

	// 保证返回非 nil 切片，避免前端 JSON 序列化时出现 null
	if replies == nil {
		replies = []model.Post{}
	}

	return post, replies, nil
}

// 删除
func (s *postService) Delete(id int64) error {
	err := dao.PostDao.UpdateColumn(id, "status", model.StatusDeleted)
	if err == nil {
		// 删掉标签文章
		PostTagService.DeleteByPostId(id)
	}
	return err
}

func (s *postService) Update(dto form.PostUpdateForm) error {
	node := dao.NodeDao.Get(dto.NodeID)
	if node == nil || node.Status != model.StatusOk {
		return util.NewErrorMsg("节点不存在")
	}
	err := dao.Tx(dao.DB(), func(tx *gorm.DB) error {
		err := dao.PostDao.Updates(dto.ID, map[string]interface{}{
			"node_id":     dto.NodeID,
			"title":       dto.Title,
			"content":     dto.Content,
			"update_time": util.NowTimestamp(),
		})
		if err != nil {
			return err
		}
		tagIds := dao.TagDao.GetOrCreates(util.ParseTagsToArray(dto.Tags)) // 创建文章对应标签
		dao.PostTagDao.DeletePostTags(dto.ID)                            // 先删掉所有的标签
		dao.PostTagDao.AddPostTags(dto.ID, tagIds)                       // 然后重新添加标签
		return nil
	})

	return err
}

// 取消删除
func (s *postService) Undelete(id int64) error {
	err := dao.PostDao.UpdateColumn(id, "status", model.StatusOk)
	if err == nil {
		// 删掉标签文章
		PostTagService.UndeleteByPostId(id)
	}
	return err
}

// 发表话题
func (s *postService) Create(dto form.PostCreateForm) (*model.Post, error) {
	nodeID := dto.NodeID
	if nodeID <= 0 {
		nodeID = SettingService.GetSetting().DefaultNodeId
		if nodeID <= 0 {
			return nil, errors.New("请配置默认节点")
		}
	}
	node := dao.NodeDao.Get(nodeID)
	if node == nil || node.Status != model.StatusOk {
		return nil, errors.New("节点不存在")
	}

	now := util.NowTimestamp()
	post := &model.Post{
		Type:            model.PostTypeNormal,
		UserId:          dto.UserID,
		NodeId:          nodeID,
		Title:           dto.Title,
		Content:         dto.Content,
		ImageList:       dto.ImageList,
		Status:          model.StatusOk,
		LastCommentTime: now,
		CreateTime:      now,
	}

	err := dao.Tx(dao.DB(), func(tx *gorm.DB) error {
		tagIds := dao.TagDao.GetOrCreates(util.ParseTagsToArray(dto.Tags))
		err := dao.PostDao.Create(post)
		if err != nil {
			return err
		}

		dao.PostTagDao.AddPostTags(post.ID, tagIds)
		return nil
	})
	if err == nil {
		// 节点话题计数
		NodeService.IncrTopicCount(nodeID)
		// 用户话题计数
		UserService.IncrTopicCount(dto.UserID)
		// 获得积分
		UserScoreService.IncrementPostPostScore(post)
	}
	return post, err
}

// 推荐
func (s *postService) SetRecommend(postId int64, recommend bool) error {
	return dao.PostDao.UpdateColumn(postId, "recommend", recommend)
}

// 话题的标签
func (s *postService) GetPostTags(postId int64) []model.Tag {
	postTags := dao.PostTagDao.Find(querybuilder.NewQueryBuilder().Where("post_id = ?", postId))

	var tagIds []int64
	for _, postTag := range postTags {
		tagIds = append(tagIds, postTag.TagId)
	}
	return cache.TagCache.GetList(tagIds)
}

// 指定标签下话题列表
func (s *postService) GetTagPosts(tagId int64, page int) (posts []model.Post, paging *querybuilder.Paging) {
	postTags, paging := dao.PostTagDao.List(querybuilder.NewQueryBuilder().
		Eq("tag_id", tagId).
		Eq("status", model.StatusOk).
		Page(page, 20).Desc("last_comment_time"))
	if len(postTags) > 0 {
		var postIds []int64
		for _, postTag := range postTags {
			postIds = append(postIds, postTag.PostId)
		}

		postsMap := s.GetPostInIds(postIds)
		if postsMap != nil {
			for _, postTag := range postTags {
				if post, found := postsMap[postTag.PostId]; found {
					posts = append(posts, post)
				}
			}
		}
	}
	return
}

// GetPostInIds 根据编号批量获取主题
func (s *postService) GetPostInIds(postIds []int64) map[int64]model.Post {
	if len(postIds) == 0 {
		return nil
	}
	var posts []model.Post
	dao.DB().Where("id in (?)", postIds).Find(&posts)

	postsMap := make(map[int64]model.Post, len(posts))
	for _, post := range posts {
		postsMap[post.ID] = post
	}
	return postsMap
}

// 浏览数+1
func (s *postService) IncrViewCount(postId int64) {
	dao.DB().Model(&model.Post{}).Where("id = ?", postId).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
}

// 当帖子被评论的时候，更新最后回复时间、回复数量+1
func (s *postService) OnComment(postId, lastCommentUserId, lastCommentTime int64) {
	dao.Tx(dao.DB(), func(tx *gorm.DB) error {
		if err := dao.DB().Model(&model.Post{}).Where("id = ?", postId).Updates(map[string]interface{}{"comment_count": gorm.Expr("comment_count + ?", 1), "last_comment_user_id": lastCommentUserId, "lastCommentTime": lastCommentTime}).Error; err != nil {
			return err
		}
		if err := dao.DB().Model(&model.PostTag{}).Where("post_id = ?", postId).Updates(map[string]interface{}{"last_comment_time": lastCommentTime}).Error; err != nil {
			return err
		}
		return nil
	})
}

// rss
func (s *postService) GenerateRss() {
	posts := dao.PostDao.Find(querybuilder.NewQueryBuilder().Where("status = ?", model.StatusOk).Desc("id").Limit(1000))

	var items []*feeds.Item
	for _, post := range posts {
		postUrl := urls.PostUrl(post.ID)
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
	atom, err := feed.ToAtom()
	if err != nil {
		log.Error(err.Error())
	} else {
		_ = util.WriteString(path.Join(viper.GetString("base.static_path"), "post_atom.xml"), atom, false)
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.Error(err.Error())
	} else {
		_ = util.WriteString(path.Join(viper.GetString("base.static_path"), "post_rss.xml"), rss, false)
	}
}

// 倒序扫描
func (s *postService) ScanDesc(dateFrom, dateTo int64, cb ScanPostCallback) {
	var cursor int64 = math.MaxInt64
	for {
		list := dao.PostDao.Find(querybuilder.NewQueryBuilder().Lt("id", cursor).
			Gte("create_time", dateFrom).Lt("create_time", dateTo).Desc("id").Limit(1000))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		cb(list)
	}
}
