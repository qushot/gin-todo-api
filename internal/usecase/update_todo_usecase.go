package usecase

import (
	"context"

	"github.com/qushot/gin-todo-api/internal/domain/model"
	"github.com/qushot/gin-todo-api/internal/domain/repository"
)

// UpdateTodoUseCase はTodoを更新するユースケースを表すインターフェース
type UpdateTodoUseCase interface {
	Execute(ctx context.Context, id string, todo model.Todo) (*model.Todo, error)
}

// updateTodoUseCase はUpdateTodoUseCaseの実装
type updateTodoUseCase struct {
	todoRepo repository.TodoRepository
}

// NewUpdateTodoUseCase はTodoを更新するユースケースのコンストラクタ
func NewUpdateTodoUseCase(todoRepo repository.TodoRepository) UpdateTodoUseCase {
	return &updateTodoUseCase{
		todoRepo: todoRepo,
	}
}

// Execute はTodoを更新する
func (uc *updateTodoUseCase) Execute(ctx context.Context, id string, todo model.Todo) (*model.Todo, error) {
	return uc.todoRepo.Update(ctx, id, todo)
}
