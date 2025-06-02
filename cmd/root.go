package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	logger  *zap.Logger
	rootCmd = &cobra.Command{
		Use:   "pomodoro",
		Short: "A CLI Pomodoro timer",
		Long:  `A command-line interface for managing Pomodoro sessions`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initLogger)
}

func initLogger() {
	var err error
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	logger, err = config.Build()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
}
