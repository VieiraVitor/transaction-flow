package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/VieiraVitor/transaction-flow/internal/domain"
	"github.com/VieiraVitor/transaction-flow/internal/infra/logger"
)

var ErrAccountNotFound = errors.New("account not found")

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *accountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) CreateAccount(ctx context.Context, account *domain.Account) (int64, error) {
	query := "INSERT INTO accounts (document_number) VALUES ($1) RETURNING id"
	var id int64
	row := r.db.QueryRow(query, account.DocumentNumber())
	err := row.Scan(&id)
	if err != nil {
		logger.Logger.ErrorContext(ctx, "error creating account", slog.String("document_number", account.DocumentNumber()), slog.String("error", err.Error()))
		return 0, fmt.Errorf("failed to create account: %w", err)
	}
	return id, err
}

func (r *accountRepository) GetAccount(ctx context.Context, accountID int64) (*domain.Account, error) {
	query := "SELECT id, document_number, created_at FROM accounts WHERE id = $1"
	row := r.db.QueryRow(query, accountID)

	account, err := r.scanAccount(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Logger.ErrorContext(ctx, "account not found", slog.Int64("account_id", accountID), slog.String("error", err.Error()))
			return nil, ErrAccountNotFound
		}
		logger.Logger.ErrorContext(ctx, "error getting account", slog.Int64("account_id", accountID), slog.String("error", err.Error()))
		return nil, err
	}

	return account, nil
}

func (r *accountRepository) scanAccount(row *sql.Row) (*domain.Account, error) {
	var (
		id             sql.NullInt64
		documentNumber sql.NullString
		createdAt      sql.NullTime
	)

	err := row.Scan(
		&id,
		&documentNumber,
		&createdAt,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to scan account: %w", err)
	}

	account := domain.NewAccount(documentNumber.String)
	account.SetID(id.Int64)
	account.SetCreatedAt(createdAt.Time)
	return account, nil
}
