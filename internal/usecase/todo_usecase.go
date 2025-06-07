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
