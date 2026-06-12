package converter

import (
	"html/template"

	"ultrathreads/cache"
	"ultrathreads/model"
	"ultrathreads/util/markdown"
	"ultrathreads/util/strtrim"
	"ultrathreads/util/hashid"
)

func ToArticle(article *model.Article) *model.ArticleResponse {
	if article == nil {
		return nil
	}
	slug := hashid.Id2Slug[model.Article](article.ID)
	rsp := &model.ArticleResponse{}
	rsp.Slug = slug
	rsp.Title = article.Title
	rsp.Summary = article.Summary
	rsp.Share = article.Share
	rsp.SourceUrl = article.SourceUrl
	rsp.ViewCount = article.ViewCount
	rsp.CreateTime = article.CreateTime

	rsp.User = ToUserDefaultIfNull(article.UserId)

	tagIds := cache.ArticleTagCache.Get(article.ID)
	tags := cache.TagCache.GetList(tagIds)
	rsp.Tags = ToTags(tags)

	if article.ContentType == model.ContentTypeMarkdown {
		mr := markdown.NewMd(markdown.MdWithTOC()).Run(article.Content)
		rsp.Content = template.HTML(ToHtmlContent(mr.ContentHtml))
		rsp.Toc = template.HTML(mr.TocHtml)
		if len(rsp.Summary) == 0 {
			rsp.Summary = mr.SummaryText
		}
	} else {
		rsp.Content = template.HTML(ToHtmlContent(article.Content))
		if len(rsp.Summary) == 0 {
			rsp.Summary = strtrim.GetTextSummary(article.Content, 256)
		}
	}

	return rsp
}

func ToArticles(articles []model.Article) []model.ArticleResponse {
	if len(articles) == 0 {
		return []model.ArticleResponse{}
	}
	responses := make([]model.ArticleResponse, 0, len(articles))
	for i := range articles {
		if r := ToArticle(&articles[i]); r != nil {
			responses = append(responses, *r)
		}
	}
	return responses
}

func ToSimpleArticle(article *model.Article) *model.ArticleSimpleResponse {
	if article == nil {
		return nil
	}

	rsp := &model.ArticleSimpleResponse{}
	rsp.Title = article.Title
	rsp.Summary = article.Summary
	rsp.Share = article.Share
	rsp.SourceUrl = article.SourceUrl
	rsp.ViewCount = article.ViewCount
	rsp.CreateTime = article.CreateTime

	rsp.User = ToUserDefaultIfNull(article.UserId)

	tagIds := cache.ArticleTagCache.Get(article.ID)
	tags := cache.TagCache.GetList(tagIds)
	rsp.Tags = ToTags(tags)

	// 列表页仅需 Summary，无需解析完整 Markdown + TOC
	if len(rsp.Summary) == 0 {
		if article.ContentType == model.ContentTypeMarkdown {
			mr := markdown.NewMd().Run(article.Content) // 不传 MdWithTOC()
			rsp.Summary = mr.SummaryText
		} else {
			rsp.Summary = strtrim.GetTextSummary(strtrim.GetHtmlText(article.Content), 256)
		}
	}

	return rsp
}

func ToSimpleArticles(articles []model.Article) []model.ArticleSimpleResponse {
	if len(articles) == 0 {
		return []model.ArticleSimpleResponse{}
	}
	responses := make([]model.ArticleSimpleResponse, 0, len(articles))
	for i := range articles {
		if r := ToSimpleArticle(&articles[i]); r != nil {
			responses = append(responses, *r)
		}
	}
	return responses
}