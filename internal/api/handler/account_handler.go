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

// CreateAccount godoc
// @Summary Create an account
// @Description Creates a new account with a document number
// @Tags Accounts
// @Accept  json
// @Produce  json
// @Param account body dto.CreateAccountRequest true "Account Data"
// @Success 201 {object} dto.CreateAccountResponse "Account ID"
// @Failure 400 {object} response.ErrorResponse "Invalid Request"
// @Failure 409 {object} response.ErrorResponse "Account Already Exists"
// @Router /accounts [post]
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

// GetAccount godoc
// @Summary Get account by ID
// @Description Retrieves account information using an account ID
// @Tags Accounts
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {object} dto.GetAccountResponse "Account Data"
// @Failure 400 {object} response.ErrorResponse "Invalid Account ID"
// @Failure 404 {object} response.ErrorResponse "Account Not Found"
// @Router /accounts/{id} [get]
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
