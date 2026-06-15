package router

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/controller"
	"ultrathreads/middleware"
)

func setupApp(e *gin.Engine) {
	api := e.Group("/api")

	// ---------- 公开接口（无需登录，也不需要 OptionalAuth） ----------
	siteController := &controller.SiteController{}
	api.Any("/stat", siteController.Stat)
	api.Any("/ping", siteController.Ping)
	api.Any("/site/config", siteController.Config)
	api.Any("/debug", siteController.Debug)

	// Auth
	api.POST("/auth/login", jwtAuth.LoginHandler)
	api.POST("/auth/login/refresh", jwtAuth.RefreshHandler)
	api.POST("/auth/register", new(controller.AuthController).Register)

	// OAuth
	oauthController := &controller.OAuthController{}
	api.GET("/oauth/:provider/authorize", oauthController.Authorize)
	api.GET("/oauth/:provider/callback", jwtOAuth.LoginHandler)

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

		// Home
		optional.GET("/threads", postController.ListThreads)

		optional.GET("/sideload", postController.ListThreads)

		// Nodes
		nodeGroup := optional.Group("/nodes")
		{
			nodeGroup.GET("", nodeController.List)
			nodeGroup.GET("/:slug", nodeController.Show)
			nodeGroup.GET("/:slug/threads", postController.ListThreads)
		}

		postApi := optional.Group("/posts")
		{
			postApi.GET("", postController.List)
			postApi.GET("/:slug", postController.Show)
			postApi.GET("/:slug/tree", postController.GetPostTree)
			postApi.GET("/:slug/flat", postController.GetPostFlat)
		}
		
		optional.GET("/posts/excellent", postController.GetPostsExcellent)
		optional.GET("/posts/recommend", postController.GetPostsRecommend)
		optional.GET("/posts/noreply", postController.GetPostsNoreply)
		optional.GET("/posts/last", postController.GetPostsLast)
		optional.GET("/posts/user/recent/:id", postController.GetUserRecent)
		optional.GET("/user/posts/:slug", postController.GetUserPosts)
		optional.GET("/post/:slug/recentlikes", postController.GetRecentLikes)

		// Tags
		tagGroup := optional.Group("/tags")
		{
			tagGroup.GET("", tagController.List)
			tagGroup.GET("/hot", tagController.HotTags)
			tagGroup.GET("/:slug", tagController.Show)
			tagGroup.GET("/:slug/threads", postController.ListTagThreads)
		}


		userGroup := optional.Group("/users/:slug")
		{
			userGroup.GET("/posts", postController.GetUserPosts)
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
		optional.GET("/users/:slug/recentwatchers", userController.GetRecentWatchers)

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
		postGroup := jwtApi.Group("/posts")
		{
			postGroup.POST("", postController.StoreRootPost)
			postGroup.POST("/:slug/replies", postController.StoreReply)
			postGroup.POST("/:slug", postController.UpdateRootPost)
			postGroup.POST("/:slug/like", postController.Like)
			postGroup.POST("/:slug/favorite", postController.Favorite)
			postGroup.Any("/:slug/view-post", postController.ViewPost)

		}

		jwtApi.POST("/replies/:slug", postController.UpdateReply)

		// Favorites
		jwtApi.GET("/favorites/favorited", favoriteController.GetFavorited)
		jwtApi.DELETE("/favorite/delete", favoriteController.Delete)

		// Tags
		jwtApi.POST("/tags/auto-complete", tagController.AutoComplete)

		// Articles（写操作）
		jwtApi.POST("/articles", articleController.Store)
		jwtApi.GET("/article/:id/edit", articleController.Edit)
		jwtApi.PUT("/article/:id", articleController.Update)
		jwtApi.POST("/article/:id/favorite", articleController.Favorite)

		// Users（个人操作）
		jwtApi.PUT("/users/:slug", userController.Update)
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