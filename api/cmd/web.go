package cmd

import (
	"fmt"
	"os"
	"strings"

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
	Name:        "web",
	Usage:       "Start UltraThreads API",
	Action:      runWeb,
	Flags:       []cli.Flag{},
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
	fmt.Printf("🚀 UltraThreads starting on :%s\n", port)
	
	if err := engine.Run(":" + port); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}