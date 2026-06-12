package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ==================== 通用错误定义 ====================
var (
	ErrParamMissing = errors.New("param missing")
	ErrParamInvalid = errors.New("param invalid")
)

func newMissingError(name string) error {
	return fmt.Errorf("%w: '%s'", ErrParamMissing, name)
}

func newInvalidError(name, value string) error {
	return fmt.Errorf("%w: '%s'='%s'", ErrParamInvalid, name, value)
}

// ==================== FormValue (POST 表单 + Query 混合) ====================
// 📌 场景: POST 表单提交、x-www-form-urlencoded、multipart/form-data

func FormInt(ctx *gin.Context, name string) (int, error) {
	str := ctx.Request.FormValue(name)
	if str == "" {
		return 0, newMissingError(name)
	}
	v, err := strconv.Atoi(str)
	if err != nil {
		return 0, newInvalidError(name, str)
	}
	return v, nil
}

func FormIntDefault(ctx *gin.Context, name string, def int) int {
	if v, err := FormInt(ctx, name); err == nil {
		return v
	}
	return def
}

func FormInt64(ctx *gin.Context, name string) (int64, error) {
	str := ctx.Request.FormValue(name)
	if str == "" {
		return 0, newMissingError(name)
	}
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, newInvalidError(name, str)
	}
	return v, nil
}

func FormInt64Default(ctx *gin.Context, name string, def int64) int64 {
	if v, err := FormInt64(ctx, name); err == nil {
		return v
	}
	return def
}

// 🆕 FormBool
// 📌 场景: 开关类参数 (is_active=true, enable=1, checked=on)
func FormBool(ctx *gin.Context, name string) (bool, error) {
	str := strings.ToLower(strings.TrimSpace(ctx.Request.FormValue(name)))
	if str == "" {
		return false, newMissingError(name)
	}
	switch str {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, newInvalidError(name, str)
	}
}

func FormBoolDefault(ctx *gin.Context, name string, def bool) bool {
	if v, err := FormBool(ctx, name); err == nil {
		return v
	}
	return def
}

// 🆕 FormString (带非空校验)
// 📌 场景: 必填文本字段 (username, email)，区别于直接 ctx.Request.FormValue 返回空串不报错
func FormString(ctx *gin.Context, name string) (string, error) {
	str := strings.TrimSpace(ctx.Request.FormValue(name))
	if str == "" {
		return "", newMissingError(name)
	}
	return str, nil
}

func FormStringDefault(ctx *gin.Context, name string, def string) string {
	if v, err := FormString(ctx, name); err == nil {
		return v
	}
	return def
}

// 🆕 FormInt64Slice
// 📌 场景: 批量操作 (ids=1,2,3 或 ids=1&ids=2&ids=3)
func FormInt64Slice(ctx *gin.Context, name string) ([]int64, error) {
	raw := ctx.Request.FormValue(name)
	if raw == "" {
		return nil, newMissingError(name)
	}
	parts := strings.Split(raw, ",")
	result := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return nil, newInvalidError(name, raw)
		}
		result = append(result, v)
	}
	if len(result) == 0 {
		return nil, newInvalidError(name, raw)
	}
	return result, nil
}

// ==================== Param (URL 路径参数) ====================
// 📌 场景: RESTful 路由 (/users/:id, /orders/:order_no)

func ParamInt64(ctx *gin.Context, name string) (int64, error) {
	v := ctx.Param(name)
	if v == "" {
		return 0, newMissingError(name)
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, newInvalidError(name, v)
	}
	return n, nil
}

func ParamInt64Default(ctx *gin.Context, name string, def int64) int64 {
	if v, err := ParamInt64(ctx, name); err == nil {
		return v
	}
	return def
}

// 🆕 ParamString (路径参数通常不为空，但做防御性校验)
func ParamString(ctx *gin.Context, name string) (string, error) {
	v := strings.TrimSpace(ctx.Param(name))
	if v == "" {
		return "", newMissingError(name)
	}
	return v, nil
}

func ParamStringDefault(ctx *gin.Context, name string, def string) string {
	if v, err := ParamString(ctx, name); err == nil {
		return v
	}
	return def
}

// ==================== Query (纯 URL 查询参数) ====================
// 📌 场景: GET 列表筛选/分页 (?page=1&size=20&status=active)
// ⚠️ 注意: 与 FormValue 的区别是 Query 只读 URL ?后面的参数，不读 POST Body

func QueryInt(ctx *gin.Context, name string) (int, error) {
	str := ctx.Query(name)
	if str == "" {
		return 0, newMissingError(name)
	}
	v, err := strconv.Atoi(str)
	if err != nil {
		return 0, newInvalidError(name, str)
	}
	return v, nil
}

func QueryIntDefault(ctx *gin.Context, name string, def int) int {
	if v, err := QueryInt(ctx, name); err == nil {
		return v
	}
	return def
}

func QueryInt64(ctx *gin.Context, name string) (int64, error) {
	str := ctx.Query(name)
	if str == "" {
		return 0, newMissingError(name)
	}
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, newInvalidError(name, str)
	}
	return v, nil
}

func QueryInt64Default(ctx *gin.Context, name string, def int64) int64 {
	if v, err := QueryInt64(ctx, name); err == nil {
		return v
	}
	return def
}

func QueryBool(ctx *gin.Context, name string) (bool, error) {
	str := strings.ToLower(strings.TrimSpace(ctx.Query(name)))
	if str == "" {
		return false, newMissingError(name)
	}
	switch str {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, newInvalidError(name, str)
	}
}

func QueryBoolDefault(ctx *gin.Context, name string, def bool) bool {
	if v, err := QueryBool(ctx, name); err == nil {
		return v
	}
	return def
}

func QueryString(ctx *gin.Context, name string) (string, error) {
	str := strings.TrimSpace(ctx.Query(name))
	if str == "" {
		return "", newMissingError(name)
	}
	return str, nil
}

func QueryStringDefault(ctx *gin.Context, name string, def string) string {
	if v, err := QueryString(ctx, name); err == nil {
		return v
	}
	return def
}