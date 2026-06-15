package admin

import (
	"github.com/gin-gonic/gin"
	"runtime"
	"time"

	"ultrathreads/config"
	"ultrathreads/controller"
	"ultrathreads/util"
)

// initTime is the time when the application was initialized.
var initTime = time.Now()

// DashboardController dashboard controller
type DashboardController struct {
	controller.BaseController
}

// GetSysteminfo get system info
func (c *DashboardController) Systeminfo(ctx *gin.Context) {
	c.Success(ctx, gin.H{
		"registerUserCount": 100,
		"postTotalCount":    200,
		"todayNewPostCount": 30,
		"appName":     		 config.AppName,
		"appVersion":  		 config.AppVersion,
		"buildTime":   		 config.BuildTime,
		"buildCommit": 		 config.BuildCommit,
		"upTime":      		 util.TimeSincePro(initTime),
		"os":          		 runtime.GOOS,
		"arch":        		 runtime.GOARCH,
		"numCpu":      		 runtime.NumCPU(),
		"goversion":         runtime.Version(),
	})
}
