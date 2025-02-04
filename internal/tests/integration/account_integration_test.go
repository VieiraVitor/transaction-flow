package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/VieiraVitor/transaction-flow/internal/api/dto"
	"github.com/VieiraVitor/transaction-flow/internal/api/response"
	"github.com/VieiraVitor/transaction-flow/internal/tests/integration/testutils"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount_WhenValidInput_ShouldReturn201(t *testing.T) {
	// Arrange
	setup := testutils.SetupTest(t)
	defer testutils.CleanupTest(t, setup)
	accountRequest := dto.CreateAccountRequest{DocumentNumber: "01101101001"}
	w, req := testutils.CreateRequest(t, http.MethodPost, "/accounts", accountRequest)

	// Act
	setup.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var accountResponse dto.CreateAccountResponse
	err := json.Unmarshal(w.Body.Bytes(), &accountResponse)
	assert.NoError(t, err)
	assert.Greater(t, accountResponse.ID, int64(0))

	assertCreateAccount(setup, t, accountRequest, accountResponse)

	setup.AccountIDs = append(setup.AccountIDs, accountResponse.ID)
}

func TestCreateAccount_WhenDocumentNumberIsEmpty_ShouldReturn422(t *testing.T) {
	// Arrange
	setup := testutils.SetupTest(t)
	defer testutils.CleanupTest(t, setup)

	w, req := testutils.CreateRequest(t, http.MethodPost, "/accounts", dto.CreateAccountRequest{DocumentNumber: ""})

	// Act
	setup.Router.ServeHTTP(w, req)

	// Arrange
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var errorResponse response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, errorResponse.StatusCode)
	assert.Equal(t, "validation failed", errorResponse.Error)
	assert.Equal(t, "document_number is mandatory", errorResponse.Description)
}

func TestCreateAccount_WhenDuplicateDocumentNumber_ShouldReturn500(t *testing.T) {
	// Arrange
	setup := testutils.SetupTest(t)
	defer testutils.CleanupTest(t, setup)

	w, req := testutils.CreateRequest(t, http.MethodPost, "/accounts", dto.CreateAccountRequest{DocumentNumber: "7777777"})
	setup.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var accountResponse dto.CreateAccountResponse
	err := json.Unmarshal(w.Body.Bytes(), &accountResponse)
	assert.NoError(t, err)
	assert.Greater(t, accountResponse.ID, int64(0))

	setup.AccountIDs = append(setup.AccountIDs, accountResponse.ID)

	// Act - Try to create a second account with the same document_number
	w, req = testutils.CreateRequest(t, http.MethodPost, "/accounts", dto.CreateAccountRequest{DocumentNumber: "7777777"})
	setup.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var errorResponse response.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, errorResponse.StatusCode)
	assert.Equal(t, "could not create account", errorResponse.Error)
	assert.Equal(t, "failed to create account: pq: duplicate key value violates unique constraint \"accounts_document_number_key\"", errorResponse.Description)
}

func TestGetAccount_WhenAccountExists_ShouldReturn200(t *testing.T) {
	// Arrange
	setup := testutils.SetupTest(t)
	defer testutils.CleanupTest(t, setup)

	w, req := testutils.CreateRequest(t, http.MethodPost, "/accounts", dto.CreateAccountRequest{DocumentNumber: "123456789"})
	setup.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.CreateAccountResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	setup.AccountIDs = append(setup.AccountIDs, response.ID)

	w, req = testutils.CreateRequest(t, http.MethodGet, fmt.Sprintf("/accounts/%d", response.ID), nil)

	// Act
	setup.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var getResponse dto.GetAccountResponse
	err = json.Unmarshal(w.Body.Bytes(), &getResponse)
	assert.NoError(t, err)
	assert.Equal(t, response.ID, getResponse.AccountID)
	assert.Equal(t, "123456789", getResponse.DocumentNumber)
}

func TestGetAccount_WhenAccountDoesNotExist_ShouldReturn404(t *testing.T) {
	// Arrange
	setup := testutils.SetupTest(t)
	defer testutils.CleanupTest(t, setup)

	w, req := testutils.CreateRequest(t, http.MethodGet, "/accounts/9999999", nil)
	// Act
	setup.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	var errorResponse response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, errorResponse.StatusCode)
	assert.Equal(t, "account not found", errorResponse.Error)
	assert.Equal(t, "account not found", errorResponse.Description)
}

func TestGetAccount_WhenInvalidInput_ShouldReturn400(t *testing.T) {
	// Arrange
	setup := testutils.SetupTest(t)
	defer testutils.CleanupTest(t, setup)

	w, req := testutils.CreateRequest(t, http.MethodGet, "/accounts/number", nil)
	// Act
	setup.Router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errorResponse response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, errorResponse.StatusCode)
	assert.Equal(t, "could not parse id", errorResponse.Error)
	assert.Equal(t, "strconv.ParseInt: parsing \"number\": invalid syntax", errorResponse.Description)
}

func assertCreateAccount(setup *testutils.TestContext,
	t *testing.T,
	requestBody dto.CreateAccountRequest,
	accountResponse dto.CreateAccountResponse) {
	var (
		accountID      int64
		documentNumber string
	)
	err := setup.DB.QueryRow(
		"SELECT id, document_number amount FROM accounts WHERE id = $1",
		accountResponse.ID,
	).Scan(&accountID, &documentNumber)

	assert.NoError(t, err)
	assert.Equal(t, accountResponse.ID, accountID)
	assert.Equal(t, requestBody.DocumentNumber, documentNumber)
}
