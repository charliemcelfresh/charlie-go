package cmd

import (
	"github.com/charliemcelfresh/charlie-go/internal/rest_server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(restServerCmd)
}

var restServerCmd = &cobra.Command{
	Use: "rest_server",
	Run: func(cmd *cobra.Command, args []string) {
		RunRestServer()
	},
}

func RunRestServer() {
	s := rest_server.NewServer()
	s.Run()
}
