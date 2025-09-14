package usecase

import (
	"context"

	"github.com/qushot/gin-todo-api/internal/domain/model"
	"github.com/qushot/gin-todo-api/internal/domain/repository"
)

// GetAllTodosUseCase はすべてのTodoを取得するユースケースを表すインターフェース
type GetAllTodosUseCase interface {
	Execute(ctx context.Context, query model.TodoQuery) ([]model.Todo, error)
}

// getAllTodosUseCase はGetAllTodosUseCaseの実装
type getAllTodosUseCase struct {
	todoRepo repository.TodoRepository
}

// NewGetAllTodosUseCase は全てのTodoを取得するユースケースのコンストラクタ
func NewGetAllTodosUseCase(todoRepo repository.TodoRepository) GetAllTodosUseCase {
	return &getAllTodosUseCase{
		todoRepo: todoRepo,
	}
}

// Execute は全てのTodoを取得する
func (uc *getAllTodosUseCase) Execute(ctx context.Context, query model.TodoQuery) ([]model.Todo, error) {
	return uc.todoRepo.FindAll(ctx, query)
}
