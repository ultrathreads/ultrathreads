package middleware

import (
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"ultrathreads/form"
	"ultrathreads/model"
	"ultrathreads/service"
	"ultrathreads/util/log"
)

// Login type constants
var (
	LoginStandard = 1
	LoginOAuth    = 2
)

// LoginDto 标准登录请求参数
type LoginDto struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Code     string `form:"code" json:"code"`
}

// LoginOAuthDto OAuth 登录请求参数
type LoginOAuthDto struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}

// LoginResponseData 标准化的登录响应数据结构
type LoginResponseData struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpireAt     time.Time `json:"expire_at"`
}

// JwtAuth 初始化 JWT 中间件
func JwtAuth(LoginType int) *jwt.GinJWTMiddleware {
	jwtMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "Jwt",
		Key:         []byte(viper.GetString("jwt.key")),
		Timeout:     time.Hour * 24,
		MaxRefresh:  time.Hour * 24 * 90,
		IdentityKey: viper.GetString("jwt.identity_key"),

		// ✅ 改造：返回结构化双 Token 响应，显式包含 success 字段
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			expiresIn := int64(time.Until(expire).Seconds())
			if expiresIn < 0 {
				expiresIn = 0
			}

			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": "login success",
				"success": true,
				"data": LoginResponseData{
					AccessToken:  token,
					RefreshToken: token,
					ExpiresIn:    expiresIn,
					ExpireAt:     expire,
				},
			})
		},

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(model.UserClaims); ok {
				return jwt.MapClaims{
					"id":    v.ID,
					"name":  v.Name,
					"uid":   v.ID,
					"uname": v.Name,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			id, _ := claims["id"].(float64)
			name, _ := claims["name"].(string)
			return model.UserClaims{
				ID:   int64(id),
				Name: name,
			}
		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			if LoginType == LoginOAuth {
				return AuthenticatorOAuth(c)
			}
			return Authenticator(c)
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(model.UserClaims); ok {
				return true
			}
			return false
		},

		// ✅ 修复：未授权时返回标准 HTTP 401 + success=false
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": message,
				"success": false,
			})
		},

		TokenLookup:   "header: Authorization, query: token, cookie: access_token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		log.Error(err.Error())
	}
	return jwtMiddleware
}

// Authenticator 标准用户名密码登录验证
func Authenticator(c *gin.Context) (interface{}, error) {
	var loginDto LoginDto
	if err := form.Bind(c, &loginDto); err != nil {
		return "", err
	}

	log.Info("loginDto.Username: %s", loginDto.Username)

	ok, err, u := service.UserService.VerifyAndReturnUserInfo(loginDto.Username, loginDto.Password)
	if ok {
		return model.UserClaims{
			ID:   u.ID,
			Name: u.Username.String,
		}, nil
	}
	return nil, err
}

// AuthenticatorOAuth OAuth 第三方登录验证
func AuthenticatorOAuth(c *gin.Context) (interface{}, error) {
	provider := c.Param("provider")

	var oauthDto LoginOAuthDto
	if err := form.Bind(c, &oauthDto); err != nil {
		return "", err
	}

	account, err := service.LoginSourceService.GetOrCreate(provider, oauthDto.Code, oauthDto.State)
	if err != nil {
		return nil, err
	}

	u, err := service.UserService.SignInByLoginSource(account)
	if err == nil {
		return model.UserClaims{
			ID:   u.ID,
			Name: u.Username.String,
		}, nil
	}

	log.Info("oauthDto.Code: %s", oauthDto.Code)
	log.Info("oauthDto.State: %s", oauthDto.State)
	return nil, err
}