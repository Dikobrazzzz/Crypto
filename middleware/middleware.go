package middleware

import (
	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := otel.Tracer("Crypto-api")
		ctx, span := tracer.Start(c.Request.Context(), "HTTP "+c.Request.Method+" "+c.FullPath())
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		span.SetAttributes(attribute.Int("http.status", c.Writer.Status()))
	}
}
