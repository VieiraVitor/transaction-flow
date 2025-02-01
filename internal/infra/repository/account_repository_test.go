package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VieiraVitor/transaction-flow/internal/domain"
	"github.com/VieiraVitor/transaction-flow/internal/infra/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AccountRepositoryTestSuite struct {
	suite.Suite
	repo *accountRepository
	mock sqlmock.Sqlmock
	db   *sql.DB
}

func (s *AccountRepositoryTestSuite) SetupTest() {
	logger.InitLogger()
	var err error
	s.db, s.mock, err = sqlmock.New()
	if err != nil {
		s.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	s.repo = NewAccountRepository(s.db)
}

func (s *AccountRepositoryTestSuite) TearDownTest() {
	s.db.Close()
}

func TestAccountRepositorySuite(t *testing.T) {
	suite.Run(t, new(AccountRepositoryTestSuite))
}

func (s *AccountRepositoryTestSuite) TestAccountRepository_CreateAccount_WhenValidInput_ShouldReturnID() {
	// Arrange
	ctx := context.Background()
	account := domain.NewAccount("12345678900", time.Now())

	s.mock.ExpectQuery("INSERT INTO accounts").
		WithArgs(account.DocumentNumber).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Act
	id, err := s.repo.CreateAccount(ctx, account)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), id)
	assert.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *AccountRepositoryTestSuite) TestAccountRepository_CreateAccount_WhenFailedToCreateAccount_ShouldReturnError() {
	// Arrange
	ctx := context.Background()
	account := domain.NewAccount("12345678900", time.Now())
	expectedError := errors.New("failed to create account")

	s.mock.ExpectQuery("INSERT INTO accounts").
		WithArgs(account.DocumentNumber).
		WillReturnError(expectedError)

	// Act
	id, err := s.repo.CreateAccount(ctx, account)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), int64(0), id)
	assert.Error(s.T(), expectedError, err)
	assert.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *AccountRepositoryTestSuite) TestAccountRepository_GetAccount_WhenAccountExists_ShouldReturnAccount() {
	// Arrange
	ctx := context.Background()

	s.mock.ExpectQuery("SELECT id, document_number, created_at FROM accounts WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "document_number", "created_at"}).
			AddRow(1, "12345678900", time.Now()))

	// Act
	account, err := s.repo.GetAccount(ctx, 1)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), account)
	assert.Equal(s.T(), int64(1), account.ID)
	assert.Equal(s.T(), "12345678900", account.DocumentNumber)
}

func (s *AccountRepositoryTestSuite) TestAccountRepository_GetAccount_WhenAccountNotFound_ShouldReturnError() {
	// Arrange
	ctx := context.Background()

	s.mock.ExpectQuery("SELECT id, document_number, created_at FROM accounts WHERE id = ?").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	// Act
	account, err := s.repo.GetAccount(ctx, 1)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), account)
	assert.Equal(s.T(), "account not found", err.Error())
}

func (s *AccountRepositoryTestSuite) TestAccountRepository_GetAccount_WhenFailToGetAccount_ShouldReturnError() {
	// Arrange
	ctx := context.Background()
	expectedError := errors.New("failed to create account")

	s.mock.ExpectQuery("SELECT id, document_number, created_at FROM accounts WHERE id = ?").
		WithArgs(1).
		WillReturnError(expectedError)

	// Act
	account, err := s.repo.GetAccount(ctx, 1)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), account)
	assert.Equal(s.T(), expectedError, err)
}
