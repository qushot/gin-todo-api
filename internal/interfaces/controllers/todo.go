package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/qushot/gin-todo-api/internal/domain/model"
	"github.com/qushot/gin-todo-api/internal/usecase"
)

// TodoController はTodo操作のためのコントローラー
type TodoController struct {
	getAllTodosUseCase usecase.GetAllTodosUseCase
	getTodoByIDUseCase usecase.GetTodoByIDUseCase
	createTodoUseCase  usecase.CreateTodoUseCase
	updateTodoUseCase  usecase.UpdateTodoUseCase
	deleteTodoUseCase  usecase.DeleteTodoUseCase
}

// NewTodoController はTodoControllerのコンストラクタ
func NewTodoController(
	getAllTodosUseCase usecase.GetAllTodosUseCase,
	getTodoByIDUseCase usecase.GetTodoByIDUseCase,
	createTodoUseCase usecase.CreateTodoUseCase,
	updateTodoUseCase usecase.UpdateTodoUseCase,
	deleteTodoUseCase usecase.DeleteTodoUseCase,
) *TodoController {
	return &TodoController{
		getAllTodosUseCase: getAllTodosUseCase,
		getTodoByIDUseCase: getTodoByIDUseCase,
		createTodoUseCase:  createTodoUseCase,
		updateTodoUseCase:  updateTodoUseCase,
		deleteTodoUseCase:  deleteTodoUseCase,
	}
}

// RegisterRoutes はルーティング設定を行う
func (c *TodoController) RegisterRoutes(router *gin.RouterGroup) {
	todoRoutes := router.Group("/todos")
	{
		todoRoutes.GET("", c.List)
		todoRoutes.POST("", c.Create)
		todoRoutes.GET("/:id", c.Read)
		todoRoutes.PUT("/:id", c.Update)
		todoRoutes.DELETE("/:id", c.Delete)
	}
}

// List は全てのTodoを取得するハンドラー
func (c *TodoController) List(ctx *gin.Context) {
	var query model.TodoQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todos, err := c.getAllTodosUseCase.Execute(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, todos)
}

// Create は新しいTodoを作成するハンドラー
func (c *TodoController) Create(ctx *gin.Context) {
	var req model.Todo
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := c.createTodoUseCase.Execute(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, todo)
}

// Read は指定されたIDのTodoを取得するハンドラー
func (c *TodoController) Read(ctx *gin.Context) {
	id := ctx.Param("id")

	todo, err := c.getTodoByIDUseCase.Execute(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, errors.New("not found")) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

// Update は指定されたIDのTodoを更新するハンドラー
func (c *TodoController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req model.Todo
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := c.updateTodoUseCase.Execute(ctx.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, errors.New("not found")) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

// Delete は指定されたIDのTodoを削除するハンドラー
func (c *TodoController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.deleteTodoUseCase.Execute(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, errors.New("not found")) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
