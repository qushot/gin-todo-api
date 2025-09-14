package usecase

import (
	"context"

	"github.com/qushot/gin-todo-api/internal/domain/model"
	"github.com/qushot/gin-todo-api/internal/domain/repository"
)

// GetTodoByID はIDによってTodoを取得するユースケースを表すインターフェース
type GetTodoByID interface {
	Execute(ctx context.Context, id string) (*model.Todo, error)
}

// getTodoByID は usecase.GetTodoByID の実装
type getTodoByID struct {
	todoRepo repository.Todo
}

// NewGetTodoByID は usecase.GetTodoByID のコンストラクタ
func NewGetTodoByID(todoRepo repository.Todo) GetTodoByID {
	return &getTodoByID{
		todoRepo: todoRepo,
	}
}

// Execute はIDによるTodoの取得
func (uc *getTodoByID) Execute(ctx context.Context, id string) (*model.Todo, error) {
	return uc.todoRepo.FindByID(ctx, id)
}
