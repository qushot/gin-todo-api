package repository

import (
	"context"

	"github.com/qushot/gin-todo-api/internal/domain/model"
)

// TodoRepository はTodoのデータ操作を担当するインターフェース
type TodoRepository interface {
	FindAll(ctx context.Context, query model.TodoQuery) ([]model.Todo, error)
	FindByID(ctx context.Context, id string) (*model.Todo, error)
	Create(ctx context.Context, todo model.Todo) (*model.Todo, error)
	Update(ctx context.Context, id string, todo model.Todo) (*model.Todo, error)
	Delete(ctx context.Context, id string) error
}
