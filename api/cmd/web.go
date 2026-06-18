package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"

	"ultrathreads/bus"
	"ultrathreads/cache"
	"ultrathreads/cron"
	"ultrathreads/dao"
	"ultrathreads/database"
	"ultrathreads/handler"
	"ultrathreads/server"
	"ultrathreads/service"
	"ultrathreads/util/hashid"
)

var CmdWeb = &cli.Command{
	Name:   "web",
	Usage:  "Start UltraThreads API",
	Action: runWeb,
}

func runWeb(c *cli.Context) error {
	// 1. 初始化日志
	zerolog.SetGlobalLevel(zerolog.Level(0))

	// 2. 加载配置
	conf := "./app.yaml"
	if c.IsSet("conf") {
		conf = c.String("conf")
	}
	viper.SetConfigFile(conf)

	content, err := os.ReadFile(conf)
	if err != nil {
		return fmt.Errorf("read conf file fail: %w", err)
	}
	if err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content)))); err != nil {
		return fmt.Errorf("parse conf file fail: %w", err)
	}

	// 3. 初始化工具组件
	hashid.Init("ultrathreads", 8)

	// 4. 初始化事件总线
	mgr := bus.NewManager()
	bus.RegisterHandlers(mgr)

	// 5. 设置 Gin 运行模式
	gin.SetMode(viper.GetString("mode"))

	// ========== 核心依赖链构建 ==========
	db, err := database.Setup()
	if err != nil {
		return fmt.Errorf("setup database failed: %w", err)
	}

	repos := dao.NewRepositories(db)
	caches := cache.NewCaches(repos)

	svcs := service.NewServices(repos, caches)

	cron.Setup(svcs.Article, svcs.Post)

	handlers := handler.NewHandlers(svcs, mgr)

	// ========== Web Server ==========
	srv := server.NewServer(server.Config{
		Port:            viper.GetString("base.port"),
		ReadTimeout:     viper.GetDuration("http.read_timeout"),
		WriteTimeout:    viper.GetDuration("http.write_timeout"),
		ShutdownTimeout: time.Duration(viper.GetInt("shutdown_timeout")) * time.Second,
		MaxHeaderBytes:  viper.GetInt("http.max_header_megabytes") << 20,
	}, handlers.Init())

	errCh := srv.Start()

	// ========== 优雅退出 ==========
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		fmt.Printf("\n⏳ Received signal [%s], shutting down gracefully...\n", sig)
	case err := <-errCh:
		return err
	}

	if err := srv.Stop(context.Background()); err != nil {
		fmt.Printf("❌ Server shutdown error: %v\n", err)
	}

	// ========== 按依赖逆序关闭资源 ==========
	cron.Stop()
	fmt.Println("✅ Cron jobs stopped")

	cache.Shutdown()
	fmt.Println("✅ Cache closed")

	if err := database.Close(db); err != nil {
		fmt.Printf("❌ Database close failed: %v\n", err)
	} else {
		fmt.Println("✅ Database closed")
	}

	fmt.Println("👋 UltraThreads exited cleanly")
	return nil
}
