package cmd

import (
	"github.com/koind/banner-rotation/api/cmd/server"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "banner-rotation",
	Short: "Microservice banner-rotation",
}

func init() {
	rootCmd.AddCommand(server.HttpServerCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
