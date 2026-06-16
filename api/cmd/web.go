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

	"ultrathreads/bus"
	"ultrathreads/cache"
	"ultrathreads/cron"
	"ultrathreads/dao"
	"ultrathreads/database"
	"ultrathreads/service"
	"ultrathreads/router"
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

	// 5. 设置 Gin 模式
	gin.SetMode(viper.GetString("mode"))

	// ========== 🆕 核心改动：显式依赖注入与资源获取 ==========
	
	// 6. 初始化数据库（返回实例，不再使用全局变量）
	db, err := database.Setup()
	if err != nil {
		return fmt.Errorf("setup database failed: %w", err)
	}

	// 7. 初始化 DAO 聚合体（注入 db 实例）
	dao.Setup(db)

	// 8. 🆕 初始化 Service 并赋值给全局 Srv
	// 过渡期暂不传参，内部仍读取 dao.XxxDao 全局变量
	service.Srv = service.NewServices(dao.NewDaos(db))

	// 8. 初始化缓存与定时任务
	cache.Setup()
	cron.Setup()

	// 9. 路由注册
	engine := gin.Default()
	router.Setup(engine, mgr, service.Srv) // 🔄 router.Setup 签名需同步调整以接收 daos

	port := viper.GetString("base.port")
	addr := ":" + port

	// ========== 优雅退出逻辑（保持原有优秀设计） ==========
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	go func() {
		fmt.Printf("🚀 UltraThreads starting on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("❌ Server listen failed: %v\n", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	fmt.Printf("\n⏳ Received signal [%s], shutting down gracefully...\n", sig)

	shutdownTimeout := time.Duration(viper.GetInt("shutdown_timeout")) * time.Second
	if shutdownTimeout == 0 {
		shutdownTimeout = 10 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("❌ HTTP server forced to shutdown: %v\n", err)
	} else {
		fmt.Println("✅ HTTP server stopped gracefully")
	}

	// ========== 🆕 按依赖逆序关闭资源 ==========
	cron.Stop()
	fmt.Println("✅ Cron jobs stopped")

	cache.Shutdown()
	fmt.Println("✅ Cache closed")

	// 使用新的 database.Close 并传入 db 实例
	if err := database.Close(db); err != nil {
		fmt.Printf("❌ Database close failed: %v\n", err)
	} else {
		fmt.Println("✅ Database closed")
	}

	fmt.Println("👋 UltraThreads exited cleanly")
	return nil
}