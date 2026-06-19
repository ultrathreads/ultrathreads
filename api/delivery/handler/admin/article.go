package admin

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"

	"ultrathreads/cache"
	"ultrathreads/delivery/handler/base"
	"ultrathreads/domain"
	"ultrathreads/dto"
	"ultrathreads/model"
	"ultrathreads/render"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/markdown"
	"ultrathreads/util/querybuilder"
	"ultrathreads/util/strtrim"
)

// ArticleHandler article controller
type ArticleHandler struct {
	base.BaseHandler
	articleSvc      service.ArticleService
	articleTagCache cache.ArticleTagCacheInterface
	tagCache        cache.TagCacheInterface
}

func NewArticleHandler(articleSvc service.ArticleService, articleTagCache cache.ArticleTagCacheInterface, tagCache cache.TagCacheInterface) *ArticleHandler {
	return &ArticleHandler{
		articleSvc:      articleSvc,
		articleTagCache: articleTagCache,
		tagCache:        tagCache,
	}
}

// Show show article
func (h *ArticleHandler) Show(ctx *gin.Context) {
	var gDto dto.IdRequest
	if h.BindAndValidate(ctx, &gDto) {
		article := h.articleSvc.Get(gDto.ID)
		if article == nil {
			h.Fail(ctx, util.NewErrorMsg("Article not found, id="+strconv.FormatInt(gDto.ID, 10)))
			return
		}
		h.Success(ctx, article)
	}
}

// Update update a article
func (h *ArticleHandler) Update(ctx *gin.Context) {
	var gDto dto.IdRequest
	if !h.BindAndValidate(ctx, &gDto) {
		return
	}
	article := h.articleSvc.Get(gDto.ID)
	if article == nil {
		h.Fail(ctx, util.NewErrorMsg("Article not found, id="+strconv.FormatInt(gDto.ID, 10)))
		return
	}

	var articleForm dto.ArticleUpdateForm
	if !h.BindAndValidate(ctx, &articleForm) {
		return
	}
	articleForm.ID = gDto.ID
	cmd := domain.UpdateArticleCommand{
		ID:      articleForm.ID,
		Title:   articleForm.Title,
		Summary: articleForm.Summary,
		Content: articleForm.Content,
		Tags:    articleForm.Tags,
	}
	err := h.articleSvc.Update(cmd)
	if err != nil {
		h.Fail(ctx, util.FromError(err))
		return
	}
	h.Success(ctx, article)
}

// Delete delete article
func (h *ArticleHandler) Delete(ctx *gin.Context) {
	var gDto dto.IdRequest
	if !h.BindAndValidate(ctx, &gDto) {
		return
	}
	h.articleSvc.Delete(gDto.ID)
	h.Success(ctx, nil)
}

// List list articles
func (h *ArticleHandler) List(ctx *gin.Context) {
	page := util.FormIntDefault(ctx, "page", 1)
	limit := util.FormIntDefault(ctx, "limit", 20)
	name := ctx.Request.FormValue("name")

	conditions := querybuilder.NewQueryBuilder()
	if len(name) > 0 {
		conditions.Like("name", name)
	}
	list, paging := h.articleSvc.List(conditions.Page(page, limit).Desc("id"))

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
		tagIds := h.articleTagCache.Get(article.ID)
		tags := h.tagCache.GetList(tagIds)
		item["tags"] = render.ToTags(tags)

		results = append(results, item)
	}

	h.Success(ctx, &querybuilder.PageResult{Results: results, Page: paging})
}
