package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"ultrathreads/model"
)

const (
	// ContextKeyUserClaims 是存储基础用户凭证的 Context Key
	// 供 CurrentUserReadState 及业务 Handler 读取
	ContextKeyUserClaims = "user_claims"

	// ContextKeyCurrentUser 是存储完整用户对象的 Context Key
	// 供 Controller 层读取
	ContextKeyCurrentUser = "CurrentUser"
)

// OptionalAuth 可选鉴权中间件
// Token 有效时注入用户信息；Token 缺失或无效时不拦截请求，仅跳过注入
func OptionalAuth(auth *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := auth.GetClaimsFromJWT(c)
		if err != nil || len(claims) == 0 {
			c.Next()
			return
		}

		userClaims := extractUserClaims(claims)
		if userClaims == nil {
			c.Next()
			return
		}

		// 写入 identityKey 供 gin-jwt 生态及 GetCurrent 使用
		identityKey := viper.GetString("jwt.identity_key")
		c.Set(identityKey, *userClaims)

		// 写入统一 Key 供下游中间件和业务层读取
		c.Set(ContextKeyUserClaims, *userClaims)

		// 获取完整用户对象并注入 Context
		if user := GetCurrent(c); user != nil {
			c.Set(ContextKeyCurrentUser, user)
		}

		c.Next()
	}
}

// extractUserClaims 将 JWT MapClaims 安全转换为 model.UserClaims
// 返回 nil 表示 Claims 格式异常或缺少必要字段
func extractUserClaims(claims jwt.MapClaims) *model.UserClaims {
	idVal, exists := claims["id"]
	if !exists {
		return nil
	}

	// JWT 数字类型默认为 float64，需安全转换
	id, ok := idVal.(float64)
	if !ok || id <= 0 {
		return nil
	}

	name, _ := claims["name"].(string)

	return &model.UserClaims{
		ID:   int64(id),
		Name: name,
	}
}