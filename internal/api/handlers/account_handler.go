package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VieiraVitor/transaction-flow/internal/api/response"
	"github.com/VieiraVitor/transaction-flow/internal/application/usecase"
	"github.com/go-chi/chi/v5"
)

type AccountHandler struct {
	useCase *usecase.AccountUseCase
}

func NewAccountHandler(useCase *usecase.AccountUseCase) *AccountHandler {
	return &AccountHandler{useCase: useCase}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var accountRequest CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&accountRequest); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, err.Error(), "invalid request")
		return
	}

	id, err := h.useCase.CreateAccount(context.Background(), accountRequest.DocumentNumber)
	if err != nil {
		response.SendErrorResponse(w, http.StatusInternalServerError, err.Error(), "could not create account")
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(id)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	accountID, err := strconv.Atoi(idParam)
	if err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, err.Error(), "could not parse id")
		return
	}

	account, err := h.useCase.GetAccount(context.Background(), int64(accountID))
	if err != nil {
		response.SendErrorResponse(w, http.StatusInternalServerError, err.Error(), "could not get account")
		return
	}

	accountResponse := &GetAccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(accountResponse)
}
