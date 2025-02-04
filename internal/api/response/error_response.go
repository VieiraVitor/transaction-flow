package response

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/VieiraVitor/transaction-flow/internal/infra/logger"
)

type ErrorResponse struct {
	StatusCode  int    `json:"status_code"`
	Error       string `json:"error"`
	Description string `json:"description"`
}

func SendErrorResponse(w http.ResponseWriter, statusCode int, errMsg, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := ErrorResponse{
		StatusCode:  statusCode,
		Error:       errMsg,
		Description: description,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Logger.ErrorContext(context.Background(), "failed to encode response", slog.String("error", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

}
