package app

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"ultrathreads/cache"
	"ultrathreads/dto"
	"ultrathreads/handler/base"
	"ultrathreads/model"
	"ultrathreads/render"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/hashid"
	"ultrathreads/util/querybuilder"
)

type ArticleHandler struct {
	base.BaseHandler
	articleSvc      service.ArticleServicer
	favoriteSvc     service.FavoriteServicer
	articleTagCache cache.ArticleTagCacheInterface
	tagCache        cache.TagCacheInterface
}

func NewArticleHandler(articleSvc service.ArticleServicer, favoriteSvc service.FavoriteServicer, articleTagCache cache.ArticleTagCacheInterface, tagCache cache.TagCacheInterface) *ArticleHandler {
	return &ArticleHandler{
		articleSvc:      articleSvc,
		favoriteSvc:     favoriteSvc,
		articleTagCache: articleTagCache,
		tagCache:        tagCache,
	}
}

// Show show article by id
func (h *ArticleHandler) Show(ctx *gin.Context) {
	var gDto dto.IdRequest
	if h.BindAndValidate(ctx, &gDto) {
		article := h.articleSvc.Get(gDto.ID)
		if article == nil || article.Status != model.StatusOk {
			h.Fail(ctx, util.ErrorArticleNotFound)
			return
		}
		h.Success(ctx, render.ToArticle(article, h.articleTagCache, h.tagCache))
	}
}

// Store 发表文章
func (h *ArticleHandler) Store(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	if user == nil {
		h.Fail(ctx, util.ErrorNotLogin)
		return
	}
	var articleForm dto.ArticleCreateForm
	if h.BindAndValidate(ctx, &articleForm) {
		articleForm.UserID = user.ID
		article, err := h.articleSvc.Create(articleForm)
		if err != nil {
			h.Fail(ctx, util.FromError(err))
			return
		}
		h.Success(ctx, render.ToArticle(article, h.articleTagCache, h.tagCache))
	}
}

// Edit 编辑时获取详情
func (h *ArticleHandler) Edit(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	if user == nil {
		h.Fail(ctx, util.ErrorNotLogin)
		return
	}

	var gDto dto.IdRequest
	if h.BindAndValidate(ctx, &gDto) {
		article := h.articleSvc.Get(gDto.ID)

		if article == nil || article.Status != model.StatusOk {
			h.Fail(ctx, util.NewErrorMsg("话题不存在或已被删除"))
			return
		}
		if article.UserId != user.ID {
			h.Fail(ctx, util.NewErrorMsg("无权限"))
			return
		}

		tags := h.articleSvc.GetArticleTags(article.ID)
		var tagNames []string
		if len(tags) > 0 {
			for _, tag := range tags {
				tagNames = append(tagNames, tag.Name)
			}
		}

		h.Success(ctx, gin.H{
			"articleId": article.ID,
			"title":     article.Title,
			"content":   article.Content,
			"tags":      tagNames,
		})
	}
}

// Update 编辑文章
func (h *ArticleHandler) Update(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	if user == nil {
		h.Fail(ctx, util.ErrorNotLogin)
		return
	}
	var gDto dto.IdRequest
	if !h.BindAndValidate(ctx, &gDto) {
		h.Fail(ctx, util.ErrorArticleNotFound)
		return
	}

	article := h.articleSvc.Get(gDto.ID)
	if article == nil || article.Status == model.StatusDeleted {
		h.Fail(ctx, util.ErrorArticleNotFound)
		return
	}

	if article.UserId != user.ID {
		h.Fail(ctx, util.NewErrorMsg("无权限"))
		return
	}

	var articleForm dto.ArticleUpdateForm
	if h.BindAndValidate(ctx, &articleForm) {
		articleForm.ID = article.ID
		err := h.articleSvc.Update(articleForm)
		if err != nil {
			h.Fail(ctx, util.FromError(err))
			return
		}
		h.Success(ctx, gin.H{
			"articleId": article.ID,
		})
	}
}

// GetRecent 最近文章
func (h *ArticleHandler) GetRecent(ctx *gin.Context) {
	articles := h.articleSvc.Find(querybuilder.NewQueryBuilder().Where("status = ?", model.StatusOk).Desc("id").Limit(10))

	h.Success(ctx, articles)
}

// List 文章列表
func (h *ArticleHandler) List(ctx *gin.Context) {
	cursor := util.FormInt64Default(ctx, "cursor", 0)
	articles, cursor := h.articleSvc.GetArticles(cursor)
	h.Success(ctx, gin.H{
		"results": render.ToSimpleArticles(articles, h.articleTagCache, h.tagCache),
		"cursor":  strconv.FormatInt(cursor, 10),
	})
}

// GetTagArticles 标签文章列表
func (h *ArticleHandler) GetTagArticles(ctx *gin.Context) {
	var gDto dto.IdRequest
	if h.BindAndValidate(ctx, &gDto) {
		cursor := util.FormInt64Default(ctx, "cursor", 0)
		articles, cursor := h.articleSvc.GetTagArticles(gDto.ID, cursor)
		h.Success(ctx, gin.H{
			"results": render.ToSimpleArticles(articles, h.articleTagCache, h.tagCache),
			"cusor":   strconv.FormatInt(cursor, 10),
		})
	}
}

// GetUserRecent 用户最近的文章
func (h *ArticleHandler) GetUserRecent(ctx *gin.Context) {
	var gDto dto.IdRequest
	if h.BindAndValidate(ctx, &gDto) {
		articles := h.articleSvc.Find(querybuilder.NewQueryBuilder().Where("user_id = ? and status = ?",
			gDto.ID, model.StatusOk).Desc("id").Limit(10))
		h.Success(ctx, render.ToSimpleArticles(articles, h.articleTagCache, h.tagCache))
	}
}

// GetUserArticles 用户的文章
func (h *ArticleHandler) GetUserArticles(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	var gDto dto.IdRequest
	if h.BindAndValidate(ctx, &gDto) {
		articles, paging := h.articleSvc.List(querybuilder.NewQueryBuilder().
			Eq("user_id", gDto.ID).
			Eq("status", model.StatusOk).
			Page(page, 20).Desc("id"))

		h.Success(ctx, gin.H{
			"results": render.ToSimpleArticles(articles, h.articleTagCache, h.tagCache),
			"page":    paging,
		})
	}
}

// GetUserNewestBy 用户最新的文章
func (h *ArticleHandler) GetUserNewestBy(ctx *gin.Context) {
	var gDto dto.IdRequest
	if h.BindAndValidate(ctx, &gDto) {
		newestArticles := h.articleSvc.GetUserNewestArticles(gDto.ID)
		h.Success(ctx, render.ToSimpleArticles(newestArticles, h.articleTagCache, h.tagCache))
	}
}

// GetRelatedBy 相关文章
func (h *ArticleHandler) GetRelatedBy(ctx *gin.Context) {
	var gDto dto.IdRequest
	if h.BindAndValidate(ctx, &gDto) {
		relatedArticles := h.articleSvc.GetRelatedArticles(gDto.ID)
		h.Success(ctx, render.ToSimpleArticles(relatedArticles, h.articleTagCache, h.tagCache))
	}
}

// Favorite 收藏文章
func (h *ArticleHandler) Favorite(ctx *gin.Context) {
	user := h.GetCurrentUser(ctx)
	var req dto.SlugRequest
	if !h.BindAndValidateUri(ctx, &req) {
		return
	}

	if err := h.favoriteSvc.AddArticleFavorite(user.ID, hashid.Slug2Id[model.Article](req.Slug)); err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}

	h.Success(ctx, nil)
}
