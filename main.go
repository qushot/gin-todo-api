package main

import (
	"log"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

type todo struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

type todos struct {
	Todos []todo `json:"todos"`
}

type todoQuery struct {
	Status string `form:"status"`
}

var defaultTodos = todos{
	Todos: []todo{
		{ID: "1", Title: "title1", Content: "content1", Done: false},
		{ID: "2", Title: "title2", Content: "content2", Done: true},
	},
}

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	baseRouter := router.Group("/api/v1")
	{
		newTodoHandler(baseRouter).handle()
	}

	if err := router.Run(); err != nil {
		log.Fatalf("router.Run error: %v", err)
	}
}

type todoHandler struct {
	r *gin.RouterGroup
}

func newTodoHandler(base *gin.RouterGroup) *todoHandler {
	return &todoHandler{
		r: base.Group("/todos"),
	}
}

func (t *todoHandler) handle() {
	t.r.GET("", t.list)
	t.r.POST("", t.create)
	t.r.GET("/:id", t.read)
	t.r.PUT("/:id", t.update)
	t.r.DELETE("/:id", t.delete)
}

func (*todoHandler) list(c *gin.Context) {
	var query todoQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if query.Status != "" {
		wasDone := true
		if query.Status != "done" {
			wasDone = false
		}
		var todos []todo
		for _, t := range defaultTodos.Todos {
			if t.Done == wasDone {
				todos = append(todos, t)
			}
		}
		c.JSON(http.StatusOK, todos)
		return
	}
	c.JSON(http.StatusOK, defaultTodos)
}

func (*todoHandler) create(c *gin.Context) {
	var req todo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t := todo{
		ID:      strconv.Itoa(len(defaultTodos.Todos) + 1),
		Title:   req.Title,
		Content: req.Content,
		Done:    false,
	}
	defaultTodos.Todos = append(defaultTodos.Todos, t)
	c.JSON(http.StatusCreated, t)
}

func (*todoHandler) read(c *gin.Context) {
	id := c.Param("id")
	idx := slices.IndexFunc(defaultTodos.Todos, func(t todo) bool {
		return t.ID == id
	})
	if idx == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, defaultTodos.Todos[idx])
}

func (*todoHandler) update(c *gin.Context) {
	id := c.Param("id")
	var req todo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	idx := slices.IndexFunc(defaultTodos.Todos, func(t todo) bool {
		return t.ID == id
	})
	if idx == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	defaultTodos.Todos[idx] = todo{
		ID:      id,
		Title:   req.Title,
		Content: req.Content,
		Done:    req.Done,
	}
	c.JSON(http.StatusOK, defaultTodos.Todos[idx])
}

func (*todoHandler) delete(c *gin.Context) {
	id := c.Param("id")
	idx := slices.IndexFunc(defaultTodos.Todos, func(t todo) bool {
		return t.ID == id
	})
	if idx == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	defaultTodos.Todos = append(defaultTodos.Todos[:idx], defaultTodos.Todos[idx+1:]...)
	c.JSON(http.StatusOK, defaultTodos)
}
