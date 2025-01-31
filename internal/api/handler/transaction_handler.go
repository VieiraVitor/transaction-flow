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
	var transactionRequest dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&transactionRequest); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, err.Error(), "invalid request")
		return
	}

	transactionID, err := h.useCase.CreateTransaction(context.Background(), transactionRequest.AccountID, transactionRequest.OperationTypeID, transactionRequest.Amount)
	if err != nil {
		response.SendErrorResponse(w, http.StatusInternalServerError, err.Error(), "could not create transaction")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(transactionID)
}
