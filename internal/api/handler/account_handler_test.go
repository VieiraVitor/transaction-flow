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
	"github.com/VieiraVitor/transaction-flow/internal/api/response"
	"github.com/VieiraVitor/transaction-flow/internal/domain"
	"github.com/VieiraVitor/transaction-flow/internal/infra/repository"
	"github.com/VieiraVitor/transaction-flow/internal/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAccountHandler_CreateAccount_WhenAccounCreatedSuccessfully_ShouldReturn201(t *testing.T) {
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

	var response dto.CreateAccountResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, response.ID)
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
	var errorResponse response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, errorResponse.StatusCode)
	assert.Equal(t, "could not create account", errorResponse.Error)
	assert.Equal(t, "could not create account", errorResponse.Description)
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
	var errorResponse response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, errorResponse.StatusCode)
	assert.Equal(t, "validation failed", errorResponse.Error)
	assert.Equal(t, "document_number is mandatory", errorResponse.Description)

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
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errorResponse response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, errorResponse.StatusCode)
	assert.Equal(t, "invalid request", errorResponse.Error)
	assert.Equal(t, "invalid character 'a' after object key:value pair", errorResponse.Description)

}

func TestAccountHandler_GetAccount_WhenAccountExists_ShouldReturn200(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAccountUseCase(ctrl)
	hdlr := NewAccountHandler(mockUseCase)

	account := domain.NewAccount("12345678900")
	account.SetID(1)

	mockUseCase.EXPECT().
		GetAccount(context.Background(), int64(1)).
		Return(account, nil)

	router := chi.NewRouter()
	router.Get("/accounts/{id}", hdlr.GetAccount)
	req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.GetAccountResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, account.ID(), response.AccountID)
	assert.Equal(t, account.DocumentNumber(), response.DocumentNumber)
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
	var errorResponse response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, errorResponse.StatusCode)
	assert.Equal(t, "account not found", errorResponse.Error)
	assert.Equal(t, "account not found", errorResponse.Description)

}

func TestAccountHandler_GetAccount_WhenFailedToGetAccount_ShouldReturn500(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAccountUseCase(ctrl)
	hdlr := NewAccountHandler(mockUseCase)

	errorExpected := errors.New("could not get account")

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
	var errorResponse response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, errorResponse.StatusCode)
	assert.Equal(t, errorExpected.Error(), errorResponse.Error)
	assert.Equal(t, errorExpected.Error(), errorResponse.Description)

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
	var errorResponse response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, errorResponse.StatusCode)
	assert.Equal(t, "could not parse id", errorResponse.Error)
	assert.Equal(t, "strconv.ParseInt: parsing \"abc\": invalid syntax", errorResponse.Description)

}
