package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"ultrathreads/dto"
	"ultrathreads/model"
	"ultrathreads/render"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/hashid"
	"ultrathreads/util/querybuilder"
)

type ArticleController struct {
	BaseController
	articleSvc  service.ArticleServicer
	favoriteSvc service.FavoriteServicer
}

func NewArticleController(articleSvc service.ArticleServicer, favoriteSvc service.FavoriteServicer) *ArticleController {
	return &ArticleController{articleSvc: articleSvc, favoriteSvc: favoriteSvc}
}

// Show show article by id
func (c *ArticleController) Show(ctx *gin.Context) {
	var gDto dto.IdRequest
	if c.BindAndValidate(ctx, &gDto) {
		article := c.articleSvc.Get(gDto.ID)
		if article == nil || article.Status != model.StatusOk {
			c.Fail(ctx, util.ErrorArticleNotFound)
			return
		}
		c.Success(ctx, render.ToArticle(article))
	}
}

// Store 发表文章
func (c *ArticleController) Store(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		c.Fail(ctx, util.ErrorNotLogin)
		return
	}
	var articleForm dto.ArticleCreateForm
	if c.BindAndValidate(ctx, &articleForm) {
		articleForm.UserID = user.ID
		article, err := c.articleSvc.Create(articleForm)
		if err != nil {
			c.Fail(ctx, util.FromError(err))
			return
		}
		c.Success(ctx, render.ToArticle(article))
	}
}

// Edit 编辑时获取详情
func (c *ArticleController) Edit(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		c.Fail(ctx, util.ErrorNotLogin)
		return
	}

	var gDto dto.IdRequest
	if c.BindAndValidate(ctx, &gDto) {
		article := c.articleSvc.Get(gDto.ID)

		if article == nil || article.Status != model.StatusOk {
			c.Fail(ctx, util.NewErrorMsg("话题不存在或已被删除"))
			return
		}
		if article.UserId != user.ID {
			c.Fail(ctx, util.NewErrorMsg("无权限"))
			return
		}

		tags := c.articleSvc.GetArticleTags(article.ID)
		var tagNames []string
		if len(tags) > 0 {
			for _, tag := range tags {
				tagNames = append(tagNames, tag.Name)
			}
		}

		c.Success(ctx, gin.H{
			"articleId": article.ID,
			"title":     article.Title,
			"content":   article.Content,
			"tags":      tagNames,
		})
	}
}

// Update 编辑文章
func (c *ArticleController) Update(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		c.Fail(ctx, util.ErrorNotLogin)
		return
	}
	var gDto dto.IdRequest
	if !c.BindAndValidate(ctx, &gDto) {
		c.Fail(ctx, util.ErrorArticleNotFound)
		return
	}

	article := c.articleSvc.Get(gDto.ID)
	if article == nil || article.Status == model.StatusDeleted {
		c.Fail(ctx, util.ErrorArticleNotFound)
		return
	}

	if article.UserId != user.ID {
		c.Fail(ctx, util.NewErrorMsg("无权限"))
		return
	}

	var articleForm dto.ArticleUpdateForm
	if c.BindAndValidate(ctx, &articleForm) {
		articleForm.ID = article.ID
		err := c.articleSvc.Update(articleForm)
		if err != nil {
			c.Fail(ctx, util.FromError(err))
			return
		}
		c.Success(ctx, gin.H{
			"articleId": article.ID,
		})
	}
}

// GetRecent 最近文章
func (c *ArticleController) GetRecent(ctx *gin.Context) {
	articles := c.articleSvc.Find(querybuilder.NewQueryBuilder().Where("status = ?", model.StatusOk).Desc("id").Limit(10))

	c.Success(ctx, articles)
}

// List 文章列表
func (c *ArticleController) List(ctx *gin.Context) {
	cursor := util.FormInt64Default(ctx, "cursor", 0)
	articles, cursor := c.articleSvc.GetArticles(cursor)
	c.Success(ctx, gin.H{
		"results": render.ToSimpleArticles(articles),
		"cursor":  strconv.FormatInt(cursor, 10),
	})
}

// GetTagArticles 标签文章列表
func (c *ArticleController) GetTagArticles(ctx *gin.Context) {
	var gDto dto.IdRequest
	if c.BindAndValidate(ctx, &gDto) {
		cursor := util.FormInt64Default(ctx, "cursor", 0)
		articles, cursor := c.articleSvc.GetTagArticles(gDto.ID, cursor)
		c.Success(ctx, gin.H{
			"results": render.ToSimpleArticles(articles),
			"cusor":   strconv.FormatInt(cursor, 10),
		})
	}
}

// GetUserRecent 用户最近的文章
func (c *ArticleController) GetUserRecent(ctx *gin.Context) {
	var gDto dto.IdRequest
	if c.BindAndValidate(ctx, &gDto) {
		articles := c.articleSvc.Find(querybuilder.NewQueryBuilder().Where("user_id = ? and status = ?",
			gDto.ID, model.StatusOk).Desc("id").Limit(10))
		c.Success(ctx, render.ToSimpleArticles(articles))
	}
}

// GetUserArticles 用户的文章
func (c *ArticleController) GetUserArticles(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	var gDto dto.IdRequest
	if c.BindAndValidate(ctx, &gDto) {
		articles, paging := c.articleSvc.List(querybuilder.NewQueryBuilder().
			Eq("user_id", gDto.ID).
			Eq("status", model.StatusOk).
			Page(page, 20).Desc("id"))

		c.Success(ctx, gin.H{
			"results": render.ToSimpleArticles(articles),
			"page":    paging,
		})
	}
}

// GetUserNewestBy 用户最新的文章
func (c *ArticleController) GetUserNewestBy(ctx *gin.Context) {
	var gDto dto.IdRequest
	if c.BindAndValidate(ctx, &gDto) {
		newestArticles := c.articleSvc.GetUserNewestArticles(gDto.ID)
		c.Success(ctx, render.ToSimpleArticles(newestArticles))
	}
}

// GetRelatedBy 相关文章
func (c *ArticleController) GetRelatedBy(ctx *gin.Context) {
	var gDto dto.IdRequest
	if c.BindAndValidate(ctx, &gDto) {
		relatedArticles := c.articleSvc.GetRelatedArticles(gDto.ID)
		c.Success(ctx, render.ToSimpleArticles(relatedArticles))
	}
}

// Favorite 收藏文章
func (c *ArticleController) Favorite(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var req dto.SlugRequest
	if !c.BindAndValidateUri(ctx, &req) {
		return
	}

	if err := c.favoriteSvc.AddArticleFavorite(user.ID, hashid.Slug2Id[model.Article](req.Slug)); err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}

	c.Success(ctx, nil)
}
