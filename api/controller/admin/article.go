package admin

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"

	"ultrathreads/render"
	"ultrathreads/cache"
	"ultrathreads/controller"
	"ultrathreads/dto"
	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/markdown"
	"ultrathreads/util/querybuilder"
	"ultrathreads/util/strtrim"
)

// ArticleController article controller
type ArticleController struct {
	controller.BaseController
	articleSvc service.ArticleServicer
}

func NewArticleController(articleSvc service.ArticleServicer) *ArticleController {
	return &ArticleController{articleSvc: articleSvc}
}

// Show show article
func (c *ArticleController) Show(ctx *gin.Context) {
	var gDto dto.IdRequest
	if c.BindAndValidate(ctx, &gDto) {
		article := c.articleSvc.Get(gDto.ID)
		if article == nil {
			c.Fail(ctx, util.NewErrorMsg("Article not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.Success(ctx, article)
	}
}

// Update update a article
func (c *ArticleController) Update(ctx *gin.Context) {
	var gDto dto.IdRequest
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	article := c.articleSvc.Get(gDto.ID)
	if article == nil {
		c.Fail(ctx, util.NewErrorMsg("Article not found, id="+strconv.FormatInt(gDto.ID, 10)))
		return
	}

	var articleForm dto.ArticleUpdateForm
	if !c.BindAndValidate(ctx, &articleForm) {
		return
	}
	articleForm.ID = gDto.ID
	err := c.articleSvc.Update(articleForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, article)
}

// Delete delete article
func (c *ArticleController) Delete(ctx *gin.Context) {
	var gDto dto.IdRequest
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	c.articleSvc.Delete(gDto.ID)
	c.Success(ctx, nil)
}

// List list articles
func (c *ArticleController) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	limit := util.FormIntDefault(ctx, "limit", 20)
	name := ctx.Request.FormValue("name")

	conditions := querybuilder.NewQueryBuilder()
	if len(name) > 0 {
		conditions.Like("name", name)
	}
	list, paging := c.articleSvc.List(conditions.Page(page, limit).Desc("id"))

	var results []map[string]interface{}
	for _, article := range list {
		item := util.StructToMap(article, "content")
		item["user"] = render.ToDefaultUser(article.UserId)

		if article.ContentType == model.ContentTypeMarkdown {
			mr := markdown.NewMd().Run(article.Content)
			if len(article.Summary) == 0 {
				item["summary"] = mr.SummaryText
			}
		} else {
			if len(article.Summary) == 0 {
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(article.Content))
				if err != nil {
					item["summary"] = strtrim.GetTextSummary(doc.Text(), 256)
				}
			}
		}
		tagIds := cache.ArticleTagCache.Get(article.ID)
		tags := cache.TagCache.GetList(tagIds)
		item["tags"] = render.ToTags(tags)

		results = append(results, item)
	}

	c.Success(ctx, &querybuilder.PageResult{Results: results, Page: paging})
}