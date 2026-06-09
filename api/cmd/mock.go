package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"ultrathreads/dao"
	"ultrathreads/mock"
	"ultrathreads/util/log"
)

var CmdMock = &cli.Command{
	Name:        "mock",
	Usage:       "Mock Data",
	Description: "A free, open-source, self-hosted forum software written in Go",
	Action:      runMock,
	Flags:       []cli.Flag{},
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

	// Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		return fmt.Errorf("parse conf file fail: %w", err)
	}

	dao.Setup()
	
	log.Info("run mock\n")
	mock.Mock()

	return nil
}