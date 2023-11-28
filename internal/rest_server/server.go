package rest_server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/charliemcelfresh/charlie-go/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	httpPort = "8080"
)

type Repository interface {
	GetItems(ctx context.Context, page int) ([]Item, error)
}

type server struct {
	Repository Repository
}

func NewServer() server {
	r := NewRepository()
	return server{
		Repository: r,
	}
}

func (s server) Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(AddUserIDToContext)
	r.Use(AddContentTypeToResponse)
	r.Get("/api/v1/items", s.GetItems)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
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

func (s server) GetItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	page := chi.URLParam(r, "page")
	pageAsInt, _ := strconv.Atoi(page)
	items, err := s.Repository.GetItems(ctx, pageAsInt)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	}
	jsonItems, err := json.Marshal(items)
	if err != nil {
		panic(err)
	}
	w.Write(jsonItems)
}
