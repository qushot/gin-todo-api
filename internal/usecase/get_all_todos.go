package usecase

import (
	"context"

	"github.com/qushot/gin-todo-api/internal/domain/model"
	"github.com/qushot/gin-todo-api/internal/domain/repository"
)

// GetAllTodos はすべてのTodoを取得するユースケースを表すインターフェース
type GetAllTodos interface {
	Execute(ctx context.Context, query model.TodoQuery) ([]model.Todo, error)
}

// getAllTodos は usecase.GetAllTodos の実装
type getAllTodos struct {
	todoRepo repository.Todo
}

// NewGetAllTodos は usecase.GetAllTodos のコンストラクタ
func NewGetAllTodos(todoRepo repository.Todo) GetAllTodos {
	return &getAllTodos{
		todoRepo: todoRepo,
	}
}

// Execute は全てのTodoを取得する
func (uc *getAllTodos) Execute(ctx context.Context, query model.TodoQuery) ([]model.Todo, error) {
	return uc.todoRepo.FindAll(ctx, query)
}
