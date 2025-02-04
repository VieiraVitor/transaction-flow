package handler

import (
	"context"
	"encoding/json"
	"fmt"
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

// CreateTransaction godoc
// @Summary Create a transaction
// @Description Registers a new financial transaction
// @Tags Transactions
// @Accept  json
// @Produce  json
// @Param transaction body dto.CreateTransactionRequest true "Transaction Request"
// @Success 201 {object} dto.CreateTransactionResponse "Transaction Created"
// @Failure 400 {object} response.ErrorResponse "Invalid Request"
// @Failure 422 {object} response.ErrorResponse "Validation Failed"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid request", fmt.Sprintf("malformed request :%v", err))
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

	transactionResponse := dto.CreateTransactionResponse{ID: transactionID}
	response.SendJSONResponse(context.Background(), w, http.StatusCreated, transactionResponse)
}
