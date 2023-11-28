package cmd

import (
	"github.com/charliemcelfresh/charlie-go/internal/config"
	"github.com/spf13/cobra"
)

var seedDatabaseCmd = &cobra.Command{
	Use: "seed_database",
	Run: func(cmd *cobra.Command, args []string) {
		seedDatabase()
	},
}

func init() {
	rootCmd.AddCommand(seedDatabaseCmd)
}

func seedDatabase() {
	dbPool := config.GetDB()
	statement := `
		WITH u AS (
			INSERT INTO users (email) VALUES ('user@example.com') RETURNING id
		), i AS (
			INSERT INTO items (name) VALUES ('Widget 01') RETURNING id
		), u_and_i AS (
			SELECT
				u.id AS user_id, i.id AS item_id
			FROM
				u
			JOIN
				i ON 1=1
		)
		INSERT INTO
			user_items (user_id, item_id)
		SELECT
			user_id, item_id FROM u_and_i
    `
	_, err := dbPool.Exec(statement)
	if err != nil {
		panic(err)
	}
}
