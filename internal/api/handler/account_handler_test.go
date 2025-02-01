package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VieiraVitor/transaction-flow/internal/api/dto"
	"github.com/VieiraVitor/transaction-flow/internal/domain"
	"github.com/VieiraVitor/transaction-flow/internal/infra/repository"
	"github.com/VieiraVitor/transaction-flow/internal/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAccountHandler_CreateAccount_WhenValidRequest_ShouldReturn201(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAccountUseCase(ctrl)
	hdlr := NewAccountHandler(mockUseCase)

	documentNumber := "12345678900"
	expectedID := int64(1)

	router := chi.NewRouter()
	router.Post("/accounts", hdlr.CreateAccount)

	reqBody, _ := json.Marshal(dto.CreateAccountRequest{DocumentNumber: documentNumber})
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUseCase.EXPECT().
		CreateAccount(gomock.Any(), documentNumber).
		Return(expectedID, nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var id int64
	err := json.Unmarshal(w.Body.Bytes(), &id)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
}

func TestAccountHandler_CreateAccount_WhenFailedToCreateAccount_ShouldReturn500(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAccountUseCase(ctrl)
	hdlr := NewAccountHandler(mockUseCase)

	router := chi.NewRouter()
	router.Post("/accounts", hdlr.CreateAccount)

	reqBody, _ := json.Marshal(dto.CreateAccountRequest{DocumentNumber: "12345678900"})
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	errorExpected := errors.New("could not create account")

	mockUseCase.EXPECT().
		CreateAccount(context.Background(), "12345678900").
		Return(int64(0), errorExpected)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAccountHandler_CreateAccount_WhenDocumentNumberIsNull_ShouldReturn422(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAccountUseCase(ctrl)
	hdlr := NewAccountHandler(mockUseCase)

	router := chi.NewRouter()
	router.Post("/accounts", hdlr.CreateAccount)

	reqBody, _ := json.Marshal(dto.CreateAccountRequest{DocumentNumber: ""})
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestAccountHandler_CreateAccount_WhenRequestIsInvalid_ShouldReturn400(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAccountUseCase(ctrl)
	hdlr := NewAccountHandler(mockUseCase)

	router := chi.NewRouter()
	router.Post("/accounts", hdlr.CreateAccount)

	invalidJSON := []byte(`{"document_number": 123abc}`)

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	var responseData map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseData)
	assert.NoError(t, err)

	assert.Contains(t, responseData, "error")
	assert.Equal(t, "invalid request", responseData["error"])
}

func TestAccountHandler_GetAccount_ShouldReturn200_WhenAccountExists(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAccountUseCase(ctrl)
	hdlr := NewAccountHandler(mockUseCase)

	account := dto.GetAccountResponse{AccountID: 1, DocumentNumber: "12345678900"}

	mockUseCase.EXPECT().
		GetAccount(context.Background(), int64(1)).
		Return(&domain.Account{ID: 1, DocumentNumber: "12345678900"}, nil)

	router := chi.NewRouter()
	router.Get("/accounts/{id}", hdlr.GetAccount)
	req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusAccepted, w.Code)

	var response dto.GetAccountResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, account, response)
}

func TestAccountHandler_GetAccount_WhenNotFoundAccount_ShouldReturn404(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAccountUseCase(ctrl)
	hdlr := NewAccountHandler(mockUseCase)

	mockUseCase.EXPECT().
		GetAccount(context.Background(), int64(1)).
		Return(nil, repository.ErrAccountNotFound)

	router := chi.NewRouter()
	router.Get("/accounts/{id}", hdlr.GetAccount)
	req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAccountHandler_GetAccount_WhenFailedToGetAccount_ShouldReturn500(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAccountUseCase(ctrl)
	hdlr := NewAccountHandler(mockUseCase)

	errorExpected := errors.New("could not create account")

	mockUseCase.EXPECT().
		GetAccount(context.Background(), int64(1)).
		Return(nil, errorExpected)

	router := chi.NewRouter()
	router.Get("/accounts/{id}", hdlr.GetAccount)
	req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAccountHandler_GetAccount_WhenInvalidID_ShouldReturn400(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAccountUseCase(ctrl)
	hdlr := NewAccountHandler(mockUseCase)

	router := chi.NewRouter()
	router.Get("/accounts/{id}", hdlr.GetAccount)

	req := httptest.NewRequest(http.MethodGet, "/accounts/abc", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var responseData map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseData)
	assert.NoError(t, err)

	assert.Contains(t, responseData, "error")
	assert.Equal(t, "could not parse id", responseData["error"].(string))
}
