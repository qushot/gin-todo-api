package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// slog setting
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}

			if a.Key == slog.LevelKey {
				a.Key = "severity"
			}

			return a
		},
	})
	slog.SetDefault(slog.New(h))

	// gin setting
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		slog.InfoContext(c.Request.Context(), "middleware")
	})

	r.GET("/ping", func(c *gin.Context) {
		s := "world"
		slog.InfoContext(c.Request.Context(), "called", slog.String("hello", s))
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
