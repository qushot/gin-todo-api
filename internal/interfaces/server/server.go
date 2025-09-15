package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/qushot/gin-todo-api/internal/infrastructure/db"
	"github.com/qushot/gin-todo-api/internal/infrastructure/persistence/postgresql"
	"github.com/qushot/gin-todo-api/internal/interfaces/controllers"
	"github.com/qushot/gin-todo-api/internal/interfaces/middleware"
	"github.com/qushot/gin-todo-api/internal/usecase"
)

// Server はAPIサーバーを表す
type Server struct {
	router *gin.Engine
	srv    *http.Server
}

// New は新しいServerを作成する
func New() *Server {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Use(middleware.SetupMiddlewares)

	return &Server{
		router: r,
	}
}

// SetupRoutes はルーティングを設定する
func (s *Server) SetupRoutes() {
	// リポジトリの作成
	todoRepo := postgresql.NewTodo(db.GetDBConn())

	// ユースケースの作成
	getAllTodosUseCase := usecase.NewGetAllTodos(todoRepo)
	getTodoByIDUseCase := usecase.NewGetTodoByID(todoRepo)
	createTodoUseCase := usecase.NewCreateTodo(todoRepo)
	updateTodoUseCase := usecase.NewUpdateTodo(todoRepo)
	deleteTodoUseCase := usecase.NewDeleteTodo(todoRepo)

	// コントローラーの作成
	todoController := controllers.NewTodo(
		getAllTodosUseCase,
		getTodoByIDUseCase,
		createTodoUseCase,
		updateTodoUseCase,
		deleteTodoUseCase,
	)

	// ルートグループの設定
	baseRouter := s.router.Group("/api/v1")
	{
		todoController.RegisterRoutes(baseRouter)
	}
}

// Start はサーバーを起動する
func (s *Server) Start() error {
	// HTTPサーバーの設定
	s.srv = &http.Server{
		Addr:              ":8080",
		Handler:           s.router,
		ReadHeaderTimeout: 5 * time.Second, // Slowloris攻撃対策
	}

	// サーバーの起動
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server.ListenAndServe error", slog.Any("error", err))
			return
		}
	}()

	return nil
}

// GracefulShutdown はサーバーを正常終了する
func (s *Server) GracefulShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		slog.Error("server.Shutdown error", slog.Any("error", err))
		return err
	}

	return nil
}
