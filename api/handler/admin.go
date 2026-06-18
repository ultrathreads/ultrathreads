package handler

import (
	"github.com/gin-gonic/gin"

	"ultrathreads/controller/admin"
	"ultrathreads/middleware"
)

func (h *Handler) initAdminAPI(e *gin.Engine) {
	adminAPI := e.Group("/api/admin")
	adminAPI.Use(h.jwtAuth.MiddlewareFunc(), middleware.CurrentUser, middleware.AdminRequired(h.services.Rbac))
	svc := h.services

	dashboardCtrl := &admin.DashboardController{}
	adminAPI.GET("/dashboard/systeminfo", dashboardCtrl.Systeminfo)

	nodeCtrl := admin.NewNodeController(svc.Node)
	adminAPI.GET("/nodes", nodeCtrl.List)
	adminAPI.GET("/nodes/:id", nodeCtrl.Show)
	adminAPI.POST("/nodes", nodeCtrl.Store)
	adminAPI.PUT("/nodes/sort", nodeCtrl.Sort)
	adminAPI.PUT("/nodes/:id", nodeCtrl.Update)
	adminAPI.DELETE("/nodes/:id", nodeCtrl.Delete)

	postCtrl := admin.NewPostController(svc.Post, svc.Node)
	adminAPI.GET("/posts", postCtrl.List)
	adminAPI.GET("/posts/:id", postCtrl.Show)
	adminAPI.PUT("/posts/:id", postCtrl.Update)
	adminAPI.DELETE("/posts/:id", postCtrl.Delete)
	adminAPI.POST("/posts/:id/recommend", postCtrl.Recommend)
	adminAPI.POST("/posts/:id/unrecommend", postCtrl.Unrecommend)
	adminAPI.POST("/posts/:id/undelete", postCtrl.Undelete)

	tagCtrl := admin.NewTagController(svc.Tag)
	adminAPI.GET("/tags", tagCtrl.List)
	adminAPI.GET("/tags/:id", tagCtrl.Show)
	adminAPI.PUT("/tags/:id", tagCtrl.Update)
	adminAPI.DELETE("/tags/:id", tagCtrl.Delete)

	articleCtrl := admin.NewArticleController(svc.Article)
	adminAPI.GET("/articles", articleCtrl.List)
	adminAPI.GET("/articles/:id", articleCtrl.Show)
	adminAPI.PUT("/articles/:id", articleCtrl.Update)
	adminAPI.DELETE("/articles/:id", articleCtrl.Delete)

	userCtrl := admin.NewUserController(svc.User)
	adminAPI.GET("/users", userCtrl.List)
	adminAPI.GET("/users/:id", userCtrl.Show)
	adminAPI.POST("/users", userCtrl.Store)
	adminAPI.PUT("/users/:id", userCtrl.Update)
	adminAPI.DELETE("/users/:id", userCtrl.Delete)

	scoreCtrl := admin.NewUserScoreController(svc.UserScore)
	adminAPI.GET("/user-scores", scoreCtrl.List)
	adminAPI.GET("/user-scores/:id", scoreCtrl.Show)

	scoreLogCtrl := admin.NewUserScoreLogController(svc.UserScoreLog)
	adminAPI.GET("/user-score-logs", scoreLogCtrl.List)
	adminAPI.GET("/user-score-logs/:id", scoreLogCtrl.Show)

	linkCtrl := admin.NewLinkController(svc.Link)
	adminAPI.GET("/links", linkCtrl.List)
	adminAPI.GET("/links/:id", linkCtrl.Show)
	adminAPI.POST("/links", linkCtrl.Store)
	adminAPI.PUT("/links/:id", linkCtrl.Update)
	adminAPI.DELETE("/links/:id", linkCtrl.Delete)

	settingCtrl := admin.NewSettingController(svc.Setting)
	adminAPI.GET("/settings", settingCtrl.List)
	adminAPI.POST("/settings", settingCtrl.Store)
}
