package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/VieiraVitor/transaction-flow/internal/api/dto"
	"github.com/VieiraVitor/transaction-flow/internal/api/response"
	"github.com/VieiraVitor/transaction-flow/internal/tests/integration/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction_WhenValidInput_ShouldReturn201(t *testing.T) {
	// Arrange
	setup := testutils.SetupTest(t)
	defer testutils.CleanupTest(t, setup)

	w, req := testutils.CreateRequest(t, http.MethodPost, "/accounts", dto.CreateAccountRequest{DocumentNumber: "01101101001"})
	setup.Router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var accountResponse dto.CreateAccountResponse
	err := json.Unmarshal(w.Body.Bytes(), &accountResponse)
	assert.NoError(t, err)
	setup.AccountIDs = append(setup.AccountIDs, accountResponse.ID)

	body := dto.CreateTransactionRequest{
		AccountID:       accountResponse.ID,
		OperationTypeID: 4,
		Amount:          100,
	}

	w, req = testutils.CreateRequest(t, http.MethodPost, "/transactions", body)

	// Act
	setup.Router.ServeHTTP(w, req)

	// Arrange
	assert.Equal(t, http.StatusCreated, w.Code)

	var transactionResponse dto.CreateTransactionResponse
	err = json.Unmarshal(w.Body.Bytes(), &transactionResponse)
	assert.NoError(t, err)
	assert.Greater(t, transactionResponse.ID, int64(0))

	assertCreateTransaction(setup, t, body, transactionResponse)
}

func TestCreateTransaction_WhenCreateMoreThanOneTransactionSuccessfully_ShouldReturn201(t *testing.T) {
	// Arrange
	setup := testutils.SetupTest(t)
	defer testutils.CleanupTest(t, setup)

	// Create new account
	w, req := testutils.CreateRequest(t, http.MethodPost, "/accounts", dto.CreateAccountRequest{DocumentNumber: "01101101001"})
	setup.Router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var response dto.CreateAccountResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	setup.AccountIDs = append(setup.AccountIDs, response.ID)

	// Act create 2 transactions
	body := dto.CreateTransactionRequest{
		AccountID:       response.ID,
		OperationTypeID: 4,
		Amount:          100,
	}

	w, req = testutils.CreateRequest(t, http.MethodPost, "/transactions", body)

	setup.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp dto.CreateTransactionResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Greater(t, resp.ID, int64(0))
	assertCreateTransaction(setup, t, body, resp)

	body2 := dto.CreateTransactionRequest{
		AccountID:       response.ID,
		OperationTypeID: 4,
		Amount:          100,
	}

	w, req = testutils.CreateRequest(t, http.MethodPost, "/transactions", body2)
	setup.Router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Arrange
	var resp2 dto.CreateTransactionResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp2)
	assert.NoError(t, err)
	assert.Greater(t, resp2.ID, int64(0))
	assertCreateTransaction(setup, t, body2, resp2)
}

func TestCreateTransaction_WhenInvalidInput_ShouldReturn201(t *testing.T) {
	// Arrange
	setup := testutils.SetupTest(t)
	defer testutils.CleanupTest(t, setup)

	body := dto.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 0,
		Amount:          100,
	}

	w, req := testutils.CreateRequest(t, http.MethodPost, "/transactions", body)

	// Act
	setup.Router.ServeHTTP(w, req)

	// Arrange
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var errorResponse response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, errorResponse.StatusCode)
	assert.Equal(t, "validation failed", errorResponse.Error)
	assert.Equal(t, "operationTypeID is mandatory", errorResponse.Description)
}

func TestCreateTransaction_WhenCreateTransactionFails_ShouldReturn500(t *testing.T) {
	// Arrange
	setup := testutils.SetupTest(t)

	w, req := testutils.CreateRequest(t, http.MethodPost, "/accounts", dto.CreateAccountRequest{DocumentNumber: "01101101001"})
	setup.Router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var accountResponse dto.CreateAccountResponse
	err := json.Unmarshal(w.Body.Bytes(), &accountResponse)
	assert.NoError(t, err)
	setup.AccountIDs = append(setup.AccountIDs, accountResponse.ID)

	body := dto.CreateTransactionRequest{
		AccountID:       1,
		OperationTypeID: 1,
		Amount:          100,
	}

	testutils.CleanupTest(t, setup) // Close the connection before trying to create the transaction
	w, req = testutils.CreateRequest(t, http.MethodPost, "/transactions", body)

	// Act
	setup.Router.ServeHTTP(w, req)

	// Arrange
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var errorResponse response.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, errorResponse.StatusCode)
	assert.Equal(t, "could not create transaction", errorResponse.Error)
	assert.Equal(t, "failed to create transaction: sql: database is closed", errorResponse.Description)
}

func assertCreateTransaction(setup *testutils.TestContext,
	t *testing.T,
	requestBody dto.CreateTransactionRequest,
	transactionResponse dto.CreateTransactionResponse) {
	var (
		transactionID   int64
		accountID       int64
		operationTypeID int
		amount          float64
	)
	err := setup.DB.QueryRow(
		"SELECT id, account_id, operation_type_id, amount FROM transactions WHERE id = $1",
		transactionResponse.ID,
	).Scan(&transactionID, &accountID, &operationTypeID, &amount)

	assert.NoError(t, err)
	assert.Equal(t, transactionResponse.ID, transactionID)
	assert.Equal(t, requestBody.AccountID, accountID)
	assert.Equal(t, requestBody.OperationTypeID, operationTypeID)
	assert.Equal(t, requestBody.Amount, amount)
}
