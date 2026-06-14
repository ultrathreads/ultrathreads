package urls

import (
	"github.com/spf13/viper"
	"net/url"
	"strconv"
	"strings"

	"ultrathreads/util/log"
)

var BaseUrl = viper.GetString("base.url")

// 是否是内部链接
func IsInternalUrl(href string) bool {
	if IsAnchor(href) {
		return true
	}
	u, err := url.Parse(BaseUrl)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	return strings.Contains(href, u.Host)
}

// 是否是锚链接
func IsAnchor(href string) bool {
	return strings.Index(href, "#") == 0
}

func AbsUrl(path string) string {
	return BaseUrl + path
}

// 用户主页
func UserUrl(userId int64) string {
	return AbsUrl("/user/" + strconv.FormatInt(userId, 10))
}

// 文章详情
func ArticleUrl(articleId int64) string {
	return AbsUrl("/article/" + strconv.FormatInt(articleId, 10))
}

// 标签文章列表
func TagArticlesUrl(tagId int64) string {
	return AbsUrl("/articles/" + strconv.FormatInt(tagId, 10))
}

// 话题详情
func PostUrl(postSlug string) string {
	return AbsUrl("/threads/" + postSlug)
}

func UrlJoin(parts ...string) string {
	sep := "/"
	var ss []string
	for i, part := range parts {
		part = strings.TrimSpace(part)
		var (
			from = 0
			to   = len(part)
		)
		if strings.Index(part, sep) == 0 {
			from = 1
		}
		if strings.LastIndex(part, sep) == len(part)-1 {
			to = len(part) - 1
		}
		part = part[from:to]

		ss = append(ss, part)
		if i != len(parts)-1 {
			ss = append(ss, sep)
		}
	}
	return strings.Join(ss, "")
}
