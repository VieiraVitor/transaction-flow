package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VieiraVitor/transaction-flow/internal/domain"
	"github.com/VieiraVitor/transaction-flow/internal/infra/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransactionRepositoryTestSuite struct {
	suite.Suite
	repo *transactionRepository
	mock sqlmock.Sqlmock
	db   *sql.DB
}

func (s *TransactionRepositoryTestSuite) SetupTest() {
	logger.InitLogger()
	var err error
	s.db, s.mock, err = sqlmock.New()
	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	s.repo = NewTransactionRepository(s.db)
}

func (s *TransactionRepositoryTestSuite) TearDownTest() {
	s.db.Close()
}

func TestTransactionRepositorySuite(t *testing.T) {
	suite.Run(t, new(TransactionRepositoryTestSuite))
}

func (s *TransactionRepositoryTestSuite) TestTransactionRepository_CreateTransaction_WhenValidInput_ShouldReturnId() {
	// Arrange
	transaction := domain.NewTransaction(int64(1), 1, 100)
	s.mock.ExpectQuery("INSERT INTO transactions").
		WithArgs(transaction.AccountID(), transaction.OperationTypeID(), transaction.Amount(), transaction.EventDate()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	ctx := context.Background()
	// Act
	id, err := s.repo.CreateTransaction(ctx, transaction)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), id)
	assert.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *TransactionRepositoryTestSuite) TestTransactionRepository_CreateTransaction_WhenFailedToCreateAccount_ShouldReturnError() {
	// Arrange
	transaction := domain.NewTransaction(int64(1), 1, 100)
	expectedError := errors.New("failed to create transaction")

	s.mock.ExpectQuery("INSERT INTO transactions").
		WithArgs(transaction.AccountID(), transaction.OperationTypeID(), transaction.Amount(), transaction.EventDate()).
		WillReturnError(expectedError)

	ctx := context.Background()
	// Act
	id, err := s.repo.CreateTransaction(ctx, transaction)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), int64(0), id)
	assert.Error(s.T(), expectedError, err)
	assert.NoError(s.T(), s.mock.ExpectationsWereMet())
}
