package testutils

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VieiraVitor/transaction-flow/internal/api/handler"
	"github.com/VieiraVitor/transaction-flow/internal/application/usecase"
	"github.com/VieiraVitor/transaction-flow/internal/infra/logger"
	"github.com/VieiraVitor/transaction-flow/internal/infra/repository"
	"github.com/go-chi/chi/v5"
	pq "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "transactions"
)

type TestContext struct {
	DB         *sql.DB
	Router     *chi.Mux
	AccountIDs []int64
}

func SetupTest(t *testing.T) *TestContext {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err, "failed to connect to database")

	logger.InitLogger()

	accountRepo := repository.NewAccountRepository(db)
	accountUseCase := usecase.NewAccountUseCase(accountRepo)
	accountHandler := handler.NewAccountHandler(accountUseCase)
	transactionRepo := repository.NewTransactionRepository(db)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase)

	router := chi.NewRouter()
	assert.NotNil(t, router, "router should not be nil")
	router.Post("/accounts", accountHandler.CreateAccount)
	router.Get("/accounts/{id}", accountHandler.GetAccount)
	router.Post("/transactions", transactionHandler.CreateTransaction)

	return &TestContext{DB: db, Router: router, AccountIDs: []int64{}}
}

func CleanupTest(t *testing.T, setup *TestContext) {
	_, err := setup.DB.Exec("DELETE FROM transactions WHERE account_id = ANY($1)", pq.Array(setup.AccountIDs))
	assert.NoError(t, err, "failed to clean up transactions")

	_, err = setup.DB.Exec("DELETE FROM accounts WHERE id = ANY($1)", pq.Array(setup.AccountIDs))
	assert.NoError(t, err, "failed to clean up accounts")

	setup.DB.Close()
}

func CreateRequest(t *testing.T, method, url string, body interface{}) (*httptest.ResponseRecorder, *http.Request) {
	var requestBody []byte
	var err error

	if body != nil {
		requestBody, err = json.Marshal(body)
		assert.NoError(t, err, "failed to marshal request body")
	}

	req := httptest.NewRequest(method, url, bytes.NewReader(requestBody))
	assert.NotNil(t, req, "Request should not be nil")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	return w, req
}
