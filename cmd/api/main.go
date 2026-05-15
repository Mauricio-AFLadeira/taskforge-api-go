package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/mauricio-reportei/taskforge-api-go/internal/config"
	"github.com/mauricio-reportei/taskforge-api-go/internal/database"
	taskforgeredis "github.com/mauricio-reportei/taskforge-api-go/internal/redis"
	"github.com/mauricio-reportei/taskforge-api-go/internal/server"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "error", err)
		os.Exit(1)
	}

	setupLogging(cfg.AppEnv)

	ctx := context.Background()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	rdb := taskforgeredis.NewClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	defer func() { _ = rdb.Close() }()

	if err := taskforgeredis.Ping(ctx, rdb); err != nil {
		slog.Error("redis", "error", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%d", cfg.AppPort)
	httpSrv := &http.Server{
		Addr:         addr,
		Handler:      server.New(addr, pool, rdb),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("listening", "addr", addr, "env", cfg.AppEnv)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	slog.Info("shutdown")
	if err := httpSrv.Shutdown(shCtx); err != nil {
		slog.Error("shutdown", "error", err)
	}
}

func setupLogging(appEnv string) {
	level := slog.LevelInfo
	if appEnv == "development" {
		level = slog.LevelDebug
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))
}
