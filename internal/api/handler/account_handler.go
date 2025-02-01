package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/VieiraVitor/transaction-flow/internal/api/dto"
	"github.com/VieiraVitor/transaction-flow/internal/api/response"
	"github.com/VieiraVitor/transaction-flow/internal/application/usecase"
	"github.com/VieiraVitor/transaction-flow/internal/infra/repository"
	"github.com/go-chi/chi/v5"
)

type AccountHandler struct {
	useCase usecase.AccountUseCase
}

func NewAccountHandler(useCase usecase.AccountUseCase) *AccountHandler {
	return &AccountHandler{useCase: useCase}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, http.StatusUnprocessableEntity, "validation failed", err.Error())
		return
	}

	id, err := h.useCase.CreateAccount(context.Background(), req.DocumentNumber)
	if err != nil {
		response.SendErrorResponse(w, http.StatusInternalServerError, "could not create account", err.Error())
		return
	}
	response := dto.CreateAccountResponse{ID: id}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	accountID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, "could not parse id", err.Error())
		return
	}

	account, err := h.useCase.GetAccount(context.Background(), accountID)
	if err != nil {
		if errors.Is(err, repository.ErrAccountNotFound) {
			response.SendErrorResponse(w, http.StatusNotFound, "account not found", err.Error())
			return
		}
		response.SendErrorResponse(w, http.StatusInternalServerError, "could not get account", err.Error())
		return
	}

	accountResponse := &dto.GetAccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(accountResponse)
}
