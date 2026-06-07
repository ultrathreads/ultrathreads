package converter

import (
	"html/template"

	"ultrathreads/cache"
	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util"
	"ultrathreads/util/log"
	"ultrathreads/util/markdown"
)

func ToPost(post *model.Post) *model.PostResponse {
	if post == nil {
		return nil
	}

	rsp := &model.PostResponse{}

	rsp.Id = post.ID
	rsp.Type = post.Type
	rsp.Title = post.Title
	rsp.User = ToUserDefaultIfNull(post.UserId)
	rsp.LastCommentTime = post.LastCommentTime
	rsp.CreateTime = post.CreateTime
	rsp.ViewCount = post.ViewCount
	rsp.LikeCount = post.LikeCount

	if post.NodeId > 0 {
		node := service.NodeService.Get(post.NodeId)
		rsp.Node = ToNode(node)
	}

	tags := service.PostService.GetPostTags(post.ID)
	rsp.Tags = ToTags(tags)

	mr := markdown.NewMd(markdown.MdWithTOC()).Run(post.Content)
	rsp.Content = template.HTML(ToHtmlContent(mr.ContentHtml))
	rsp.Toc = template.HTML(mr.TocHtml)

	if len(post.ImageList) > 0 {
		if err := util.ParseJson(post.ImageList, &rsp.ImageList); err != nil {
			log.Error(err.Error())
		}
	}

	return rsp
}

func ToSimplePost(post *model.Post) *model.PostSimpleResponse {
	if post == nil {
		return nil
	}

	rsp := &model.PostSimpleResponse{}

	rsp.Id = post.ID
	rsp.ThreadId = post.ThreadId
	rsp.ParentId = post.ParentId
	rsp.Type = post.Type
	rsp.Title = post.Title
	rsp.User = ToUserDefaultIfNull(post.UserId)
	rsp.LastCommentUser = ToUserDefaultIfNull(post.LastCommentUserId)
	rsp.LastCommentTime = post.LastCommentTime
	rsp.CreateTime = post.CreateTime
	rsp.ViewCount = post.ViewCount
	rsp.LikeCount = post.LikeCount

	if len(post.ImageList) > 0 {
		if err := util.ParseJson(post.ImageList, &rsp.ImageList); err != nil {
			log.Error(err.Error())
		}
	}

	if post.NodeId > 0 {
		node := cache.NodeCache.Get(post.NodeId)
		rsp.Node = ToNode(node)
	}

	tags := service.PostService.GetPostTags(post.ID)
	rsp.Tags = ToTags(tags)
	return rsp
}

func ToSimplePosts(posts []model.Post) []model.PostSimpleResponse {
	if posts == nil || len(posts) == 0 {
		return nil
	}
	var responses []model.PostSimpleResponse
	for _, post := range posts {
		responses = append(responses, *ToSimplePost(&post))
	}
	return responses
}
