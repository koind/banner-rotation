package cmd

import (
	"github.com/koind/banner-rotation/api/cmd/server"
	"github.com/spf13/cobra"
	"log"
)

// Declaring root commands
var rootCmd = &cobra.Command{
	Use:   "banner-rotation",
	Short: "Microservice banner-rotation",
}

// Adds http and grpc server commands during initialization
func init() {
	rootCmd.AddCommand(server.RunServerCmd)
}

// Runs the application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
