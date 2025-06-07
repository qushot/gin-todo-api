package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/propagation"
)

// DumpRequestBody は HTTP Request Body をダンプするミドルウェア
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

// Recover は panic から回復するミドルウェア
func Recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func(ctx context.Context) {
			if rcr := recover(); rcr != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("%+v", rcr))
				// TODO: エラーレスポンスを返す
				http.Error(w, "panic", http.StatusInternalServerError)
			}
		}(r.Context())

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// SetTraceContext is a middleware that extracts the trace context from the incoming request and sets it in the request context
func SetTraceContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tc := propagation.TraceContext{}
		ctx := tc.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SetupMiddlewares はGin用のミドルウェアを設定する
func SetupMiddlewares(c *gin.Context) {
	hf := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		c.Request = r
		c.Next()
	})
	m := SetTraceContext(DumpRequestBody(Recover(hf)))

	// Gin の ResponseWriter, Request を渡す
	m.ServeHTTP(c.Writer, c.Request)
}
