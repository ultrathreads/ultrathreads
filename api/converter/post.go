package converter

import (
	"html/template"

	"ultrathreads/cache"
	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/hashid"
	"ultrathreads/util/log"
	"ultrathreads/util/markdown"
)

// basePostFields 提取两个响应共有的字段赋值逻辑
func basePostFields(rsp *model.PostSimpleResponse, post *model.Post) {
	slug := hashid.Id2Slug[model.Post](post.ID)
	parentSlug := hashid.Id2Slug[model.Post](post.ParentId)
	threadSlug := hashid.Id2Slug[model.Post](post.ThreadId)
	rsp.Slug = slug
	rsp.Type = post.Type
	rsp.ThreadSlug = threadSlug
	rsp.ParentSlug = parentSlug
	rsp.IsRoot = post.IsRoot()
	rsp.Title = post.Title
	rsp.IsPinned = post.IsPinned
	rsp.User = ToUserDefaultIfNull(post.UserId)
	rsp.LastCommentTime = post.LastCommentTime
	rsp.CreateTime = post.CreateTime
	rsp.ViewCount = post.ViewCount
	rsp.LikeCount = post.LikeCount

	if post.IsRoot() {
		tags := service.Srv.PostService.GetPostTags(post.ID)
		rsp.Tags = ToTags(tags)
	}
}

func ToPost(post *model.Post) *model.PostResponse {
	if post == nil {
		return nil
	}

	rsp := &model.PostResponse{}
	basePostFields(&rsp.PostSimpleResponse, post)

	if len(post.ImageList) > 0 {
		if err := util.ParseJson(post.ImageList, &rsp.ImageList); err != nil {
			log.Error(err.Error())
		}
	}

	if post.NodeId > 0 {
		node := service.Srv.NodeService.Get(post.NodeId)
		rsp.Node = ToNode(node)
	}

	rsp.RawContent = post.Content // 供编辑使用
	mr := markdown.NewMd(markdown.MdWithTOC()).Run(post.Content)
	rsp.Content = template.HTML(ToHtmlContent(mr.ContentHtml))
	rsp.Toc = template.HTML(mr.TocHtml)

	return rsp
}

func ToSimplePost(post *model.Post) *model.PostSimpleResponse {
	if post == nil {
		return nil
	}

	rsp := &model.PostSimpleResponse{}
	basePostFields(rsp, post)

	// 列表页特有：LastCommentUser
	rsp.LastCommentUser = ToUserDefaultIfNull(post.LastCommentUserId)
	rsp.NodeSlug = hashid.Id2Slug[model.Node](post.NodeId)

	// 列表页特有：Node 走 Cache（高性能）
	if post.NodeId > 0 {
		node := cache.NodeCache.Get(post.NodeId)
		rsp.Node = ToNode(node)
	}

	return rsp
}

// ToSimplePosts 返回列表页帖子切片
func ToSimplePosts(posts []model.Post) []model.PostSimpleResponse {
	if len(posts) == 0 {
		return []model.PostSimpleResponse{}
	}
	responses := make([]model.PostSimpleResponse, 0, len(posts))
	for i := range posts {
		if r := ToSimplePost(&posts[i]); r != nil {
			responses = append(responses, *r)
		}
	}
	return responses
}

// ToPosts 返回详情页帖子切片（注意性能警告 ⚠️）
func ToPosts(posts []model.Post) []model.PostResponse {
	if len(posts) == 0 {
		return []model.PostResponse{}
	}
	responses := make([]model.PostResponse, 0, len(posts))
	for i := range posts {
		if r := ToPost(&posts[i]); r != nil {
			responses = append(responses, *r)
		}
	}
	return responses
}