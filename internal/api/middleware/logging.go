package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/VieiraVitor/transaction-flow/internal/infra/logger"
	"github.com/google/uuid"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		traceID := uuid.New().String()

		logger.Logger.Info("Request received",
			slog.String("traceID", traceID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remoteAddr", r.RemoteAddr),
		)

		next.ServeHTTP(w, r)

		logger.Logger.Info("Request finished",
			slog.String("traceID", traceID),
			slog.String("duration", time.Since(startTime).String()),
		)
	})
}
