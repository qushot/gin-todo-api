package main

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
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
	// ログメッセージのkeyをCloud Logging向けに変更
	replace := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.MessageKey {
			// メッセージが空文字の場合は除去する
			if a.Value.String() == "" {
				return slog.Attr{}
			}
			a.Key = "message"
		}
		if a.Key == slog.LevelKey {
			a.Key = "severity"
			if a.Value.Any().(slog.Level) == slog.LevelWarn {
				a.Value = slog.StringValue("WARNING")
			}
		}
		if a.Key == slog.SourceKey {
			a.Key = "logging.googleapis.com/sourceLocation"
		}
		return a
	}
	logger := slog.New(NewCustomJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: replace,
	}))
	slog.SetDefault(logger)

	slog.Info("Logger initialized")
}

func init() {
	var err error
	conn, err = pgx.Connect(context.Background(), "postgres://postgres:pass@localhost:5432/postgres")
	if err != nil {
		slog.Error("pgx.Connect error", slog.Any("error", err))
		return
	}

	slog.Info("Database connected")
}

// customJSONHandler is a custom handler for slog.JSONHandler
type customJSONHandler struct {
	*slog.JSONHandler
}

// Handle is a override implementation of slog.Handler.Handle
func (h *customJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	sc := trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		projectID := cmp.Or(os.Getenv("GOOGLE_CLOUD_PROJECT"), "unknown")
		r.AddAttrs(
			slog.String("logging.googleapis.com/trace", fmt.Sprintf("projects/%s/traces/%s", projectID, sc.TraceID().String())),
			slog.String("logging.googleapis.com/spanId", sc.SpanID().String()),
		)
	}

	return h.JSONHandler.Handle(ctx, r)
}

// NewCustomJSONHandler is a factory method for SlogHandler
func NewCustomJSONHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	return &customJSONHandler{
		JSONHandler: slog.NewJSONHandler(w, opts),
	}
}

// HTTP Request Body をダンプするミドルウェア
func DumpRequestBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		defer next.ServeHTTP(w, r)

		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				slog.WarnContext(ctx, err.Error())
				return
			}

			// リクエストボディは一度読み込むとCloseされて読み込めなくなるため、読み込んだ内容を再度Bodyにセットする
			r.Body = io.NopCloser(bytes.NewReader(body))

			// JSON形式として不正な場合はそのまま文字列として出力する
			if json.Valid(body) {
				m := make(map[string]any)
				if err := json.Unmarshal(body, &m); err != nil {
					slog.WarnContext(ctx, err.Error())
				}
				slog.InfoContext(ctx, "", slog.Any("requestBody", m))
				return
			}
			slog.InfoContext(ctx, "", slog.String("requestBody", string(body)))
		}
	})
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
		slog.Error("router.Run error", slog.Any("error", err))
		return
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
	if _, err := conn.Exec(context.Background(), "DELETE FROM todo WHERE id = $1", id); err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
