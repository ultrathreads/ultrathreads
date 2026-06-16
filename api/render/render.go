package render

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"

	"ultrathreads/model"
	"ultrathreads/util"
	//"ultrathreads/util/avatar"
	"ultrathreads/util/hashid"
	"ultrathreads/util/urls"
)

func ToNotification(notification *model.Notification) *model.NotificationResponse {
	if notification == nil {
		return nil
	}

	detailUrl := ""
	icon := ""
	if notification.Type == model.MsgTypeComment {
		entityType := gjson.Get(notification.ExtraData, "entityType")
		entityId := gjson.Get(notification.ExtraData, "entityId")
		if entityType.String() == model.EntityTypeArticle {
			detailUrl = urls.ArticleUrl(entityId.Int())
		} else if entityType.String() == model.EntityTypePost {
			postSlug := hashid.Id2Slug[model.Post](entityId.Int())
			detailUrl = urls.PostUrl(postSlug)
		}
		icon = "comment"
	} else if notification.Type == model.MsgTypePostLike {
		entityId := gjson.Get(notification.ExtraData, "entityId")
		postSlug := hashid.Id2Slug[model.Post](entityId.Int())
		detailUrl = urls.PostUrl(postSlug)
		icon = "heart"
	} else if notification.Type == model.MsgTypeUserWatch {
		entityId := gjson.Get(notification.ExtraData, "entityId")
		detailUrl = urls.UserUrl(entityId.Int())
		icon = "eye"
	}
	/*
	from := ToUserDefaultIfNull(notification.FromId)
	if notification.FromId <= 0 {
		from.Nickname = "系统通知"
		from.Avatar = avatar.DefaultAvatar
	}
	*/

	return &model.NotificationResponse{
		MessageId:    notification.ID,
		//From:         from,
		UserId:       notification.UserId,
		Content:      notification.Content,
		QuoteContent: notification.QuoteContent,
		Type:         notification.Type,
		Icon:         icon,
		DetailUrl:    detailUrl,
		ExtraData:    notification.ExtraData,
		Status:       notification.Status,
		CreateTime:   notification.CreateTime,
	}
}

func ToNotifications(notifications []model.Notification) []model.NotificationResponse {
	if len(notifications) == 0 {
		return nil
	}
	var responses []model.NotificationResponse
	for _, notification := range notifications {
		responses = append(responses, *ToNotification(&notification))
	}
	return responses
}

func ToHtmlContent(htmlContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent
	}

	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		href := selection.AttrOr("href", "")

		if len(href) == 0 {
			return
		}

		// 不是内部链接
		if !urls.IsInternalUrl(href) {
			selection.SetAttr("target", "_blank")
			selection.SetAttr("rel", "external nofollow") // 标记站外链接，搜索引擎爬虫不传递权重值
		}

		// 如果是锚链接
		if urls.IsAnchor(href) {
			selection.ReplaceWithHtml(selection.Text())
		}

		// 如果a标签没有title，那么设置title
		title := selection.AttrOr("title", "")
		if len(title) == 0 {
			selection.SetAttr("title", selection.Text())
		}
	})

	// 处理图片
	doc.Find("img").Each(func(i int, selection *goquery.Selection) {
		src := selection.AttrOr("src", "")
		// 处理第三方图片
		if strings.Contains(src, "qpic.cn") {
			src = util.ParseUrl("/api/img/proxy").AddQuery("url", src).BuildStr()
			// selection.SetAttr("src", src)
		}

		// 处理lazyload
		selection.SetAttr("data-src", src)
		selection.RemoveAttr("src")
	})

	html, err := doc.Find("body").Html()
	if err != nil {
		return htmlContent
	}
	return html
}
