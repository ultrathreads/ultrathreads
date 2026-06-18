package handler

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/handler/app"
	"ultrathreads/middleware"
)

// initAppAPI 注册所有前台路由，JWT 通过参数注入而非全局变量
func (h *Handler) initAppAPI(e *gin.Engine) {
	api := e.Group("/api")
	svc := h.services

	// --- 创建所有 Handler 实例 ---
	siteHandler := app.NewSiteHandler(svc.Setting, svc.Appinfo, svc.UserReadState)
	authHandler := app.NewAuthHandler(svc.User)
	oauthHandler := &app.OAuthHandler{}
	captchaHandler := &app.CaptchaHandler{}
	postHandler := app.NewPostHandler(svc.Post, svc.User, svc.PostLike, svc.Favorite)
	nodeHandler := app.NewNodeHandler(svc.Node, svc.UserReadState)
	tagHandler := app.NewTagHandler(svc.Tag)
	articleHandler := app.NewArticleHandler(svc.Article, svc.Favorite)
	userHandler := app.NewUserHandler(svc.User, svc.Post, svc.UserScore, svc.UserScoreLog, svc.Notification, svc.Favorite, svc.Article, svc.UserWatch, svc.Rbac)
	linkHandler := app.NewLinkHandler(svc.Link)
	favoriteHandler := app.NewFavoriteHandler(svc.Favorite)
	uploadHandler := &app.UploadHandler{}

	// ---------- 公开接口（无需登录，也不需要 OptionalAuth） ----------
	api.Any("/stat", siteHandler.Stat)
	api.Any("/ping", siteHandler.Ping)
	api.Any("/site/config", siteHandler.Config)
	api.Any("/debug", siteHandler.Debug)

	// Auth
	api.POST("/auth/login", h.jwtAuth.LoginHandler)
	api.POST("/auth/login/refresh", h.jwtAuth.RefreshHandler)
	api.POST("/auth/register", authHandler.Register)

	// OAuth
	api.GET("/oauth/:provider/authorize", oauthHandler.Authorize)
	api.GET("/oauth/:provider/callback", h.jwtOAuth.LoginHandler)

	// Captcha
	api.GET("/captcha/request", captchaHandler.GetRequest)
	api.GET("/captcha/show/:captchaId", captchaHandler.Show)

	// ---------- Optional Auth 组 ----------
	optional := api.Group("/")
	optional.Use(middleware.OptionalAuth(h.jwtAuth))
	optional.Use(middleware.CurrentUserReadState(svc.UserReadState))
	{
		// Home
		optional.GET("/threads", postHandler.ListThreads)

		// Nodes
		nodeGroup := optional.Group("/nodes")
		{
			nodeGroup.GET("", nodeHandler.List)
			nodeGroup.GET("/:slug", nodeHandler.Show)
			nodeGroup.GET("/:slug/threads", postHandler.ListThreads)
		}

		// Posts
		postApi := optional.Group("/posts")
		{
			postApi.GET("/:slug", postHandler.Show)
			postApi.GET("/:slug/tree", postHandler.GetPostTree)
			postApi.GET("/:slug/flat", postHandler.GetPostFlat)
		}
		optional.GET("/posts/user/recent/:id", postHandler.GetUserRecent)
		optional.GET("/user/posts/:slug", postHandler.GetUserPosts)

		// Tags
		tagGroup := optional.Group("/tags")
		{
			tagGroup.GET("", tagHandler.List)
			tagGroup.GET("/hot", tagHandler.HotTags)
			tagGroup.GET("/:slug", tagHandler.Show)
			tagGroup.GET("/:slug/threads", postHandler.ListTagThreads)
		}

		// Users (public profile)
		userGroup := optional.Group("/users/:slug")
		{
			userGroup.GET("/posts", postHandler.GetUserPosts)
		}

		// Articles
		optional.GET("/articles", articleHandler.List)
		optional.GET("/article/:id", articleHandler.Show)
		optional.GET("/articles/related/:id", articleHandler.GetRelatedBy)
		optional.GET("/articles/tag/:slug", articleHandler.GetTagArticles)
		optional.GET("/articles/user/newest/:id", articleHandler.GetUserNewestBy)
		optional.GET("/articles/recent", articleHandler.GetRecent)
		optional.GET("/articles/user/recent/:id", articleHandler.GetUserRecent)
		optional.GET("/user/articles/:id", articleHandler.GetUserArticles)

		// User score & profile
		optional.GET("/profile/:slug", userHandler.Show)
		optional.GET("/user/score/rank", userHandler.GetScoreRank)

		// Links
		optional.GET("/links/top", linkHandler.GetToplinks)
		optional.GET("/links", linkHandler.List)
	}

	// ---------- 需登录接口（强制鉴权） ----------
	jwtApi := api.Group("/")
	jwtApi.Use(h.jwtAuth.MiddlewareFunc(), middleware.CurrentUser)
	{
		// Nodes
		jwtApi.POST("/nodes/:slug/mark-as-read", nodeHandler.MarkAsRead)

		// Posts (write)
		postGroup := jwtApi.Group("/posts")
		{
			postGroup.POST("", postHandler.StoreRootPost)
			postGroup.POST("/:slug/replies", postHandler.StoreReply)
			postGroup.POST("/:slug", postHandler.UpdateRootPost)
			postGroup.POST("/:slug/like", postHandler.Like)
			postGroup.POST("/:slug/favorite", postHandler.Favorite)
			postGroup.Any("/:slug/view-post", postHandler.ViewPost)
		}
		jwtApi.POST("/replies/:slug", postHandler.UpdateReply)

		// Favorites
		jwtApi.GET("/favorites/favorited", favoriteHandler.GetFavorited)
		jwtApi.DELETE("/favorite/delete", favoriteHandler.Delete)

		// Tags
		jwtApi.POST("/tags/auto-complete", tagHandler.AutoComplete)

		// Articles (write)
		jwtApi.POST("/articles", articleHandler.Store)
		jwtApi.GET("/article/:id/edit", articleHandler.Edit)
		jwtApi.PUT("/article/:id", articleHandler.Update)
		jwtApi.POST("/article/:id/favorite", articleHandler.Favorite)

		// Users (personal)
		jwtApi.PUT("/users/:slug", userHandler.Update)
		jwtApi.GET("/user/current", userHandler.GetCurrent)
		jwtApi.GET("/user/scorelogs", userHandler.GetScorelogs)
		jwtApi.GET("/user/notifications/recent", userHandler.GetNotificationsRecent)
		jwtApi.GET("/user/notifications", userHandler.GetNotifications)
		jwtApi.GET("/user/favorites", userHandler.GetFavorites)
		jwtApi.PUT("/user/update/avatar", userHandler.UpdateAvatar)
		jwtApi.PUT("/user/set/username", userHandler.SetUsername)
		jwtApi.PUT("/user/set/email", userHandler.SetEmail)
		jwtApi.PUT("/user/set/password", userHandler.SetPassword)
		jwtApi.PUT("/user/change/password", userHandler.ChangePassword)
		jwtApi.POST("/users/:id/watch", userHandler.Watch)
		jwtApi.GET("/watch/watched", userHandler.GetWatched)
		jwtApi.DELETE("/watch/delete", userHandler.WatchDelete)

		// Upload
		jwtApi.POST("/upload", uploadHandler.Upload)
		jwtApi.POST("/upload/editor", uploadHandler.UploadFromEditor)
		jwtApi.POST("/upload/fetch", uploadHandler.UploadFromURL)
	}
}
