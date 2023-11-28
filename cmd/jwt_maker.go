package cmd

import (
	"github.com/charliemcelfresh/charlie-go/internal/jwt_maker"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(generateJwtCmd)
}

var generateJwtCmd = &cobra.Command{
	Use: "generate_jwt",
	Run: func(cmd *cobra.Command, args []string) {
		// args[0] == duration
		// args[1] == users.id from database
		jwt_maker.Run(args[0], args[1])
	},
}
