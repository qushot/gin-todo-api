package di

import (
	"sync"

	"github.com/qushot/gin-todo-api/internal/domain/repository"
	"github.com/qushot/gin-todo-api/internal/infrastructure/db"
	"github.com/qushot/gin-todo-api/internal/infrastructure/persistence/postgresql"
	"github.com/qushot/gin-todo-api/internal/interfaces/controllers"
	"github.com/qushot/gin-todo-api/internal/usecase"
)

var (
	once sync.Once
	c    *container
)

type container struct {
	TodoRepo repository.Todo

	GetAllTodosUseCase usecase.GetAllTodos
	GetTodoByIDUseCase usecase.GetTodoByID
	CreateTodoUseCase  usecase.CreateTodo
	UpdateTodoUseCase  usecase.UpdateTodo
	DeleteTodoUseCase  usecase.DeleteTodo

	TodoController *controllers.Todo
}

func GetContainer() *container {
	once.Do(func() {
		// repositories
		dbConn := db.GetDBConn()
		todoRepo := postgresql.NewTodo(dbConn)

		// use cases
		getAllTodosUseCase := usecase.NewGetAllTodos(todoRepo)
		getTodoByIDUseCase := usecase.NewGetTodoByID(todoRepo)
		createTodoUseCase := usecase.NewCreateTodo(todoRepo)
		updateTodoUseCase := usecase.NewUpdateTodo(todoRepo)
		deleteTodoUseCase := usecase.NewDeleteTodo(todoRepo)

		// controllers
		todoController := controllers.NewTodo(
			getAllTodosUseCase,
			getTodoByIDUseCase,
			createTodoUseCase,
			updateTodoUseCase,
			deleteTodoUseCase,
		)

		c = &container{
			TodoRepo: todoRepo,

			GetAllTodosUseCase: getAllTodosUseCase,
			GetTodoByIDUseCase: getTodoByIDUseCase,
			CreateTodoUseCase:  createTodoUseCase,
			UpdateTodoUseCase:  updateTodoUseCase,
			DeleteTodoUseCase:  deleteTodoUseCase,

			TodoController: todoController,
		}
	})

	return c
}
