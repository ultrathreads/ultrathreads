package admin

import (
	"runtime"
	"time"

	"github.com/gin-gonic/gin"

	"ultrathreads/config"
	"ultrathreads/delivery/handler/base"
	"ultrathreads/util"
)

// initTime is the time when the application was initialized.
var initTime = time.Now()

// DashboardHandler dashboard controller
type DashboardHandler struct {
	base.BaseHandler
}

// GetSysteminfo get system info
func (h *DashboardHandler) Systeminfo(ctx *gin.Context) {
	h.Success(ctx, gin.H{
		"registerUserCount": 100,
		"postTotalCount":    200,
		"todayNewPostCount": 30,
		"appName":           config.AppName,
		"appVersion":        config.AppVersion,
		"buildTime":         config.BuildTime,
		"buildCommit":       config.BuildCommit,
		"upTime":            util.TimeSincePro(initTime),
		"os":                runtime.GOOS,
		"arch":              runtime.GOARCH,
		"numCpu":            runtime.NumCPU(),
		"goversion":         runtime.Version(),
	})
}
