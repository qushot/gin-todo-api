package usecase

import (
	"context"

	"github.com/qushot/gin-todo-api/internal/domain/model"
	"github.com/qushot/gin-todo-api/internal/domain/repository"
)

// CreateTodo は新しいTodoを作成するユースケースを表すインターフェース
type CreateTodo interface {
	Execute(ctx context.Context, todo model.Todo) (*model.Todo, error)
}

// createTodo は usecase.CreateTodo の実装
type createTodo struct {
	todoRepo repository.Todo
}

// NewCreateTodo は usecase.CreateTodo のコンストラクタ
func NewCreateTodo(todoRepo repository.Todo) CreateTodo {
	return &createTodo{
		todoRepo: todoRepo,
	}
}

// Execute は新しいTodoを作成する
func (uc *createTodo) Execute(ctx context.Context, todo model.Todo) (*model.Todo, error) {
	return uc.todoRepo.Create(ctx, todo)
}
