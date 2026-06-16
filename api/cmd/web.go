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
	// 1. 初始化日志（TODO: 改为从 viper 读取 log.level，避免生产环境硬编码 Debug）
	zerolog.SetGlobalLevel(zerolog.Level(0))

	// 2. 加载配置（支持 ${ENV_VAR} 环境变量展开）
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

	// ========== 核心依赖链构建（显式 DI） ==========

	// 6. 初始化数据库（返回实例，不再使用全局变量）
	db, err := database.Setup()
	if err != nil {
		return fmt.Errorf("setup database failed: %w", err)
	}

	// 7. 初始化 DAO 聚合体（注入 db 实例）
	repos := dao.NewRepositories(db)

	// 8. 初始化缓存层（注入 repos 作为缓存 Miss 时的降级数据源）
	caches := cache.NewCaches(repos)

	// 9. 初始化 Service 层
	// ⚠️ 过渡期：暂保留全局 Srv 赋值，供 bus handler / cron job 等未完成 DI 改造的模块使用
	// TODO: 所有消费者改为构造注入后，删除 service.Srv 全局变量
	service.Srv = service.NewServices(repos, caches)

	// 10. 启动定时任务
	cron.Setup()

	// 11. 路由注册（通过参数接收服务实例，路由层不再直接 import service 包）
	engine := gin.Default()
	router.Setup(engine, mgr, service.Srv)

	port := viper.GetString("base.port")
	addr := ":" + port

	// ========== 优雅退出逻辑 ==========
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

	// ========== 按依赖逆序关闭资源 ==========
	cron.Stop()
	fmt.Println("✅ Cron jobs stopped")

	// ⚠️ 过渡期：仍为全局函数调用，后续应改为 caches.Shutdown()
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