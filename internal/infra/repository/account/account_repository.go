package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/VieiraVitor/transaction-flow/internal/domain"
	"github.com/VieiraVitor/transaction-flow/pkg/logger"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) CreateAccount(ctx context.Context, account domain.Account) (int, error) {
	query := "INSERT INTO accounts (document_number, created_at) VALUES ($1, NOW()) RETURNING id"
	var id int
	row := r.db.QueryRow(query, account.DocumentNumber)
	err := row.Scan(&id)
	if err != nil {
		logger.Logger.Error("error creating account", slog.String("error", err.Error()))
		return 0, fmt.Errorf("failed to create account: %w", err)
	}
	return id, err
}

func (r *AccountRepository) GetAccount(ctx context.Context, accountID int64) (*domain.Account, error) {
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

func (r *AccountRepository) scanAccount(row *sql.Row) (domain.Account, error) {
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
