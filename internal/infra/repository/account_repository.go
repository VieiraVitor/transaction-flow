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

var ErrDuplicateAccount = errors.New("account already exists")

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *accountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) CreateAccount(ctx context.Context, account domain.Account) (int64, error) {
	query := "INSERT INTO accounts (document_number, created_at) VALUES ($1, NOW()) RETURNING id"
	var id int64
	row := r.db.QueryRow(query, account.DocumentNumber)
	err := row.Scan(&id)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "accounts_document_number_key"` {
			logger.Logger.Error("account already exists", slog.String("error", err.Error()))
			return 0, ErrDuplicateAccount
		}
		logger.Logger.Error("error creating account", slog.String("error", err.Error()))
		return 0, fmt.Errorf("failed to create account: %w", err)
	}
	return id, err
}

func (r *accountRepository) GetAccount(ctx context.Context, accountID int64) (*domain.Account, error) {
	query := "SELECT id, document_number, created_at FROM accounts WHERE id = $1"
	row := r.db.QueryRow(query, accountID)

	account, err := r.scanAccount(row)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Logger.Error("account not found", slog.String("error", err.Error()))
			return nil, fmt.Errorf("account not found")
		}
		logger.Logger.Error("error getting account", slog.String("error", err.Error()))
		return nil, err
	}

	return &account, err
}

func (r *accountRepository) scanAccount(row *sql.Row) (domain.Account, error) {
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
		return domain.Account{}, fmt.Errorf("unable to scan account: %w", err)
	}

	return domain.Account{
		ID:             id.Int64,
		DocumentNumber: documentNumber.String,
		CreatedAt:      createdAt.Time,
	}, nil
}
