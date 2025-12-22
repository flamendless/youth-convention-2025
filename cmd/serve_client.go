package cmd

import (
	"youth-convention-2025/client"
	"youth-convention-2025/internal/ctx"
	"youth-convention-2025/internal/logs"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var ctxClient ctx.ClientFlags

func init() {
	f := serveClientCmd.Flags
	f().BoolVarP(&ctxClient.Secure, "secure", "s", true, "Use secure session")
	f().BoolVarP(&ctxClient.TSC, "transport security", "t", false, "Use transport security")
	f().BoolVarP(&ctxClient.DevMode, "dev mode", "d", false, "Dev Mode")
	f().StringVarP(&ctxClient.Address, "address", "a", "localhost", "Address")
	f().StringVarP(&ctxClient.Port, "port", "p", ":3001", "Port of the address")

	rootCmd.AddCommand(serveClientCmd)
}

var serveClientCmd = &cobra.Command{
	Use:   "serve_client",
	Short: "Run the client server",
	Run: func(cmd *cobra.Command, args []string) {
		logs.Log().Info(
			"Run the client server",
			zap.Bool("secure", ctxClient.Secure),
			zap.Bool("transport security", ctxClient.TSC),
			zap.Bool("dev mode", ctxClient.DevMode),
			zap.String("address", ctxClient.Address),
			zap.String("port", ctxClient.Port),
		)
		client.Serve(&ctxClient)
	},
}
