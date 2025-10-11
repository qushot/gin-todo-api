package inmemory

import (
	"context"
	"errors"
	"slices"

	"github.com/google/uuid"

	"github.com/qushot/gin-todo-api/internal/domain/model"
	"github.com/qushot/gin-todo-api/internal/domain/repository"
)

// Todo はPostgreSQLを使ったTodoの実装
type Todo struct {
	todos []model.Todo
}

// NewTodo は repository.Todo のコンストラクタ
func NewTodo() repository.Todo {
	return &Todo{
		todos: []model.Todo{
			{ID: "00000000-0000-4000-a000-000000000001", Title: "掃除", Content: "掃除をする", Done: true},
			{ID: "00000000-0000-4000-a000-000000000002", Title: "洗濯", Content: "洗濯をする", Done: false},
			{ID: "00000000-0000-4000-a000-000000000003", Title: "料理", Content: "料理をする", Done: false},
		},
	}
}

// FindAll は全てのTodoを取得する
func (r *Todo) FindAll(ctx context.Context, _ model.TodoQuery) ([]model.Todo, error) {
	return r.todos, nil
}

// FindByID はIDによるTodoの取得
func (r *Todo) FindByID(ctx context.Context, id string) (*model.Todo, error) {
	i := slices.IndexFunc(r.todos, func(t model.Todo) bool {
		return t.ID == id
	})
	if i == -1 {
		return nil, errors.New("not found")
	}
	return &r.todos[i], nil
}

// Create は新しいTodoを作成する
func (r *Todo) Create(ctx context.Context, todo model.Todo) (*model.Todo, error) {
	t := model.Todo{
		ID:      uuid.New().String(),
		Title:   todo.Title,
		Content: todo.Content,
		Done:    todo.Done,
	}
	r.todos = append(r.todos, t)
	return &t, nil
}

// Update はTodoを更新する
func (r *Todo) Update(ctx context.Context, id string, todo model.Todo) (*model.Todo, error) {
	i := slices.IndexFunc(r.todos, func(t model.Todo) bool {
		return t.ID == id
	})
	if i == -1 {
		return nil, errors.New("not found")
	}
	r.todos[i] = model.Todo{
		ID:      id,
		Title:   todo.Title,
		Content: todo.Content,
		Done:    todo.Done,
	}
	return &r.todos[i], nil
}

// Delete はTodoを削除する
func (r *Todo) Delete(ctx context.Context, id string) error {
	i := slices.IndexFunc(r.todos, func(t model.Todo) bool {
		return t.ID == id
	})
	if i == -1 {
		return errors.New("not found")
	}
	r.todos = append(r.todos[:i], r.todos[i+1:]...)
	return nil
}
