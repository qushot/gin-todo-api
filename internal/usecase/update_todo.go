package usecase

import (
	"context"

	"github.com/qushot/gin-todo-api/internal/domain/model"
	"github.com/qushot/gin-todo-api/internal/domain/repository"
)

// UpdateTodo はTodoを更新するユースケースを表すインターフェース
type UpdateTodo interface {
	Execute(ctx context.Context, id string, todo model.Todo) (*model.Todo, error)
}

// updateTodo は usecase.UpdateTodo の実装
type updateTodo struct {
	todoRepo repository.Todo
}

// NewUpdateTodo は usecase.UpdateTodo のコンストラクタ
func NewUpdateTodo(todoRepo repository.Todo) UpdateTodo {
	return &updateTodo{
		todoRepo: todoRepo,
	}
}

// Execute はTodoを更新する
func (uc *updateTodo) Execute(ctx context.Context, id string, todo model.Todo) (*model.Todo, error) {
	return uc.todoRepo.Update(ctx, id, todo)
}
