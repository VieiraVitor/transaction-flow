package response

import (
	"encoding/json"
	"net/http"
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

	json.NewEncoder(w).Encode(resp)
}
