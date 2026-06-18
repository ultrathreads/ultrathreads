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
	"ultrathreads/dto"
	"ultrathreads/model"
	"ultrathreads/repository"
	"ultrathreads/util"
	"ultrathreads/util/log"
	"ultrathreads/util/querybuilder"
	"ultrathreads/util/urls"
)

type ScanArticleCallback func(articles []model.Article)

// ArticleService 文章业务契约
type ArticleService interface {
	Get(id int64) *model.Article
	Find(cnd *querybuilder.QueryBuilder) []model.Article
	List(cnd *querybuilder.QueryBuilder) ([]model.Article, *querybuilder.Paging)
	Create(req dto.ArticleCreateForm) (*model.Article, error)
	Update(req dto.ArticleUpdateForm) error
	Delete(id int64) error
	GetArticleInIds(articleIds []int64) []model.Article
	GetArticles(cursor int64) ([]model.Article, int64)
	GetArticleTags(articleId int64) []model.Tag
	GetTagArticles(tagId int64, cursor int64) ([]model.Article, int64)
	GetRelatedArticles(articleId int64) []model.Article
	GetUserNewestArticles(userId int64) []model.Article
	ScanDesc(dateFrom, dateTo int64, cb ScanArticleCallback)
	GenerateRss()
}

func NewArticleService(repo repository.ArticleRepository, tagRepo repository.TagRepository, articleTagRepo repository.ArticleTagRepository, articleTagSvc ArticleTagService, articleTagCache cache.ArticleTagCacheInterface, tagCache cache.TagCacheInterface, userCache cache.UserCacheInterface, settingCache cache.SettingCacheInterface, db *gorm.DB) ArticleService {
	return &articleService{
		repo:            repo,
		tagRepo:         tagRepo,
		articleTagRepo:  articleTagRepo,
		articleTagSvc:   articleTagSvc,
		articleTagCache: articleTagCache,
		tagCache:        tagCache,
		userCache:       userCache,
		settingCache:    settingCache,
		db:              db,
	}
}

type articleService struct {
	repo            repository.ArticleRepository
	tagRepo         repository.TagRepository
	articleTagRepo  repository.ArticleTagRepository
	articleTagSvc   ArticleTagService
	articleTagCache cache.ArticleTagCacheInterface
	tagCache        cache.TagCacheInterface
	userCache       cache.UserCacheInterface
	settingCache    cache.SettingCacheInterface
	db              *gorm.DB
}

func (s *articleService) Get(id int64) *model.Article {
	return s.repo.Get(id)
}

func (s *articleService) Find(cnd *querybuilder.QueryBuilder) []model.Article {
	return s.repo.Find(cnd)
}

func (s *articleService) List(cnd *querybuilder.QueryBuilder) ([]model.Article, *querybuilder.Paging) {
	return s.repo.List(cnd)
}

func (s *articleService) Create(req dto.ArticleCreateForm) (*model.Article, error) {
	article := &model.Article{
		UserId:      req.UserID,
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		ContentType: model.ContentTypeMarkdown,
		Status:      model.StatusOk,
		Share:       false,
		SourceUrl:   "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		tagIDs := s.tagRepo.GetOrCreates(util.ParseTagsToArray(req.Tags))

		if err := tx.Create(article).Error; err != nil {
			return err
		}

		s.articleTagRepo.AddArticleTags(article.ID, tagIDs)
		return nil
	})

	return article, err
}

func (s *articleService) Update(req dto.ArticleUpdateForm) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Article{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
			"title":      req.Title,
			"content":    req.Content,
			"updated_at": time.Now(),
		}).Error; err != nil {
			return err
		}

		tagIds := s.tagRepo.GetOrCreates(util.ParseTagsToArray(req.Tags))
		s.articleTagRepo.DeleteArticleTags(req.ID)
		s.articleTagRepo.AddArticleTags(req.ID, tagIds)
		return nil
	})

	s.articleTagCache.Invalidate(req.ID)
	return err
}

func (s *articleService) Delete(id int64) error {
	err := s.repo.UpdateColumn(id, "status", model.StatusDeleted)
	if err == nil {
		s.articleTagSvc.DeleteByArticleId(id)
	}
	return err
}

func (s *articleService) GetArticleInIds(articleIds []int64) []model.Article {
	if len(articleIds) == 0 {
		return nil
	}
	var articles []model.Article
	if err := s.db.Where("id IN (?)", articleIds).Find(&articles).Error; err != nil {
		log.Error("GetArticleInIds failed: %v", err)
		return nil
	}
	return articles
}

func (s *articleService) GetArticles(cursor int64) (articles []model.Article, nextCursor int64) {
	cnd := querybuilder.NewQueryBuilder().Eq("status", model.StatusOk).Desc("id").Limit(20)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	articles = s.repo.Find(cnd)
	if len(articles) > 0 {
		nextCursor = articles[len(articles)-1].ID
	} else {
		nextCursor = cursor
	}
	return
}

func (s *articleService) GetArticleTags(articleId int64) []model.Tag {
	articleTags := s.articleTagRepo.Find(querybuilder.NewQueryBuilder().Where("article_id = ?", articleId))
	var tagIds []int64
	for _, articleTag := range articleTags {
		tagIds = append(tagIds, articleTag.TagId)
	}
	return s.tagCache.GetList(tagIds)
}

func (s *articleService) GetTagArticles(tagId int64, cursor int64) (articles []model.Article, nextCursor int64) {
	cnd := querybuilder.NewQueryBuilder().Eq("tag_id", tagId).Eq("status", model.StatusOk).Desc("id").Limit(20)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	nextCursor = cursor
	articleTags := s.articleTagRepo.Find(cnd)
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

func (s *articleService) GetRelatedArticles(articleId int64) []model.Article {
	tagIds := s.articleTagCache.Get(articleId)
	if len(tagIds) == 0 {
		return nil
	}
	var articleTags []model.ArticleTag
	if err := s.db.Where("tag_id IN (?)", tagIds).Limit(30).Find(&articleTags).Error; err != nil {
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

func (s *articleService) GetUserNewestArticles(userId int64) []model.Article {
	return s.repo.Find(querybuilder.NewQueryBuilder().
		Where("user_id = ? AND status = ?", userId, model.StatusOk).
		Desc("id").Limit(10))
}

func (s *articleService) ScanDesc(dateFrom, dateTo int64, cb ScanArticleCallback) {
	var cursor int64 = math.MaxInt64
	for {
		list := s.repo.Find(querybuilder.NewQueryBuilder("id", "status", "created_at", "updated_at").
			Lt("id", cursor).Gte("created_at", dateFrom).Lt("created_at", dateTo).Desc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		cb(list)
	}
}

func (s *articleService) GenerateRss() {
	articles := s.repo.Find(querybuilder.NewQueryBuilder().
		Where("status = ?", model.StatusOk).Desc("id").Limit(1000))

	var items []*feeds.Item
	for _, article := range articles {
		articleUrl := urls.ArticleUrl(article.ID)
		user := s.userCache.Get(article.UserId)
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
			Created:     article.CreatedAt,
		}
		items = append(items, item)
	}

	siteTitle := s.settingCache.GetValue(model.SettingSiteTitle)
	siteDescription := s.settingCache.GetValue(model.SettingSiteDescription)
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
