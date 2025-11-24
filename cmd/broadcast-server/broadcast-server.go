package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/qs-lzh/broadcast-server/internal/client"
	"github.com/qs-lzh/broadcast-server/internal/server"
)

func main() {
	Execute()
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(connectCmd)
}

var rootCmd = &cobra.Command{
	Use:   "broadcast-server",
	Short: "A simple broadcast WebSocket server",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the broadcast server",
	Run: func(cmd *cobra.Command, args []string) {
		addr := "localhost:9000"
		if err := server.StartBCServer(addr); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	},
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a broadcast server as a client",
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.StartClient(); err != nil {
			log.Fatalf("failed to start client: %v", err)
		}
	},
}
