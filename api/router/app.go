package router

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/controller"
	"ultrathreads/middleware"
)

func setupApp(e *gin.Engine) {
	api := e.Group("/api")

	// ---------- 公开接口（无需登录，也不需要 OptionalAuth） ----------
	api.Any("/stat", new(controller.SiteController).Stat)
	api.Any("/ping", new(controller.SiteController).Ping)

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

	// Captcha
	captchaController := &controller.CaptchaController{}
	api.GET("/captcha/request", captchaController.GetRequest)
	api.GET("/captcha/show/:captchaId", captchaController.Show)

	// ---------- Optional Auth 组（未登录可访问，已登录自动注入用户上下文与已读状态） ----------
	optional := api.Group("/")
	optional.Use(middleware.OptionalAuth(jwtAuth))
	optional.Use(middleware.CurrentUserReadState())
	{
		postController := &controller.PostController{}
		nodeController := &controller.NodeController{}
		tagController := &controller.TagController{}
		articleController := &controller.ArticleController{}
		userController := &controller.UserController{}
		linkController := &controller.LinkController{}

		// Nodes
		optional.GET("/nodes", nodeController.List)
		optional.GET("/node/:slug", nodeController.Show)

		optional.GET("/threads", postController.ListThreads)
		nodeGroup := optional.Group("/nodes/:slug")
		{
			nodeGroup.GET("/threads", postController.ListThreads)
		}

		threadGroup := optional.Group("/threads/:slug")
		{
			threadGroup.GET("/with-thread", postController.GetPostWithThread)
		}

		optional.GET("/posts", postController.List)
		optional.GET("/post/:slug", postController.Show)
		optional.GET("/post/:slug/with-thread", postController.GetPostWithThread)
		optional.GET("/post/:slug/flat", postController.GetPostsFlat)
		optional.GET("/posts/excellent", postController.GetPostsExcellent)
		optional.GET("/posts/recommend", postController.GetPostsRecommend)
		optional.GET("/posts/noreply", postController.GetPostsNoreply)
		optional.GET("/posts/last", postController.GetPostsLast)
		optional.GET("/posts/user/recent/:id", postController.GetUserRecent)
		optional.GET("/user/posts/:slug", postController.GetUserPosts)
		optional.GET("/post/:slug/recentlikes", postController.GetRecentLikes)

		// Tags
		optional.GET("/tag/:slug", tagController.Show)
		optional.GET("/tags", tagController.List)
		optional.GET("/tags/hot", tagController.HotTags)

		tagGroup := optional.Group("/tags/:slug")
		{
			 tagGroup.GET("/threads", postController.ListTagThreads)
		}

		// Articles
		optional.GET("/articles", articleController.List)
		optional.GET("/article/:id", articleController.Show)
		optional.GET("/articles/related/:id", articleController.GetRelatedBy)
		optional.GET("/articles/tag/:slug", articleController.GetTagArticles)
		optional.GET("/articles/user/newest/:id", articleController.GetUserNewestBy)
		optional.GET("/articles/recent", articleController.GetRecent)
		optional.GET("/articles/user/recent/:id", articleController.GetUserRecent)
		optional.GET("/user/articles/:id", articleController.GetUserArticles)

		// Users（公开资料）
		optional.GET("/profile/:slug", userController.Show)
		optional.GET("/user/score/rank", userController.GetScoreRank)
		optional.GET("/users/:id/recentwatchers", userController.GetRecentWatchers)

		// Links
		optional.GET("/links/top", linkController.GetToplinks)
		optional.GET("/links", linkController.List)
	}

	// ---------- 需登录接口（强制鉴权） ----------
	jwtApi := api.Group("/")
	jwtApi.Use(jwtAuth.MiddlewareFunc(), middleware.CurrentUser)
	{
		nodeController := &controller.NodeController{}
		postController := &controller.PostController{}
		favoriteController := &controller.FavoriteController{}
		tagController := &controller.TagController{}
		articleController := &controller.ArticleController{}
		userController := &controller.UserController{}
		uploadController := &controller.UploadController{}

		// Nodes
		jwtApi.POST("/nodes/:slug/mark-as-read", nodeController.MarkAsRead)

		// Posts（写操作）
		jwtApi.POST("/posts", postController.Store)
		jwtApi.GET("/post/:slug/edit", postController.Edit)
		jwtApi.PUT("/post/:slug", postController.Update)
		jwtApi.Any("/post/:slug/view-post", postController.ViewPost)
		jwtApi.POST("/post/:slug/like", postController.Like)
		jwtApi.POST("/post/:slug/favorite", postController.Favorite)

		// Favorites
		jwtApi.GET("/favorites/favorited", favoriteController.GetFavorited)
		jwtApi.DELETE("/favorite/delete", favoriteController.Delete)

		// Tags
		jwtApi.POST("/tag/auto-complete", tagController.AutoComplete)

		// Articles（写操作）
		jwtApi.POST("/articles", articleController.Store)
		jwtApi.GET("/article/:id/edit", articleController.Edit)
		jwtApi.PUT("/article/:id", articleController.Update)
		jwtApi.POST("/article/:id/favorite", articleController.Favorite)

		// Users（个人操作）
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
		jwtApi.POST("/upload", uploadController.Upload)
		jwtApi.POST("/upload/editor", uploadController.UploadFromEditor)
		jwtApi.POST("/upload/fetch", uploadController.UploadFromURL)
	}
}