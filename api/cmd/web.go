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
	"ultrathreads/database"
	"ultrathreads/delivery/handler"
	"ultrathreads/model"
	"ultrathreads/repository"
	"ultrathreads/server"
	"ultrathreads/service"
	"ultrathreads/util/hashid"
	"ultrathreads/util/querybuilder"
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

	// 5. 设置 Gin 运行模式
	gin.SetMode(viper.GetString("mode"))

	// ========== 核心依赖链构建 ==========
	db, err := database.Setup()
	if err != nil {
		return fmt.Errorf("setup database failed: %w", err)
	}

	repos := repository.NewRepositories(db)

	// 构建 CacheLoaders，提供数据加载函数给 cache 层
	// 这样 cache 层不需要直接依赖 dao 层
	loaders := &cache.CacheLoaders{
		NodeLoader: func(nodeId int64) *model.Node {
			return repos.Node.Get(nodeId)
		},
		AllNodesLoader: func() []model.Node {
			return repos.Node.Find(querybuilder.NewQueryBuilder().
				Eq("status", model.StatusOk).
				Asc("sort_no").Desc("id"))
		},
		PostRecLoader: func() []model.Post {
			return repos.Post.Find(querybuilder.NewQueryBuilder().
				Eq("recommend", true).
				Eq("status", model.StatusOk).
				Limit(20).
				Desc("last_comment_time"))
		},
		TagLoader: func(tagId int64) *model.Tag {
			return repos.Tag.Get(tagId)
		},
		HotTagsLoader: func() []model.Tag {
			return repos.Tag.Find(querybuilder.NewQueryBuilder().
				Eq("status", model.StatusOk).
				Desc("id").
				Limit(10))
		},
		PostTagsLoader: func(postId int64) []model.Tag {
			postTags := repos.PostTag.Find(
				querybuilder.NewQueryBuilder().Where("post_id = ?", postId),
			)
			var tags []model.Tag
			for _, pt := range postTags {
				if tag := repos.Tag.Get(pt.TagId); tag != nil {
					tags = append(tags, *tag)
				}
			}
			return tags
		},
		UserLoader: func(userId int64) *model.User {
			return repos.User.Get(userId)
		},
		UserScoreLoader: func(userId int64) int {
			userScore := repos.UserScore.FindOne(querybuilder.NewQueryBuilder().Eq("user_id", userId))
			if userScore == nil {
				return 0
			}
			return userScore.Score
		},
		ReadStateLoader: func(userID, nodeID int64) int64 {
			return repos.UserReadState.GetLastReadAt(userID, nodeID)
		},
		UserStatesLoader: func(userID int64) map[int64]int64 {
			return repos.UserReadState.GetAllReadStates(userID)
		},
		UserCountLoader: func() int {
			return int(repos.User.Count(querybuilder.NewQueryBuilder()))
		},
		PostCountLoader: func() int {
			return int(repos.Post.Count(querybuilder.NewQueryBuilder()))
		},
		SettingLoader: func(key string) *model.Setting {
			return repos.Setting.GetByKey(key)
		},
		ArticleTagLoader: func(articleId int64) []int64 {
			articleTags := repos.ArticleTag.FindByArticleId(articleId)
			var tagIds []int64
			for _, articleTag := range articleTags {
				tagIds = append(tagIds, articleTag.TagId)
			}
			return tagIds
		},
	}

	caches := cache.NewCaches(loaders)

	svcs := service.NewServices(repos, caches, db)

	cron.Setup(svcs.Article, svcs.Post)

	handlers := handler.NewHandlers(svcs, caches, mgr)

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
