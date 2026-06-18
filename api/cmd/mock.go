package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"ultrathreads/database"
	"ultrathreads/mock"
	"ultrathreads/repository"
	"ultrathreads/util/log"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var CmdMock = &cli.Command{
	Name:   "mock",
	Usage:  "Mock Data",
	Action: runMock,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "yes",
			Aliases: []string{"y"},
			Usage:   "Skip confirmation prompt",
		},
	},
}

func runMock(c *cli.Context) error {
	// 1. Set up log level
	zerolog.SetGlobalLevel(zerolog.Level(0))

	conf := "./app.yaml"
	if c.String("conf") != "" {
		conf = c.String("conf")
	}

	// 2. Set up configuration
	viper.SetConfigFile(conf)

	content, err := os.ReadFile(conf)
	if err != nil {
		return fmt.Errorf("read conf file fail: %w", err)
	}

	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		return fmt.Errorf("parse conf file fail: %w", err)
	}

	// 6. 初始化数据库（返回实例，不再使用全局变量）
	db, err := database.Setup()
	if err != nil {
		return fmt.Errorf("setup database failed: %w", err)
	}

	// 7. 初始化 DAO 聚合体（注入 db 实例）
	_ = repository.NewRepositories(db)

	// ✅ 交互确认（支持 --yes / -y 跳过）
	if !c.Bool("yes") {
		fmt.Print("⚠️  This will overwrite existing mock data. Continue? [Y/n]: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if input != "" && input != "y" && input != "yes" {
			fmt.Println("❌ Mock cancelled.")
			return nil
		}
	}

	log.Info("run mock\n")
	mock.Mock(db)

	return nil
}
