package cmd

import (
	"context"
	"fmt"
	"github.com/charliemcelfresh/charlie-go/internal/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charliemcelfresh/charlie-go/rpc/charlie-go"

	"github.com/charliemcelfresh/charlie-go/internal/twirp_server"

	"github.com/spf13/cobra"

	"github.com/twitchtv/twirp"
)

const (
	httpPort = ":8081"
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

	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: mux,
	}

	// POST http(s)://<host>/api/v1/charlie_go.CharlieGo/CreateItem
	// POST http(s)://<host>/api/v1/charlie_go.CharlieGo/GetItem
	handler := charlie_go.NewCharlieGoServer(provider, twirp.WithServerPathPrefix("/api/v1"), chainHooks)
	mux.Handle(
		handler.PathPrefix(), twirp_server.AddJwtTokenToContext(
			handler,
		),
	)
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Printf("Server listening on %s", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("ListenAndServe: %v", err))
		}
	}()

	<-stop

	config.GetLogger().Info("Shutting down httpServer ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		panic(fmt.Sprintf("Server graceful shutdown failed: %v", err))
	}

	config.GetLogger().Info("Server shutdown")
}
