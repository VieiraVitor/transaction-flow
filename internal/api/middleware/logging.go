package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/VieiraVitor/transaction-flow/internal/infra/logger"
	"github.com/google/uuid"
)

// ResponseWriter customizado para capturar status code e tamanho da resposta
type responseLogger struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rl *responseLogger) WriteHeader(code int) {
	rl.statusCode = code
	rl.ResponseWriter.WriteHeader(code)
}

func (rl *responseLogger) Write(p []byte) (int, error) {
	rl.body.Write(p)
	return rl.ResponseWriter.Write(p)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		traceID := uuid.NewString()

		var reqBody bytes.Buffer
		if r.Body != nil && r.Method != http.MethodGet {
			bodyBytes, _ := io.ReadAll(r.Body)
			reqBody.Write(bodyBytes)
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes)) // Reset body
		}

		respLogger := &responseLogger{ResponseWriter: w}

		logger.Logger.Info("Request received",
			slog.String("traceID", traceID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("requestBody", reqBody.String()),
		)

		next.ServeHTTP(respLogger, r)

		logger.Logger.Info("Request finished",
			slog.String("traceID", traceID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("status", http.StatusText(respLogger.statusCode)),
			slog.Int("statusCode", respLogger.statusCode),
			slog.String("responseBody", respLogger.body.String()),
			slog.Duration("duration", time.Since(startTime)),
		)
	})
}
