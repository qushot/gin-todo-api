package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/qushot/gin-todo-api/internal/infrastructure/db"
	"github.com/qushot/gin-todo-api/internal/infrastructure/logger"
	"github.com/qushot/gin-todo-api/internal/infrastructure/server"
)

func main() {
	// ロガーの初期化
	logger.InitializeLogger()

	// データベースへの接続
	_, err := db.InitializeDB("postgres://postgres:pass@localhost:5432/postgres")
	if err != nil {
		slog.Error("Failed to connect to database", slog.Any("error", err))
		return
	}

	// サーバーの作成と起動
	srv := server.NewServer()
	srv.SetupRoutes()

	if err := srv.Start(); err != nil {
		slog.Error("Failed to start server", slog.Any("error", err))
		return
	}

	// graceful shutdown
	if err := srv.GracefulShutdown(); err != nil {
		slog.Error("Failed to shutdown server", slog.Any("error", err))
		return
	}

	// Close the database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.CloseDB(ctx); err != nil {
		slog.Error("Failed to close database connection", slog.Any("error", err))
		return
	}

	slog.Info("Server exiting")
}
