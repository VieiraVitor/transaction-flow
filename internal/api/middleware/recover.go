package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/VieiraVitor/transaction-flow/internal/infra/logger"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Logger.Error("Panic",
					slog.String("error", fmt.Sprintf("%v", err)),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("stacktrace", string(debug.Stack())),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
