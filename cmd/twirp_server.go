package cmd

import (
	"net/http"

	charlie_go "github.com/charliemcelfresh/go-by-charlie/rpc/charlie-go"

	"github.com/charliemcelfresh/go-by-charlie/internal/twirp_server"

	"github.com/spf13/cobra"

	"github.com/twitchtv/twirp"
)

func init() {
	rootCmd.AddCommand(twirpServerCmd)
}

var twirpServerCmd = &cobra.Command{
	Use: "twirp_server",
	Run: func(cmd *cobra.Command, args []string) {
		Run()
	},
}

func Run() {
	provider := twirp_server.NewProvider()

	chainHooks := twirp.ChainHooks(
		provider.AuthHooks(),
	)

	mux := http.NewServeMux()

	// POST http(s)://<host>/go_by_charlie.GoByCharlie/CreateItem
	// POST http(s)://<host>/go_by_charlie.GoByCharlie/GetItem
	handler := charlie_go.NewCharlieGoServer(provider, chainHooks)
	mux.Handle(
		handler.PathPrefix(), twirp_server.AddJwtTokenToContext(
			handler,
		),
	)

	http.ListenAndServe(":8080", mux)
}
