package usecase

import (
	"context"

	"github.com/qushot/gin-todo-api/internal/domain/repository"
)

// DeleteTodo はTodoを削除するユースケースを表すインターフェース
type DeleteTodo interface {
	Execute(ctx context.Context, id string) error
}

// deleteTodo は usecase.DeleteTodo の実装
type deleteTodo struct {
	todoRepo repository.Todo
}

// NewDeleteTodo は usecase.DeleteTodo のコンストラクタ
func NewDeleteTodo(todoRepo repository.Todo) DeleteTodo {
	return &deleteTodo{
		todoRepo: todoRepo,
	}
}

// Execute はTodoを削除する
func (uc *deleteTodo) Execute(ctx context.Context, id string) error {
	return uc.todoRepo.Delete(ctx, id)
}
