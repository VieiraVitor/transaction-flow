package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VieiraVitor/transaction-flow/internal/infra/logger"
	"github.com/stretchr/testify/assert"
)

func TestRecoverMiddleware_ShouldHandlePanicAndReturn500(t *testing.T) {
	// Arrange
	logger.InitLogger()
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	recoveryMiddleware := RecoverMiddleware(panicHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Act
	recoveryMiddleware.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Internal Server Error")
}
