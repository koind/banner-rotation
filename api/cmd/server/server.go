package server

import (
	"github.com/koind/banner-rotation/api/internal/config"
	"github.com/koind/banner-rotation/api/internal/transport"
	"github.com/spf13/cobra"
	"os"
)

// Declaring commands to start server
var RunServerCmd = &cobra.Command{
	Use:   "run",
	Short: "Run server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Init(config.Path)
		s := transport.NewServer(cfg)
		serverType := os.Getenv("SERVER_TYPE")

		s.Run(serverType)
	},
}

// When initializing parse the path to the configuration
func init() {
	RunServerCmd.Flags().StringVarP(
		&config.Path,
		"config",
		"c",
		"config/development/config.toml",
		"Path to toml configuration file",
	)
}
