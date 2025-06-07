package logger

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

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

// InitializeLogger はロガーを初期化する
func InitializeLogger() {
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
