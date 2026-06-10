package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"ultrathreads/cache"
	"ultrathreads/cron"
	"ultrathreads/dao"
	"ultrathreads/middleware"
	"ultrathreads/router"
)

var CmdWeb = &cli.Command{
	Name:   "web",
	Usage:  "Start UltraThreads API",
	Action: runWeb,
	Flags:  []cli.Flag{}, // ✅ 保持为空，conf 由 App 全局继承
}

func runWeb(c *cli.Context) error {
	// 1. Set up log level
	zerolog.SetGlobalLevel(zerolog.Level(0))

	conf := "./app.yaml"
	if c.IsSet("conf") {
		conf = c.String("conf")
	}

	// 2. Set up configuration
	viper.SetConfigFile(conf)

	content, err := os.ReadFile(conf)
	if err != nil {
		return fmt.Errorf("read conf file fail: %w", err)
	}

	// Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		return fmt.Errorf("parse conf file fail: %w", err)
	}

	// 3. Set up run mode
	mode := viper.GetString("mode")
	gin.SetMode(mode)

	// 4. Set up database connection
	dao.Setup()

	// 5. Set up cache
	cache.Setup()

	// 6. Set up cron
	cron.Setup()

	// 7. Initialize language
	middleware.InitLang()

	engine := gin.Default()
	router.Setup(engine)

	port := viper.GetString("base.port")
	addr := ":" + port

	// ========== 优雅退出核心逻辑 ==========
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	// 在 goroutine 中启动服务，避免阻塞主线程
	go func() {
		fmt.Printf("🚀 UltraThreads starting on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("❌ Server listen failed: %v\n", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号 (SIGINT / SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	fmt.Printf("\n⏳ Received signal [%s], shutting down gracefully...\n", sig)

	// 设置超时上下文，防止关闭过程无限等待
	shutdownTimeout := time.Duration(viper.GetInt("shutdown_timeout")) * time.Second
	if shutdownTimeout == 0 {
		shutdownTimeout = 10 * time.Second // 默认 10 秒
	}
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// 1) 停止接收新请求，等待已有请求处理完成
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("❌ HTTP server forced to shutdown: %v\n", err)
	} else {
		fmt.Println("✅ HTTP server stopped gracefully")
	}

	// 2) 关闭定时任务
	cron.Stop()
	fmt.Println("✅ Cron jobs stopped")

	// 3) 关闭缓存连接
	cache.Shutdown()
	fmt.Println("✅ Cache closed")

	// 4) 关闭数据库连接
	dao.Shutdown()
	fmt.Println("✅ Database closed")

	fmt.Println("👋 UltraThreads exited cleanly")
	return nil
}