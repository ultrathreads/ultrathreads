package service

import (
	"math"
	"path"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gorilla/feeds"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"ultrathreads/cache"
	"ultrathreads/dao"
	"ultrathreads/form"
	"ultrathreads/model"
	"ultrathreads/util"
	"ultrathreads/util/log"
	"ultrathreads/util/querybuilder"
	"ultrathreads/util/urls"
)

type ScanArticleCallback func(articles []model.Article)

var ArticleService = newArticleService()

func newArticleService() *articleService {
	return &articleService{}
}

type articleService struct{}

func (s *articleService) Get(id int64) *model.Article {
	return dao.ArticleDao.Get(id)
}

func (s *articleService) Find(cnd *querybuilder.QueryBuilder) []model.Article {
	return dao.ArticleDao.Find(cnd)
}

func (s *articleService) List(cnd *querybuilder.QueryBuilder) (list []model.Article, paging *querybuilder.Paging) {
	return dao.ArticleDao.List(cnd)
}

// Create 发表文章
func (s *articleService) Create(dto form.ArticleCreateForm) (*model.Article, error) {
	article := &model.Article{
		UserId:      dto.UserID,
		Title:       dto.Title,
		Summary:     dto.Summary,
		Content:     dto.Content,
		ContentType: model.ContentTypeMarkdown,
		Status:      model.StatusOk,
		Share:       false,
		SourceUrl:   "",
		CreateTime:  util.NowTimestamp(),
		UpdateTime:  util.NowTimestamp(),
	}

	// ✅ v2 事务：Transaction + 闭包，tx 替代原来的 dao.DB()
	err := dao.DB().Transaction(func(tx *gorm.DB) error {
		tagIDs := dao.TagDao.GetOrCreates(util.ParseTagsToArray(dto.Tags))

		// ⚠️ 注意：如果 ArticleDao.Create 内部仍使用全局 db，
		// 则此处的 tx 不会生效。需要改造 DAO 支持传入 tx，
		// 或在此处直接使用 tx.Create(article)
		if err := tx.Create(article).Error; err != nil {
			return err
		}

		dao.ArticleTagDao.AddArticleTags(article.ID, tagIDs)
		return nil
	})

	return article, err
}

// Update 编辑文章
func (s *articleService) Update(dto form.ArticleUpdateForm) error {
	err := dao.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Article{}).Where("id = ?", dto.ID).Updates(map[string]interface{}{
			"title":       dto.Title,
			"content":     dto.Content,
			"update_time": util.NowTimestamp(),
		}).Error; err != nil {
			return err
		}

		tagIds := dao.TagDao.GetOrCreates(util.ParseTagsToArray(dto.Tags))
		dao.ArticleTagDao.DeleteArticleTags(dto.ID)
		dao.ArticleTagDao.AddArticleTags(dto.ID, tagIds)
		return nil
	})

	cache.ArticleTagCache.Invalidate(dto.ID)
	return err
}

func (s *articleService) Delete(id int64) error {
	err := dao.ArticleDao.UpdateColumn(id, "status", model.StatusDeleted)
	if err == nil {
		ArticleTagService.DeleteByArticleId(id)
	}
	return err
}

// GetArticleInIds 根据文章编号批量获取文章
func (s *articleService) GetArticleInIds(articleIds []int64) []model.Article {
	if len(articleIds) == 0 {
		return nil
	}
	var articles []model.Article
	// ✅ v2 必须处理错误；Find 即使无结果也不返回 ErrRecordNotFound
	if err := dao.DB().Where("id IN (?)", articleIds).Find(&articles).Error; err != nil {
		log.Error("GetArticleInIds failed: %v", err)
		return nil
	}
	return articles
}

// GetArticles 游标分页文章列表
func (s *articleService) GetArticles(cursor int64) (articles []model.Article, nextCursor int64) {
	cnd := querybuilder.NewQueryBuilder().Eq("status", model.StatusOk).Desc("id").Limit(20)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	articles = dao.ArticleDao.Find(cnd)
	if len(articles) > 0 {
		nextCursor = articles[len(articles)-1].ID
	} else {
		nextCursor = cursor
	}
	return
}

