package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VieiraVitor/transaction-flow/internal/api/dto"
	"github.com/VieiraVitor/transaction-flow/internal/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTransactionHandler_CreateTransaction_WhenValidRequest_ShouldReturn201(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTransactionUseCase(ctrl)
	hdlr := NewTransactionHandler(mockUseCase)

	accountID := int64(123)
	operationTypeID := 1
	amount := 100.0
	expectedID := int64(1)

	router := chi.NewRouter()
	router.Post("/transactions", hdlr.CreateTransaction)

	reqBody, _ := json.Marshal(dto.CreateTransactionRequest{AccountID: accountID, OperationTypeID: operationTypeID, Amount: amount})
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockUseCase.EXPECT().
		CreateTransaction(gomock.Any(), accountID, operationTypeID, amount).
		Return(int64(1), nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var id int64
	err := json.Unmarshal(w.Body.Bytes(), &id)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
}

func TestTransactionHandler_CreateTransaction_WhenInvalidRequest_ShouldReturn400(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTransactionUseCase(ctrl)
	hdlr := NewTransactionHandler(mockUseCase)

	router := chi.NewRouter()
	router.Post("/transactions", hdlr.CreateTransaction)

	invalidJSON := []byte(`{"account_id": 123abc, "operation_type_id": "mock", "amount": 0}`)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var responseData map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &responseData)
	assert.NoError(t, err)

	assert.Contains(t, responseData, "error")
	assert.Equal(t, "invalid request", responseData["error"])
}

func TestTransactionHandler_CreateTransaction_WhenFailedToCreateTransaction_ShouldReturn500(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTransactionUseCase(ctrl)
	hdlr := NewTransactionHandler(mockUseCase)

	accountID := int64(123)
	operationTypeID := 1
	amount := 100.0

	router := chi.NewRouter()
	router.Post("/transactions", hdlr.CreateTransaction)

	reqBody, _ := json.Marshal(dto.CreateTransactionRequest{AccountID: accountID, OperationTypeID: operationTypeID, Amount: amount})
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	errorExpected := errors.New("could not create account")

	mockUseCase.EXPECT().
		CreateTransaction(gomock.Any(), accountID, operationTypeID, amount).
		Return(int64(0), errorExpected)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestTransactionHandler_CreateTransaction_InvalidInputs_ShouldReturn422(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTransactionUseCase(ctrl)
	hdlr := NewTransactionHandler(mockUseCase)

	router := chi.NewRouter()
	router.Post("/transactions", hdlr.CreateTransaction)

	testCases := []struct {
		name           string
		requestBody    dto.CreateTransactionRequest
		expectedStatus int
	}{
		{
			name:           "When AccountID is 0",
			requestBody:    dto.CreateTransactionRequest{AccountID: 0, OperationTypeID: 1, Amount: 100},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "When OperationTypeID is 0",
			requestBody:    dto.CreateTransactionRequest{AccountID: 1, OperationTypeID: 0, Amount: 100},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "When Amount is 0",
			requestBody:    dto.CreateTransactionRequest{AccountID: 1, OperationTypeID: 1, Amount: 0},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
