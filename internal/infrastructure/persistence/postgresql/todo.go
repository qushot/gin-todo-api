package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/qushot/gin-todo-api/internal/domain/model"
	"github.com/qushot/gin-todo-api/internal/domain/repository"
)

// Todo はPostgreSQLを使ったTodoの実装
type Todo struct {
	conn *pgx.Conn
}

// NewTodo は repository.Todo のコンストラクタ
func NewTodo(conn *pgx.Conn) repository.Todo {
	return &Todo{
		conn: conn,
	}
}

// FindAll は全てのTodoを取得する
func (r *Todo) FindAll(ctx context.Context, _ model.TodoQuery) ([]model.Todo, error) {
	var todos []model.Todo
	rows, err := r.conn.Query(ctx, "SELECT id, title, content, done FROM todo")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t model.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Content, &t.Done); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// FindByID はIDによるTodoの取得
func (r *Todo) FindByID(ctx context.Context, id string) (*model.Todo, error) {
	var t model.Todo
	err := r.conn.QueryRow(ctx, "SELECT id, title, content, done FROM todo WHERE id = $1", id).
		Scan(&t.ID, &t.Title, &t.Content, &t.Done)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("not found")
	} else if err != nil {
		return nil, err
	}

	return &t, nil
}

// Create は新しいTodoを作成する
func (r *Todo) Create(ctx context.Context, todo model.Todo) (*model.Todo, error) {
	var t model.Todo
	err := r.conn.QueryRow(ctx,
		"INSERT INTO todo (title, content, done) VALUES ($1, $2, $3) RETURNING id, title, content, done",
		todo.Title, todo.Content, todo.Done).
		Scan(&t.ID, &t.Title, &t.Content, &t.Done)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Update はTodoを更新する
func (r *Todo) Update(ctx context.Context, id string, todo model.Todo) (*model.Todo, error) {
	var t model.Todo
	err := r.conn.QueryRow(ctx,
		"UPDATE todo SET title = $2, content = $3, done = $4 WHERE id = $1 RETURNING id, title, content, done",
		id, todo.Title, todo.Content, todo.Done).
		Scan(&t.ID, &t.Title, &t.Content, &t.Done)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("not found")
	} else if err != nil {
		return nil, err
	}

	return &t, nil
}

// Delete はTodoを削除する
func (r *Todo) Delete(ctx context.Context, id string) error {
	cmdTag, err := r.conn.Exec(ctx, "DELETE FROM todo WHERE id = $1", id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("not found")
	}

	return nil
}
