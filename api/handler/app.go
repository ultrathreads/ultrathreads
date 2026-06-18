package handler

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/controller"
	"ultrathreads/middleware"
)

// initAppAPI 注册所有前台路由，JWT 通过参数注入而非全局变量
func (h *Handler) initAppAPI(e *gin.Engine) {
	api := e.Group("/api")
	svc := h.services

	// --- 创建所有 Controller 实例 ---
	siteController := controller.NewSiteController(svc.Setting, svc.Appinfo, svc.UserReadState)
	authController := controller.NewAuthController(svc.User)
	oauthController := &controller.OAuthController{}
	captchaController := &controller.CaptchaController{}
	postController := controller.NewPostController(svc.Post, svc.User, svc.PostLike, svc.Favorite)
	nodeController := controller.NewNodeController(svc.Node, svc.UserReadState)
	tagController := controller.NewTagController(svc.Tag)
	articleController := controller.NewArticleController(svc.Article, svc.Favorite)
	userController := controller.NewUserController(svc.User, svc.Post, svc.UserScore, svc.UserScoreLog, svc.Notification, svc.Favorite, svc.Article, svc.UserWatch, svc.Rbac)
	linkController := controller.NewLinkController(svc.Link)
	favoriteController := controller.NewFavoriteController(svc.Favorite)
	uploadController := &controller.UploadController{}

	// ---------- 公开接口（无需登录，也不需要 OptionalAuth） ----------
	api.Any("/stat", siteController.Stat)
	api.Any("/ping", siteController.Ping)
	api.Any("/site/config", siteController.Config)
	api.Any("/debug", siteController.Debug)

	// Auth
	api.POST("/auth/login", h.jwtAuth.LoginHandler)
	api.POST("/auth/login/refresh", h.jwtAuth.RefreshHandler)
	api.POST("/auth/register", authController.Register)

	// OAuth
	api.GET("/oauth/:provider/authorize", oauthController.Authorize)
	api.GET("/oauth/:provider/callback", h.jwtOAuth.LoginHandler)

	// Captcha
	api.GET("/captcha/request", captchaController.GetRequest)
	api.GET("/captcha/show/:captchaId", captchaController.Show)

	// ---------- Optional Auth 组 ----------
	optional := api.Group("/")
	optional.Use(middleware.OptionalAuth(h.jwtAuth))
	optional.Use(middleware.CurrentUserReadState(svc.UserReadState))
	{
		// Home
		optional.GET("/threads", postController.ListThreads)

		// Nodes
		nodeGroup := optional.Group("/nodes")
		{
			nodeGroup.GET("", nodeController.List)
			nodeGroup.GET("/:slug", nodeController.Show)
			nodeGroup.GET("/:slug/threads", postController.ListThreads)
		}

		// Posts
		postApi := optional.Group("/posts")
		{
			postApi.GET("/:slug", postController.Show)
			postApi.GET("/:slug/tree", postController.GetPostTree)
			postApi.GET("/:slug/flat", postController.GetPostFlat)
		}
		optional.GET("/posts/user/recent/:id", postController.GetUserRecent)
		optional.GET("/user/posts/:slug", postController.GetUserPosts)

		// Tags
		tagGroup := optional.Group("/tags")
		{
			tagGroup.GET("", tagController.List)
			tagGroup.GET("/hot", tagController.HotTags)
			tagGroup.GET("/:slug", tagController.Show)
			tagGroup.GET("/:slug/threads", postController.ListTagThreads)
		}

		// Users (public profile)
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

		// User score & profile
		optional.GET("/profile/:slug", userController.Show)
		optional.GET("/user/score/rank", userController.GetScoreRank)

		// Links
		optional.GET("/links/top", linkController.GetToplinks)
		optional.GET("/links", linkController.List)
	}

	// ---------- 需登录接口（强制鉴权） ----------
	jwtApi := api.Group("/")
	jwtApi.Use(h.jwtAuth.MiddlewareFunc(), middleware.CurrentUser)
	{
		// Nodes
		jwtApi.POST("/nodes/:slug/mark-as-read", nodeController.MarkAsRead)

		// Posts (write)
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

		// Articles (write)
		jwtApi.POST("/articles", articleController.Store)
		jwtApi.GET("/article/:id/edit", articleController.Edit)
		jwtApi.PUT("/article/:id", articleController.Update)
		jwtApi.POST("/article/:id/favorite", articleController.Favorite)

		// Users (personal)
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
