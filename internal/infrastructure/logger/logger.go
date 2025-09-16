package logger

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

// customHandler is a custom handler for slog.JSONHandler
type customHandler struct {
	slog.Handler
}

// Handle is a override implementation of slog.Handler.Handle
func (h *customHandler) Handle(ctx context.Context, r slog.Record) error {
	sc := trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		projectID := cmp.Or(os.Getenv("GOOGLE_CLOUD_PROJECT"), "unknown")
		r.AddAttrs(
			slog.String("logging.googleapis.com/trace", fmt.Sprintf("projects/%s/traces/%s", projectID, sc.TraceID().String())),
			slog.String("logging.googleapis.com/spanId", sc.SpanID().String()),
		)
	}

	return h.Handler.Handle(ctx, r)
}

// WithAttrs is a override implementation of slog.Handler.WithAttrs
func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &customHandler{
		Handler: h.Handler.WithAttrs(attrs),
	}
}

// WithGroup is a override implementation of slog.Handler.WithGroup
func (h *customHandler) WithGroup(name string) slog.Handler {
	return &customHandler{
		Handler: h.Handler.WithGroup(name),
	}
}

// Initialize はロガーを初期化する
func Initialize() {
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
	logger := slog.New(&customHandler{
		Handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource:   true,
			Level:       slog.LevelDebug,
			ReplaceAttr: replace,
		}),
	})
	slog.SetDefault(logger)

	slog.Info("Logger initialized")
}
