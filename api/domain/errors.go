package domain

import "errors"

// 通用错误
var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrInvalidInput  = errors.New("invalid input")
)

// 用户相关错误
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserDisabled      = errors.New("user is disabled")
	ErrPasswordIncorrect = errors.New("password incorrect")
	ErrEmailInUse        = errors.New("email already in use")
)

// 帖子相关错误
var (
	ErrPostNotFound = errors.New("post not found")
	ErrPostClosed   = errors.New("post is closed")
)

// 文章相关错误
var (
	ErrArticleNotFound = errors.New("article not found")
)

// 节点相关错误
var (
	ErrNodeNotFound = errors.New("node not found")
)

// 标签相关错误
var (
	ErrTagNotFound = errors.New("tag not found")
)

// 收藏相关错误
var (
	ErrFavoriteNotFound = errors.New("favorite not found")
)

// 通知相关错误
var (
	ErrNotificationNotFound = errors.New("notification not found")
)
