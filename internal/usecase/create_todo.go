package usecase

import (
	"context"

	"github.com/qushot/gin-todo-api/internal/domain/model"
	"github.com/qushot/gin-todo-api/internal/domain/repository"
)

// CreateTodoUseCase は新しいTodoを作成するユースケースを表すインターフェース
type CreateTodoUseCase interface {
	Execute(ctx context.Context, todo model.Todo) (*model.Todo, error)
}

// createTodoUseCase はCreateTodoUseCaseの実装
type createTodoUseCase struct {
	todoRepo repository.TodoRepository
}

// NewCreateTodoUseCase は新しいTodoを作成するユースケースのコンストラクタ
func NewCreateTodoUseCase(todoRepo repository.TodoRepository) CreateTodoUseCase {
	return &createTodoUseCase{
		todoRepo: todoRepo,
	}
}

// Execute は新しいTodoを作成する
func (uc *createTodoUseCase) Execute(ctx context.Context, todo model.Todo) (*model.Todo, error) {
	return uc.todoRepo.Create(ctx, todo)
}
