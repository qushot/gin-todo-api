package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type todo struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

type todoQuery struct {
	Status string `form:"status"`
}

var conn *pgx.Conn

func init() {
	var err error
	conn, err = pgx.Connect(context.Background(), "postgres://postgres:pass@localhost:5432/postgres")
	if err != nil {
		log.Fatalf("pgx.Connect error: %v", err)
	}
}

func main() {
	defer conn.Close(context.Background())

	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	baseRouter := router.Group("/api/v1")
	{
		newTodoHandler(baseRouter).handle()
	}

	// TODO: graceful shutdown

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

	var todos []todo
	rows, err := conn.Query(context.Background(), "SELECT id, title, content, done FROM todo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for rows.Next() {
		var t todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Content, &t.Done); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, t)
	}
	c.JSON(http.StatusOK, todos)
}

func (*todoHandler) create(c *gin.Context) {
	var req todo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var t todo
	if err := conn.QueryRow(context.Background(), "INSERT INTO todo (title, content, done) VALUES ($1, $2, $3) RETURNING id, title, content, done", req.Title, req.Content, req.Done).Scan(&t.ID, &t.Title, &t.Content, &t.Done); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, t)
}

func (*todoHandler) read(c *gin.Context) {
	id := c.Param("id")

	var t todo
	if err := conn.QueryRow(context.Background(), "SELECT id, title, content, done FROM todo WHERE id = $1", id).Scan(&t.ID, &t.Title, &t.Content, &t.Done); err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (*todoHandler) update(c *gin.Context) {
	id := c.Param("id")
	var req todo
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var t todo
	if err := conn.QueryRow(context.Background(), "UPDATE todo SET title = $2, content = $3, done = $4 WHERE id = $1 RETURNING id, title, content, done", id, req.Title, req.Content, req.Done).Scan(&t.ID, &t.Title, &t.Content, &t.Done); err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (*todoHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := conn.QueryRow(context.Background(), "DELETE FROM todo WHERE id = $1 RETURNING id", id).Scan(); err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
