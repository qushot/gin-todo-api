package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DumpRequestBody は HTTP Request Body をダンプするミドルウェア
func DumpRequestBody(c *gin.Context) {
	ctx := c.Request.Context()

	defer c.Next()
	if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			slog.WarnContext(ctx, err.Error())
			return
		}

		// リクエストボディは一度読み込むとCloseされて読み込めなくなるため、読み込んだ内容を再度Bodyにセットする
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

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
}