// GetArticleTags 获取文章对应的标签
func (s *articleService) GetArticleTags(articleId int64) []model.Tag {
	articleTags := dao.ArticleTagDao.Find(querybuilder.NewQueryBuilder().Where("article_id = ?", articleId))
	var tagIds []int64
	for _, articleTag := range articleTags {
		tagIds = append(tagIds, articleTag.TagId)
	}
	return cache.TagCache.GetList(tagIds)
}

// GetTagArticles 标签文章列表
func (s *articleService) GetTagArticles(tagId int64, cursor int64) (articles []model.Article, nextCursor int64) {
	cnd := querybuilder.NewQueryBuilder().Eq("tag_id", tagId).Eq("status", model.StatusOk).Desc("id").Limit(20)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	nextCursor = cursor
	articleTags := dao.ArticleTagDao.Find(cnd)
	if len(articleTags) > 0 {
		var articleIds []int64
		for _, articleTag := range articleTags {
			articleIds = append(articleIds, articleTag.ArticleId)
			nextCursor = articleTag.ID
		}
		articles = s.GetArticleInIds(articleIds)
	}
	return
}

// GetRelatedArticles 相关文章
func (s *articleService) GetRelatedArticles(articleId int64) []model.Article {
	tagIds := cache.ArticleTagCache.Get(articleId)
	if len(tagIds) == 0 {
		return nil
	}
	var articleTags []model.ArticleTag
	// ✅ 补充错误处理
	if err := dao.DB().Where("tag_id IN (?)", tagIds).Limit(30).Find(&articleTags).Error; err != nil {
		log.Error("GetRelatedArticles find tags failed: %v", err)
		return nil
	}

	set := hashset.New()
	for _, articleTag := range articleTags {
		set.Add(articleTag.ArticleId)
	}

	var articleIds []int64
	for i, aid := range set.Values() {
		if i >= 10 {
			break
		}
		articleIds = append(articleIds, aid.(int64))
	}

	return s.GetArticleInIds(articleIds)
}

// GetUserNewestArticles 用户最新文章
func (s *articleService) GetUserNewestArticles(userId int64) []model.Article {
	return dao.ArticleDao.Find(querybuilder.NewQueryBuilder().
		Where("user_id = ? AND status = ?", userId, model.StatusOk).
		Desc("id").Limit(10))
}

// ScanDesc 倒序扫描文章
func (s *articleService) ScanDesc(dateFrom, dateTo int64, cb ScanArticleCallback) {
	var cursor int64 = math.MaxInt64
	for {
		list := dao.ArticleDao.Find(querybuilder.NewQueryBuilder("id", "status", "create_time", "update_time").
			Lt("id", cursor).Gte("create_time", dateFrom).Lt("create_time", dateTo).Desc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		cb(list)
	}
}

// GenerateRss 生成 RSS / Atom 订阅文件
func (s *articleService) GenerateRss() {
	articles := dao.ArticleDao.Find(querybuilder.NewQueryBuilder().
		Where("status = ?", model.StatusOk).Desc("id").Limit(1000))

	var items []*feeds.Item
	for _, article := range articles {
		articleUrl := urls.ArticleUrl(article.ID)
		user := cache.UserCache.Get(article.UserId)
		if user == nil {
			continue
		}
		description := ""
		if article.ContentType == model.ContentTypeMarkdown {
			description = util.GetMarkdownSummary(article.Content)
		} else {
			description = util.GetHtmlSummary(article.Content)
		}
		item := &feeds.Item{
			Title:       article.Title,
			Link:        &feeds.Link{Href: articleUrl},
			Description: description,
			Author:      &feeds.Author{Name: user.Avatar, Email: user.Email.String},
			Created:     util.TimeFromTimestamp(article.CreateTime),
		}
		items = append(items, item)
	}

	siteTitle := cache.SettingCache.GetValue(model.SettingSiteTitle)
	siteDescription := cache.SettingCache.GetValue(model.SettingSiteDescription)
	feed := &feeds.Feed{
		Title:       siteTitle,
		Link:        &feeds.Link{Href: viper.GetString("base.url")},
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
		_ = util.WriteString(path.Join(staticPath, "atom.xml"), atom, false)
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.Error("GenerateRss ToRss failed: %v", err)
	} else {
		_ = util.WriteString(path.Join(staticPath, "rss.xml"), rss, false)
	}
}