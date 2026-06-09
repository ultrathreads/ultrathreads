package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"ultrathreads/model"
	"ultrathreads/service"
)

// OptionalAuth 可选鉴权中间件
// 与强制鉴权的区别：Token 无效或缺失时不拦截请求，仅跳过用户注入
func OptionalAuth(auth *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := auth.GetClaimsFromJWT(c)
		if err != nil || len(claims) == 0 {
			c.Next()
			return
		}

		// 从 JWT Claims 中提取用户ID并构造 Service 层期望的 UserClaims
		userClaims := extractUserClaims(claims)
		if userClaims == nil {
			c.Next()
			return
		}

		// 写入 Service 层读取的 Key，使 GetCurrent 能正常工作
		identityKey := viper.GetString("jwt.identity_key")
		c.Set(identityKey, *userClaims)

		// 调用 Service 获取完整用户对象，并写入 Controller 层读取的 Key
		user := service.UserService.GetCurrent(c)
		if user != nil {
			c.Set("CurrentUser", user)
		}

		c.Next()
	}
}

// extractUserClaims 将 JWT MapClaims 转换为项目自定义的 UserClaims
// ⚠️ 请根据 model.UserClaims 的实际字段调整此处逻辑
func extractUserClaims(claims jwt.MapClaims) *model.UserClaims {
	idVal, exists := claims["id"]
	if !exists {
		return nil
	}

	// JWT 解析出的数字默认为 float64，需安全转换为 int64
	f, ok := idVal.(float64)
	if !ok || f <= 0 {
		return nil
	}

	return &model.UserClaims{
		ID: int64(f),
		// 如 UserClaims 还有其他必要字段，在此补充：
		// Username: fmt.Sprintf("%v", claims["username"]),
	}
}