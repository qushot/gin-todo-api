package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/propagation"
)

// TraceContext is a middleware that extracts the trace context from the incoming request and sets it in the request context
func TraceContext(c *gin.Context) {
	tc := propagation.TraceContext{}
	ctx := tc.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}
