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

// basePostFields 提取两个响应共有的字段赋值逻辑
// simpleRsp 作为基础载体传入，避免重复定义公共字段
func basePostFields(rsp *model.PostSimpleResponse, post *model.Post) {
	rsp.Id = post.ID
	rsp.Type = post.Type
	rsp.ThreadId = post.ThreadId
	rsp.ParentId = post.ParentId
	rsp.Title = post.Title
	rsp.IsPinned = post.IsPinned
	rsp.User = ToUserDefaultIfNull(post.UserId)
	rsp.LastCommentTime = post.LastCommentTime
	rsp.CreateTime = post.CreateTime
	rsp.ViewCount = post.ViewCount
	rsp.LikeCount = post.LikeCount

	if len(post.ImageList) > 0 {
		if err := util.ParseJson(post.ImageList, &rsp.ImageList); err != nil {
			log.Error(err.Error())
		}
	}

	tags := service.PostService.GetPostTags(post.ID)
	rsp.Tags = ToTags(tags)
}

func ToPost(post *model.Post) *model.PostResponse {
	if post == nil {
		return nil
	}

	rsp := &model.PostResponse{}

	// 填充公共字段（PostResponse 内嵌或包含 PostSimpleResponse 的所有公共字段）
	basePostFields(&rsp.PostSimpleResponse, post)

	// 详情页特有：Node 走 Service（可能需要实时数据）
	if post.NodeId > 0 {
		node := service.NodeService.Get(post.NodeId)
		rsp.Node = ToNode(node)
	}

	// 详情页特有：Markdown 渲染
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

	// 列表页特有：Node 走 Cache（高性能）
	if post.NodeId > 0 {
		node := cache.NodeCache.Get(post.NodeId)
		rsp.Node = ToNode(node)
	}

	return rsp
}

func ToSimplePosts(posts []model.Post) []model.PostSimpleResponse {
	if len(posts) == 0 {
		return nil
	}
	responses := make([]model.PostSimpleResponse, 0, len(posts))
	for i := range posts {
		responses = append(responses, *ToSimplePost(&posts[i]))
	}
	return responses
}

func ToPosts(posts []model.Post) []model.PostResponse {
	if len(posts) == 0 {
		return nil
	}
	responses := make([]model.PostResponse, 0, len(posts))
	for i := range posts {
		responses = append(responses, *ToPost(&posts[i]))
	}
	return responses
}