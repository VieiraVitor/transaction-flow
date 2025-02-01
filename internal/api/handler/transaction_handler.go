package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/VieiraVitor/transaction-flow/internal/api/dto"
	"github.com/VieiraVitor/transaction-flow/internal/api/response"
	"github.com/VieiraVitor/transaction-flow/internal/application/usecase"
)

type TransactionHandler struct {
	useCase usecase.TransactionUseCase
}

func NewTransactionHandler(useCase usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		useCase: useCase,
	}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, http.StatusUnprocessableEntity, "validation failed", err.Error())
		return
	}

	transactionID, err := h.useCase.CreateTransaction(context.Background(), req.AccountID, req.OperationTypeID, req.Amount)
	if err != nil {
		response.SendErrorResponse(w, http.StatusInternalServerError, "could not create transaction", err.Error())
		return
	}

	response := dto.CreateTransactionResponse{ID: transactionID}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
