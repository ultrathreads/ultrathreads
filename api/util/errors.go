package util

import (
	"strconv"
)

var (
	ErrorNotLogin         = NewError(401, "请先登录")
	ErrorPermissionDenied = NewError(403, "Permission denied.")

	ErrorResourceNotFound = NewError(404, "请求的资源不存在")
	ErrorPostNotFound     = NewError(404, "话题不存在")
	ErrorArticleNotFound  = NewError(404, "文章不存在")
	ErrorTagNotFound      = NewError(404, "标签不存在")

	ErrorCaptchaWrong     = NewError(422, "验证码错误")
)

func NewError(code int, text string) *CodeError {
	return &CodeError{code, text, nil}
}

func NewErrorMsg(text string) *CodeError {
	return &CodeError{-1, text, nil}
}

func NewErrorData(code int, text string, data interface{}) *CodeError {
	return &CodeError{code, text, data}
}

func FromError(err error) *CodeError {
	if err == nil {
		return nil
	}
	return &CodeError{-1, err.Error(), nil}
}

type CodeError struct {
	Code    int
	Message string
	Data    interface{}
}

func (e *CodeError) Error() string {
	return strconv.Itoa(e.Code) + ": " + e.Message
}
