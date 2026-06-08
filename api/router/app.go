package router

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/controller"
	"ultrathreads/middleware"
)

func setupApp(e *gin.Engine) {
	api := e.Group("/api")

	// ---------- 公开接口（无需登录） ----------
	api.Any("/stat", new(controller.SiteController).Stat)

	// Auth
	api.POST("/auth/login", jwtAuth.LoginHandler)
	api.POST("/auth/login/refresh", jwtAuth.RefreshHandler)
	api.POST("/auth/register", new(controller.AuthController).Register)

	// OAuth
	oauthController := &controller.OAuthController{}
	api.GET("/oauth/:provider/authorize", oauthController.Authorize)
	api.GET("/oauth/:provider/callback", jwtOAuth.LoginHandler)

	// Configs
	api.GET("/config/site-config", new(controller.ConfigController).List)

	// Nodes
	nodeController := &controller.NodeController{}
	api.GET("/nodes", nodeController.List)
	api.GET("/node/:id", nodeController.Show)

	// Posts（公开）
	postController := &controller.PostController{}
	api.GET("/posts", postController.List)
	api.GET("/posts/threads", postController.ListWithReplies)
	api.GET("/post/:id", postController.Show)
	api.GET("/post/:id/with-thread", postController.GetPostWithThread)
	api.GET("/post/:id/flat", postController.GetPostsFlat)
	api.GET("/posts/node", postController.GetNodePosts)
	api.GET("/posts/excellent", postController.GetPostsExcellent)
	api.GET("/posts/recommend", postController.GetPostsRecommend)
	api.GET("/posts/noreply", postController.GetPostsNoreply)
	api.GET("/posts/last", postController.GetPostsLast)
	api.GET("/posts/tag", postController.GetTagPosts)
	api.GET("/posts/user/recent/:id", postController.GetUserRecent)
	api.GET("/user/posts/:id", postController.GetUserPosts)
	api.GET("/post/:id/recentlikes", postController.GetRecentLikes)

	// Tags
	tagController := &controller.TagController{}
	api.GET("/tag/:id", tagController.Show)
	api.GET("/tags", tagController.List)
	api.GET("/tags/hot", tagController.HotTags)

	// Articles（公开）
	articleController := &controller.ArticleController{}
	api.GET("/articles", articleController.List)
	api.GET("/article/:id", articleController.Show)
	api.GET("/articles/related/:id", articleController.GetRelatedBy)
	api.GET("/articles/tag/:id", articleController.GetTagArticles)
	api.GET("/articles/user/newest/:id", articleController.GetUserNewestBy)
	api.GET("/articles/recent", articleController.GetRecent)
	api.GET("/articles/user/recent/:id", articleController.GetUserRecent)
	api.GET("/user/articles/:id", articleController.GetUserArticles)

	// Users（公开）
	userController := &controller.UserController{}
	api.GET("/profile/:id", userController.Show)
	api.GET("/user/score/rank", userController.GetScoreRank)
	api.GET("/users/:id/recentwatchers", userController.GetRecentWatchers)

	// Links
	linkController := &controller.LinkController{}
	api.GET("/links/top", linkController.GetToplinks)
	api.GET("/links", linkController.List)

	// Captcha
	captchaController := &controller.CaptchaController{}
	api.GET("/captcha/request", captchaController.GetRequest)
	api.GET("/captcha/show/:captchaId", captchaController.Show)

	// ---------- 需登录接口 ----------
	jwtApi := api.Group("/")
	jwtApi.Use(jwtAuth.MiddlewareFunc(), middleware.CurrentUser)

	// Posts（鉴权）
	jwtApi.POST("/posts", postController.Store)
	jwtApi.GET("/post/:id/edit", postController.Edit)
	jwtApi.PUT("/post/:id", postController.Update)
	jwtApi.POST("/post/:id/like", postController.Like)
	jwtApi.POST("/post/:id/favorite", postController.Favorite)

	// Favorites
	favoriteController := &controller.FavoriteController{}
	jwtApi.GET("/favorites/favorited", favoriteController.GetFavorited)
	jwtApi.DELETE("/favorite/delete", favoriteController.Delete)

	// Tags（鉴权）
	jwtApi.POST("/tag/auto-complete", tagController.AutoComplete)

	// Articles（鉴权）
	jwtApi.POST("/articles", articleController.Store)
	jwtApi.GET("/article/:id/edit", articleController.Edit)
	jwtApi.PUT("/article/:id", articleController.Update)
	jwtApi.POST("/article/:id/favorite", articleController.Favorite)

	// Users（鉴权）
	jwtApi.PUT("/users/:id", userController.Update)
	jwtApi.GET("/user/current", userController.GetCurrent)
	jwtApi.GET("/user/scorelogs", userController.GetScorelogs)
	jwtApi.GET("/user/notifications/recent", userController.GetNotificationsRecent)
	jwtApi.GET("/user/notifications", userController.GetNotifications)
	jwtApi.GET("/user/favorites", userController.GetFavorites)
	jwtApi.PUT("/user/update/avatar", userController.UpdateAvatar)
	jwtApi.PUT("/user/set/username", userController.SetUsername)
	jwtApi.PUT("/user/set/email", userController.SetEmail)
	jwtApi.PUT("/user/set/password", userController.SetPassword)
	jwtApi.PUT("/user/change/password", userController.ChangePassword)
	jwtApi.POST("/users/:id/watch", userController.Watch)
	jwtApi.GET("/watch/watched", userController.GetWatched)
	jwtApi.DELETE("/watch/delete", userController.WatchDelete)

	// Upload
	uploadController := &controller.UploadController{}
	jwtApi.POST("/upload", uploadController.Upload)
	jwtApi.POST("/upload/editor", uploadController.UploadFromEditor)
	jwtApi.POST("/upload/fetch", uploadController.UploadFromURL)
}