package usecase

import (
	"context"

	"github.com/qushot/gin-todo-api/internal/domain/model"
	"github.com/qushot/gin-todo-api/internal/domain/repository"
)

// GetTodoByIDUseCase はIDによってTodoを取得するユースケースを表すインターフェース
type GetTodoByIDUseCase interface {
	Execute(ctx context.Context, id string) (*model.Todo, error)
}

// getTodoByIDUseCase はGetTodoByIDUseCaseの実装
type getTodoByIDUseCase struct {
	todoRepo repository.TodoRepository
}

// NewGetTodoByIDUseCase はIDによるTodoの取得ユースケースのコンストラクタ
func NewGetTodoByIDUseCase(todoRepo repository.TodoRepository) GetTodoByIDUseCase {
	return &getTodoByIDUseCase{
		todoRepo: todoRepo,
	}
}

// Execute はIDによるTodoの取得
func (uc *getTodoByIDUseCase) Execute(ctx context.Context, id string) (*model.Todo, error) {
	return uc.todoRepo.FindByID(ctx, id)
}
