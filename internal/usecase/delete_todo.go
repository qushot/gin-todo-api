package usecase

import (
	"context"

	"github.com/qushot/gin-todo-api/internal/domain/repository"
)

// DeleteTodoUseCase はTodoを削除するユースケースを表すインターフェース
type DeleteTodoUseCase interface {
	Execute(ctx context.Context, id string) error
}

// deleteTodoUseCase はDeleteTodoUseCaseの実装
type deleteTodoUseCase struct {
	todoRepo repository.TodoRepository
}

// NewDeleteTodoUseCase はTodoを削除するユースケースのコンストラクタ
func NewDeleteTodoUseCase(todoRepo repository.TodoRepository) DeleteTodoUseCase {
	return &deleteTodoUseCase{
		todoRepo: todoRepo,
	}
}

// Execute はTodoを削除する
func (uc *deleteTodoUseCase) Execute(ctx context.Context, id string) error {
	return uc.todoRepo.Delete(ctx, id)
}
