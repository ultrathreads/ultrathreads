package router

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/controller/admin"
	"ultrathreads/middleware"
)

func setupAdmin(e *gin.Engine) {
	adminAPI := e.Group("/api/admin")
	adminAPI.Use(jwtAuth.MiddlewareFunc(), middleware.CurrentUser, middleware.AdminRequired())

	// Dashboard
	dashboardController := &admin.DashboardController{}
	adminAPI.GET("/dashboard/systeminfo", dashboardController.Systeminfo)

	// Node
	adminNodeController := &admin.NodeController{}
	adminAPI.GET("/nodes", adminNodeController.List)
	adminAPI.GET("/nodes/:id", adminNodeController.Show)
	adminAPI.POST("/nodes", adminNodeController.Store)
	adminAPI.PUT("/nodes/sort", adminNodeController.Sort)  // ✅ 固定路径优先
	adminAPI.PUT("/nodes/:id", adminNodeController.Update)
	adminAPI.DELETE("/nodes/:id", adminNodeController.Delete)

	// Post
	adminPostController := &admin.PostController{}
	adminAPI.GET("/posts", adminPostController.List)
	adminAPI.GET("/posts/:id", adminPostController.Show)
	adminAPI.PUT("/posts/:id", adminPostController.Update)
	adminAPI.DELETE("/posts/:id", adminPostController.Delete)
	adminAPI.POST("/posts/:id/recommend", adminPostController.Recommend)
	adminAPI.POST("/posts/:id/unrecommend", adminPostController.Unrecommend)
	adminAPI.POST("/posts/:id/undelete", adminPostController.Undelete)

	// Tag
	adminTagController := &admin.TagController{}
	adminAPI.GET("/tags", adminTagController.List)
	adminAPI.GET("/tags/:id", adminTagController.Show)
	adminAPI.PUT("/tags/:id", adminTagController.Update)
	adminAPI.DELETE("/tags/:id", adminTagController.Delete)

	// Article
	adminArticleController := &admin.ArticleController{}
	adminAPI.GET("/articles", adminArticleController.List)
	adminAPI.GET("/articles/:id", adminArticleController.Show)
	adminAPI.PUT("/articles/:id", adminArticleController.Update)
	adminAPI.DELETE("/articles/:id", adminArticleController.Delete)

	// User
	adminUserController := &admin.UserController{}
	adminAPI.GET("/users", adminUserController.List)
	adminAPI.GET("/users/:id", adminUserController.Show)
	adminAPI.POST("/users", adminUserController.Store)
	adminAPI.PUT("/users/:id", adminUserController.Update)
	adminAPI.DELETE("/users/:id", adminUserController.Delete)

	// User Score
	adminUserScoreController := &admin.UserScoreController{}
	adminAPI.GET("/user-scores", adminUserScoreController.List)
	adminAPI.GET("/user-scores/:id", adminUserScoreController.Show)

	// User Score Log
	adminUserScoreLogController := &admin.UserScoreLogController{}
	adminAPI.GET("/user-score-logs", adminUserScoreLogController.List)
	adminAPI.GET("/user-score-logs/:id", adminUserScoreLogController.Show)

	// Link
	adminLinkController := &admin.LinkController{}
	adminAPI.GET("/links", adminLinkController.List)
	adminAPI.GET("/links/:id", adminLinkController.Show)
	adminAPI.POST("/links", adminLinkController.Store)
	adminAPI.PUT("/links/:id", adminLinkController.Update)
	adminAPI.DELETE("/links/:id", adminLinkController.Delete)

	// Settings
	adminSettingController := &admin.SettingController{}
	adminAPI.GET("/settings", adminSettingController.List)
	adminAPI.POST("/settings", adminSettingController.Store)
}