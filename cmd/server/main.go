package main

import (
	"context"
	"github.com/wyw14/cry058/internal/config"
	"github.com/wyw14/cry058/internal/platform/logger"
	"github.com/wyw14/cry058/internal/repository/memory"
	transport "github.com/wyw14/cry058/internal/transport/http"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()
	ports := memory.NewPorts()
	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: transport.NewRouter(ports).Handler(), ReadHeaderTimeout: 5 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go srv.ListenAndServe()
	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdown)
	logger.New().Info("server stopped")
}
