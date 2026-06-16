package render

import (
	"strings"

	"github.com/PuerkitoBio/goquery"

	"ultrathreads/model"
	"ultrathreads/util/hashid"
	"ultrathreads/util/strtrim"
	"ultrathreads/util"
	"ultrathreads/util/urls"
)

// FavoriteContext 封装 ToFavorite 所需的全部外部数据
// 由调用方（handler/service）提前查询好后注入，render 本身不做任何 I/O
type FavoriteContext struct {
	Article *model.Article // EntityType == Article 时提供
	Post    *model.Post    // EntityType == Post 时提供
	User    *model.User    // 关联实体的作者信息
}

// ToFavorite 纯函数：根据收藏记录 + 预加载的上下文构建响应
func ToFavorite(favorite *model.Favorite, ctx *FavoriteContext) *model.FavoriteResponse {
	if favorite == nil {
		return nil
	}

	rsp := &model.FavoriteResponse{
		Slug:       hashid.Id2Slug[model.Favorite](favorite.ID),
		EntityType: favorite.EntityType,
		EntityId:   favorite.EntityId,
		CreateTime: favorite.CreateTime,
	}

	// 没有上下文 → 关联实体已被删除或查询失败
	if ctx == nil {
		rsp.Deleted = true
		return rsp
	}

	switch favorite.EntityType {
	case model.EntityTypeArticle:
		if ctx.Article == nil || ctx.Article.Status != model.StatusOk {
			rsp.Deleted = true
			return rsp
		}
		rsp.Url = urls.ArticleUrl(ctx.Article.ID)
		rsp.Title = ctx.Article.Title
		rsp.User = ToUser(ctx.User)
		if ctx.Article.ContentType == model.ContentTypeMarkdown {
			rsp.Content = util.GetMarkdownSummary(ctx.Article.Content)
		} else {
			rsp.Content = getHtmlTextSummary(ctx.Article.Content, 256)
		}

	case model.EntityTypePost:
		if ctx.Post == nil || ctx.Post.Status != model.StatusOk {
			rsp.Deleted = true
			return rsp
		}
		postSlug := hashid.Id2Slug[model.Post](ctx.Post.ID)
		rsp.Url = urls.PostUrl(postSlug)
		rsp.Title = ctx.Post.Title
		rsp.User = ToUser(ctx.User)
		rsp.Content = util.GetMarkdownSummary(ctx.Post.Content)

	default:
		rsp.Deleted = true
	}

	return rsp
}

// ToFavorites 批量构建：调用方需提前准备好每个 Favorite 对应的 Context
// contexts 与 favorites 按索引一一对应；某项为 nil 则标记为 Deleted
func ToFavorites(favorites []model.Favorite, contexts []*FavoriteContext) []model.FavoriteResponse {
	if len(favorites) == 0 {
		return nil
	}

	responses := make([]model.FavoriteResponse, 0, len(favorites))
	for i, fav := range favorites {
		var ctx *FavoriteContext
		if i < len(contexts) {
			ctx = contexts[i]
		}
		if r := ToFavorite(&fav, ctx); r != nil {
			responses = append(responses, *r)
		}
	}
	return responses
}

// ---------- 内部辅助函数 ----------

// getHtmlTextSummary 从 HTML 中提取纯文本摘要（纯函数，不依赖 service/cache）
func getHtmlTextSummary(html string, maxLen int) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	return strtrim.GetTextSummary(doc.Text(), maxLen)
}