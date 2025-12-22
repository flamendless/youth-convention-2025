package cmd

import (
	"os"
	"time"
	"youth-convention-2025/internal/logs"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use: "app",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	start := time.Now()
	defer func() {
		logs.Log().Info(
			"Metrics",
			zap.Duration("execution time", time.Since(start)),
		)
	}()
	if err := rootCmd.Execute(); err != nil {
		logs.Log().Error("error", zap.Error(err))
		os.Exit(1)
	}
}
